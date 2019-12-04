package email

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"

	core "github.com/jacob-ebey/graphql-core"
)

type Client interface {
	SendMail(to string, subject string, message string) error
}

func NewSmtpClient(from string, server string, auth smtp.Auth) *SmtpClient {
	return &SmtpClient{
		Auth:   auth,
		From:   from,
		Server: server,
	}
}

type SmtpClient struct {
	Auth   smtp.Auth
	From   string
	Server string
}

func (hook *SmtpClient) PreExecute(ctx context.Context, req core.GraphQLRequest) context.Context {
	return context.WithValue(ctx, "email", hook)
}

func (client *SmtpClient) SendMail(to string, subject string, message string) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	sbj := "Subject: " + subject + "\n"
	msg := []byte(sbj + mime + "\n" + message)

	err := smtp.SendMail(client.Server, client.Auth, client.From, []string{to}, msg)

	if err != nil {
		fmt.Println(err)
		return &core.WrappedError{
			Message:       "Could not send email.",
			InternalError: err,
		}
	}

	return nil
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}
