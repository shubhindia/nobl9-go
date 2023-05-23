package sdk

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nobl9/nobl9-go/sdk/retryhttp"
)

func TestCredentials_SetAuthorizationHeader(t *testing.T) {
	t.Run("access token is not set, do not set the header", func(t *testing.T) {
		creds := &Credentials{AccessToken: ""}
		req := &http.Request{}
		creds.SetAuthorizationHeader(req)
		assert.Empty(t, req.Header)
	})

	t.Run("set the header", func(t *testing.T) {
		creds := &Credentials{AccessToken: "123"}
		req := &http.Request{}
		creds.SetAuthorizationHeader(req)
		require.Contains(t, req.Header, HeaderAuthorization)
		assert.Equal(t, "Bearer 123", req.Header.Get(HeaderAuthorization))
	})
}

func TestCredentials_RefreshAccessToken(t *testing.T) {
	t.Run("don't run in offline mode", func(t *testing.T) {
		creds := &Credentials{offlineMode: true}
		tokenUpdated, err := creds.RefreshAccessToken(context.Background())
		require.NoError(t, err)
		assert.False(t, tokenUpdated)
	})

	for name, test := range map[string]struct {
		// Make sure 'exp' claim is float64, otherwise claims.VerifyExpiresAt will always return false.
		Claims       jwt.MapClaims
		TokenFetched bool
	}{
		"request new token for the first time": {
			// Claims are not set, we need a new token.
			Claims:       nil,
			TokenFetched: true,
		},
		"refresh the token if it's expired": {
			// If the expiry is set to now, the offset should still catch it.
			Claims:       jwt.MapClaims{"exp": float64(time.Now().Add(tokenExpiryOffset).Unix())},
			TokenFetched: true,
		},
		"don't refresh the token if it's not expired": {
			Claims:       jwt.MapClaims{"exp": float64(time.Now().Add(time.Hour).Unix())},
			TokenFetched: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			tokenProvider := &mockTokenProvider{}
			tokenParser := &mockTokenParser{}
			creds := &Credentials{
				claims:        test.Claims,
				TokenProvider: tokenProvider,
				TokenParser:   tokenParser,
			}
			_, err := creds.RefreshAccessToken(context.Background())
			require.NoError(t, err)
			expectedCalledTimes := 0
			if test.TokenFetched {
				expectedCalledTimes = 1
			}
			assert.Equal(t, expectedCalledTimes, tokenProvider.calledTimes)
			assert.Equal(t, expectedCalledTimes, tokenParser.calledTimes)
		})
	}

	t.Run("golden path, m2m token", func(t *testing.T) {
		tokenProvider := &mockTokenProvider{
			token: "access-token",
		}
		tokenParser := &mockTokenParser{
			claims: jwt.MapClaims{
				"m2mProfile": map[string]interface{}{
					"user":         "test@user.com",
					"environment":  "app.nobl9.com",
					"organization": "my-org",
				},
			},
		}
		hookCalled := false
		creds := &Credentials{
			ClientID:        "client-id",
			ClientSecret:    "my-secret",
			TokenProvider:   tokenProvider,
			TokenParser:     tokenParser,
			PostRequestHook: func(token string) error { hookCalled = true; return nil },
		}
		tokenUpdated, err := creds.RefreshAccessToken(context.Background())
		require.NoError(t, err)
		assert.True(t, tokenUpdated, "AccessToken must be updated")
		assert.True(t, hookCalled, "PostRequestHook must be called")
		assert.Equal(t, "my-secret", tokenProvider.calledWithClientSecret)
		assert.Equal(t, "client-id", tokenProvider.calledWithClientID)
		assert.Equal(t, "access-token", tokenParser.calledWithToken)
		assert.Equal(t, "client-id", tokenParser.calledWithClientID)
		assert.Equal(t, "access-token", creds.AccessToken)
		assert.Equal(t, "app.nobl9.com", creds.Environment)
		assert.Equal(t, "my-org", creds.Organization)
		assert.Equal(t, tokenTypeM2M, creds.tokenType)
		assert.Equal(t, accessTokenM2MProfile{
			User:         "test@user.com",
			Environment:  "app.nobl9.com",
			Organization: "my-org",
		}, creds.m2mProfile)
		assert.Equal(t, tokenParser.claims, creds.claims)
	})

	t.Run("golden path, agent token", func(t *testing.T) {
		tokenProvider := &mockTokenProvider{
			token: "access-token",
		}
		tokenParser := &mockTokenParser{
			claims: jwt.MapClaims{
				"agentProfile": map[string]interface{}{
					"user":         "test@user.com",
					"environment":  "app.nobl9.com",
					"organization": "my-org",
					"name":         "my-agent",
					"project":      "default",
				},
			},
		}
		hookCalled := false
		creds := &Credentials{
			ClientID:        "client-id",
			ClientSecret:    "my-secret",
			TokenProvider:   tokenProvider,
			TokenParser:     tokenParser,
			PostRequestHook: func(token string) error { hookCalled = true; return nil },
		}
		tokenUpdated, err := creds.RefreshAccessToken(context.Background())
		require.NoError(t, err)
		assert.True(t, tokenUpdated, "AccessToken must be updated")
		assert.True(t, hookCalled, "PostRequestHook must be called")
		assert.Equal(t, "my-secret", tokenProvider.calledWithClientSecret)
		assert.Equal(t, "client-id", tokenProvider.calledWithClientID)
		assert.Equal(t, "access-token", tokenParser.calledWithToken)
		assert.Equal(t, "client-id", tokenParser.calledWithClientID)
		assert.Equal(t, "access-token", creds.AccessToken)
		assert.Equal(t, "app.nobl9.com", creds.Environment)
		assert.Equal(t, "my-org", creds.Organization)
		assert.Equal(t, tokenTypeAgent, creds.tokenType)
		assert.Equal(t, accessTokenAgentProfile{
			User:         "test@user.com",
			Environment:  "app.nobl9.com",
			Organization: "my-org",
			Name:         "my-agent",
			Project:      "default",
		}, creds.agentProfile)
		assert.Equal(t, tokenParser.claims, creds.claims)
	})
}

func TestCredentials_SetAccessToken(t *testing.T) {
	tokenParser := &mockTokenParser{
		claims: jwt.MapClaims{
			"m2mProfile": map[string]interface{}{
				"user":         "test@user.com",
				"environment":  "app.nobl9.com",
				"organization": "my-org",
			},
		},
	}
	creds := &Credentials{
		ClientID:        "client-id",
		TokenParser:     tokenParser,
		PostRequestHook: func(token string) error { return errors.New("hook should not be called!") },
	}
	err := creds.SetAccessToken("access-token")
	require.NoError(t, err)
	assert.Equal(t, 1, tokenParser.calledTimes)
	assert.Equal(t, "access-token", tokenParser.calledWithToken)
	assert.Equal(t, "client-id", tokenParser.calledWithClientID)
	assert.Equal(t, "access-token", creds.AccessToken)
	assert.Equal(t, "app.nobl9.com", creds.Environment)
	assert.Equal(t, "my-org", creds.Organization)
	assert.Equal(t, tokenTypeM2M, creds.tokenType)
	assert.Equal(t, accessTokenM2MProfile{
		User:         "test@user.com",
		Environment:  "app.nobl9.com",
		Organization: "my-org",
	}, creds.m2mProfile)
	assert.Equal(t, tokenParser.claims, creds.claims)

	t.Run("offline mode", func(t *testing.T) {
		creds.offlineMode = true
		err = creds.SetAccessToken("new-access-token")
		require.NoError(t, err)
		assert.Equal(t, "access-token", creds.AccessToken)
		assert.Equal(t, 1, tokenParser.calledTimes)
	})
}

func TestCredentials_setNewToken(t *testing.T) {
	t.Run("don't call hook if parser fails", func(t *testing.T) {
		parserErr := errors.New("parser failed!")
		hookCalled := false
		creds := &Credentials{
			TokenParser:     &mockTokenParser{err: parserErr},
			PostRequestHook: func(token string) error { hookCalled = true; return nil },
		}
		err := creds.setNewToken("", true)
		require.Error(t, err)
		assert.Equal(t, parserErr, err)
		assert.False(t, hookCalled, "PostRequestHook should not be called")
	})

	t.Run("don't call hook If we can't decode m2mProfile", func(t *testing.T) {
		hookCalled := false
		creds := &Credentials{
			TokenParser:     &mockTokenParser{claims: jwt.MapClaims{"m2mProfile": "should be a map..."}},
			PostRequestHook: func(token string) error { hookCalled = true; return nil },
		}
		err := creds.setNewToken("", true)
		assert.Error(t, err)
		assert.False(t, hookCalled, "PostRequestHook should not be called")
	})

	t.Run("don't update Credentials state if hook fails", func(t *testing.T) {
		tokenParser := &mockTokenParser{
			claims: jwt.MapClaims{},
		}
		hookErr := errors.New("hook failed!")
		creds := &Credentials{
			TokenParser:     tokenParser,
			PostRequestHook: func(token string) error { return hookErr },
		}
		err := creds.setNewToken("my-token", true)
		require.Error(t, err)
		assert.ErrorIs(t, err, hookErr)
		assert.Empty(t, creds.AccessToken)
		assert.Empty(t, creds.m2mProfile)
		assert.Empty(t, creds.claims)
	})
}

// nolint: bodyclose
func TestCredentials_RoundTrip(t *testing.T) {
	t.Run("wrap errors with NonRetryableError", func(t *testing.T) {
		tokenProvider := &mockTokenProvider{err: errors.New("token fetching failed!")}
		creds := &Credentials{
			TokenProvider: tokenProvider,
			TokenParser:   &mockTokenParser{},
		}

		req := &http.Request{}
		_, err := creds.RoundTrip(req)
		require.Error(t, err)
		_, isNonRetryableError := err.(retryhttp.NonRetryableError)
		assert.True(t, isNonRetryableError, "err is of type NonRetryableError")
	})

	t.Run("set auth header if not present", func(t *testing.T) {
		creds := &Credentials{
			AccessToken:   "my-token",
			claims:        jwt.MapClaims{"exp": float64(time.Now().Add(time.Hour).Unix())},
			TokenProvider: &mockTokenProvider{},
			TokenParser:   &mockTokenParser{},
		}

		req := &http.Request{}
		_, _ = creds.RoundTrip(req)
		require.Contains(t, req.Header, HeaderAuthorization)
		assert.Equal(t, "Bearer my-token", req.Header.Get(HeaderAuthorization))
	})

	t.Run("update auth header if token was updated", func(t *testing.T) {
		creds := &Credentials{
			AccessToken:   "my-old-token",
			claims:        jwt.MapClaims{"exp": float64(time.Now().Unix())}, // expired
			TokenProvider: &mockTokenProvider{token: "my-new-token"},
			TokenParser:   &mockTokenParser{},
		}

		req := &http.Request{Header: http.Header{HeaderAuthorization: []string{"Bearer my-old-token"}}}
		_, _ = creds.RoundTrip(req)
		require.Contains(t, req.Header, HeaderAuthorization)
		assert.Equal(t, "Bearer my-new-token", req.Header.Get(HeaderAuthorization))
	})

	t.Run("don't update auth header if token was not updated", func(t *testing.T) {
		creds := &Credentials{
			AccessToken:   "my-old-token",
			claims:        jwt.MapClaims{"exp": float64(time.Now().Add(time.Hour).Unix())}, // not expired
			TokenProvider: &mockTokenProvider{token: "my-new-token"},
			TokenParser:   &mockTokenParser{},
		}

		req := &http.Request{Header: http.Header{HeaderAuthorization: []string{"Bearer my-old-token"}}}
		_, _ = creds.RoundTrip(req)
		require.Contains(t, req.Header, HeaderAuthorization)
		assert.Equal(t, "Bearer my-old-token", req.Header.Get(HeaderAuthorization))
	})
}

type mockTokenProvider struct {
	calledTimes            int
	calledWithClientID     string
	calledWithClientSecret string

	token string
	err   error
}

func (m *mockTokenProvider) RequestAccessToken(
	_ context.Context,
	clientID, clientSecret string,
) (token string, err error) {
	m.calledTimes++
	m.calledWithClientID = clientID
	m.calledWithClientSecret = clientSecret
	return m.token, m.err
}

type mockTokenParser struct {
	calledTimes        int
	calledWithClientID string
	calledWithToken    string

	claims jwt.MapClaims
	err    error
}

func (m *mockTokenParser) Parse(token, clientID string) (jwt.MapClaims, error) {
	m.calledTimes++
	m.calledWithToken = token
	m.calledWithClientID = clientID
	return m.claims, m.err
}