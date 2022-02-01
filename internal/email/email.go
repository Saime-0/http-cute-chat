package email

import (
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/validator"
	"net/smtp"
	"strconv"
)

type SMTPSender struct {
	msgAuthor string
	from      string
	pass      string
	host      string
	address   string
}

func NewSMTPSender(author, from, pass, host string, port int) (*SMTPSender, error) {
	if !validator.ValidateEmail(from) {
		return nil, errors.New("NewSMTPSender: smtp sender email not valid")
	}
	return &SMTPSender{
		msgAuthor: author,
		from:      from,
		pass:      pass,
		host:      host,
		address:   host + ":" + strconv.Itoa(port),
	}, nil
}

func (s *SMTPSender) Send(subject string, msgBody string, to ...string) error {
	if len(to) == 0 {
		return errors.New("empty to")
	}

	if subject == "" || msgBody == "" {
		return errors.New("empty subject/body")
	}

	for _, rec := range to {
		if !validator.ValidateEmail(rec) {
			return errors.New("invalid to email")
		}
	}

	// SetState up authentication information.
	auth := smtp.PlainAuth("", s.from, s.pass, s.host)
	msg := []byte(
		"From: " + s.msgAuthor + " <" + s.from + ">\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			msgBody + "\r\n",
	)
	err := smtp.SendMail(s.address, auth, s.from, to, msg)
	if err != nil {
		return errors.Wrap(err, "failed to sent email via smtp")
	}

	return nil
}
