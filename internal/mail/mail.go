// Package mail facilitates the interaction between the lorafication daemon and
// the configured SMTP server.
package mail

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
)

// Mailer is a type that holds the SMTP auth, ready to send email using it's receiver
// functions after proper initialization using NewMailer.
type Mailer struct {
	addr string
	auth smtp.Auth
	from string
}

// NewMailer configures a Mailer to be used given the credentials for carrying out
// authentication.
func NewMailer(host string, port int, user, pass string) *Mailer {
	return &Mailer{
		addr: fmt.Sprintf("%s:%d", host, port),
		from: user,
		auth: smtp.PlainAuth("", user, pass, host),
	}
}

// DefaultHeaders returns the default headers necessary to send a properly formed
// email in string form.
func (m *Mailer) DefaultHeaders(to, subject string) string {
	headers := make(map[string]string)

	headers["From"] = m.from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", "text/html")
	headers["Content-Disposition"] = "inline"
	headers["Content-Transfer-Encoding"] = "quoted-printable"

	header := ""
	for key, value := range headers {
		header += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	return header
}

// Send takes a recipient address, subject, and message and uses them to send an email
// using the underlying receiver type, Mailer.
func (m *Mailer) Send(to, subject, msg string) error {
	var body bytes.Buffer
	qpw := quotedprintable.NewWriter(&body)

	if _, err := qpw.Write([]byte(msg)); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	if err := qpw.Close(); err != nil {
		return fmt.Errorf("close writer: %w", err)
	}

	// Add headers to body.
	finalBody := m.DefaultHeaders(to, subject) + "\r\n" + body.String()

	return smtp.SendMail(m.addr, m.auth, m.from, []string{to}, []byte(finalBody))
}
