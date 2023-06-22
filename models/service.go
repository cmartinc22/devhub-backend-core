package models

type AuthResult struct {
	IsValid  bool
	Enabled  bool
	Identity *string
}
