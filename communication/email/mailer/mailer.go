package mailer

type AuthInfo struct {
	Identity string
	Username string
	Password string
	Host     string
}

type Mailer interface {
	Send(request *SendRequest, response *SendResponse) error
}

type SendRequest struct {
	From    string
	To      string
	Cc      string
	Subject string
	Body    string
	Bcc     []string
}

type SendResponse struct {
}
