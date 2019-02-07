package auth

type Service interface {
	Login(request *LoginRequest, response *LoginResponse) error
	Logout(request *LogoutRequest, response *LogoutResponse) error
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

type LoginRequest struct {
	UsernameOrEmailAddress string
	Password               string
}

type LoginResponse struct {
	Jwt string
}
