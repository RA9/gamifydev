// Package email sends transactional mail over SMTP. It is optional: when SMTP is
// not configured, Configured() returns false and callers fall back to showing
// the action link directly (e.g. an admin copies an invite link by hand).
package email

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// ErrNotConfigured is returned by Send when no SMTP host is configured.
var ErrNotConfigured = errors.New("email: SMTP not configured")

// Mailer holds SMTP settings read from the environment:
//
//	SMTP_HOST       smtp.example.com   (required to enable email)
//	SMTP_PORT       587                (default)
//	SMTP_USERNAME   login user         (optional; enables auth)
//	SMTP_PASSWORD   login password
//	SMTP_FROM       from address       (default: SMTP_USERNAME)
//	SMTP_FROM_NAME  display name       (default: GamifyDev)
type Mailer struct {
	host, port, user, pass, from, fromName string
}

// New builds a Mailer from environment variables.
func New() *Mailer {
	user := os.Getenv("SMTP_USERNAME")
	return &Mailer{
		host:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
		port:     envOr("SMTP_PORT", "587"),
		user:     user,
		pass:     os.Getenv("SMTP_PASSWORD"),
		from:     envOr("SMTP_FROM", user),
		fromName: envOr("SMTP_FROM_NAME", "GamifyDev"),
	}
}

// Configured reports whether enough is set to actually send mail.
func (m *Mailer) Configured() bool { return m.host != "" && m.from != "" }

// Send delivers a plain-text email. Returns ErrNotConfigured if SMTP is unset.
func (m *Mailer) Send(to, subject, body string) error {
	if !m.Configured() {
		return ErrNotConfigured
	}
	var auth smtp.Auth
	if m.user != "" {
		auth = smtp.PlainAuth("", m.user, m.pass, m.host)
	}
	addr := m.host + ":" + m.port
	return smtp.SendMail(addr, auth, m.from, []string{to}, m.message(to, subject, body))
}

func (m *Mailer) message(to, subject, body string) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s <%s>\r\n", m.fromName, m.from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	b.WriteString("\r\n")
	return []byte(b.String())
}

func envOr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
