package claims

type LoginClaims struct {
	Username string `json:"username"`
	SystemRole string `json:"role"`
	IssuedAtTime int64 `json:"issueTime"`
	ExpirationTime int64 `json:"expiry"`
}
