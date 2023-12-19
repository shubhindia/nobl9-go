// Code generated by "generate-object-impl Direct,PublicDirect"; DO NOT EDIT.

package direct

import (
	"github.com/nobl9/nobl9-go/manifest"
	"github.com/nobl9/nobl9-go/manifest/v1alpha"
)

// Ensure interfaces are implemented.
var _ manifest.Object = Direct{}
var _ manifest.ProjectScopedObject = Direct{}
var _ v1alpha.ObjectContext = Direct{}

func (d Direct) GetVersion() string {
	return d.APIVersion
}

func (d Direct) GetKind() manifest.Kind {
	return d.Kind
}

func (d Direct) GetName() string {
	return d.Metadata.Name
}

func (d Direct) Validate() error {
	if err := validate(d); err != nil {
		return err
	}
	return nil
}

func (d Direct) GetManifestSource() string {
	return d.ManifestSource
}

func (d Direct) SetManifestSource(src string) manifest.Object {
	d.ManifestSource = src
	return d
}

func (d Direct) GetProject() string {
	return d.Metadata.Project
}

func (d Direct) SetProject(project string) manifest.Object {
	d.Metadata.Project = project
	return d
}

func (d Direct) GetOrganization() string {
	return d.Organization
}

func (d Direct) SetOrganization(org string) manifest.Object {
	d.Organization = org
	return d
}