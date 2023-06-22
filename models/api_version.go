package models

type (
	ApiVersion   int64
	ContextParam string
)

const (
	ApiVersionContext ContextParam = "apiVersion"
	IdentityContext   ContextParam = "identity"
	NotSupported      ApiVersion   = -1
	Alpha1            ApiVersion   = iota
)

func (v ApiVersion) String() string {
	switch v {
	case Alpha1:
		return "alpha1"
	}
	return "not supported"
}

func ToApiVersion(version string) ApiVersion {
	switch version {
	case "alpha1":
		return Alpha1
	}
	return NotSupported
}
