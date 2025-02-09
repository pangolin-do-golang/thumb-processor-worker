package email

type Adapter interface {
	Send(to []string, subject string, body string) error
}
