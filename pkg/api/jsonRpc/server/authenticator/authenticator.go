package authenticator

type Authenticator interface {
	Login(request *LoginRequest) (*LoginResponse, error)
	Logout(request *LogoutRequest) (*LogoutResponse, error)
}

const ServiceProvider = "Authenticator-Service"
const LoginService = ServiceProvider + ".Login"
const LogoutService = ServiceProvider + ".Logout"

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
