package provider

type Name string

type Provider interface {
	Name() Name
	MethodRequiresAuthorization(string) bool
}
