package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jhillyerd/enmime"
	"github.com/joho/godotenv"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/emersion/go-smtp"
)

// The Backend implements SMTP server methods.
type Backend struct{}

// Login handles a login command with username and password.
func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username != "username" || password != "password" {
		return nil, errors.New("Invalid username or password")
	}
	return &Session{}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	//return nil, smtp.ErrAuthRequired
	return &Session{}, nil
}

// A Session is returned after successful login.
type Session struct {
	receiver []string
	from     string
	message  string
	subject  string
}

func (s *Session) Mail(from string) error {
	log.Println("Mail from:", from)
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	s.receiver = append(s.receiver, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {

	// Parse message body with enmime.
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		fmt.Print(err)
		return err
	}

	s.subject = env.GetHeader("Subject")
	s.message = env.Text

	sendMail(s)

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

// spClient
var client sp.Client

func main() {

	godotenv.Load()

	be := &Backend{}

	err := getClient()
	if err != nil {
		return
	}
	s := smtp.NewServer(be)

	s.Domain = os.Getenv("SPARKPOST_DOMAIN")
	s.Addr = os.Getenv("ADDRESS")

	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	//s.AuthDisabled = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func sendMail(s *Session) {
	// Create a Transmission using an inline Recipient List
	// and inline email Content.
	tx := &sp.Transmission{
		Recipients: s.receiver,
		Content: sp.Content{
			HTML:    s.message,
			From:    s.from,
			Subject: s.subject,
		},
	}
	id, _, err := client.Send(tx)
	if err != nil {
		log.Fatal(err)
	}

	// The second value returned from Send
	// has more info about the HTTP response, in case
	// you'd like to see more than the Transmission id.
	log.Printf("Transmission sent with id [%s]\n", id)
}

func getClient() error {
	// Get our API key from the environment; configure.

	apiKey := os.Getenv("SPARKPOST_API_KEY")

	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     apiKey,
		ApiVersion: 1,
	}

	err := client.Init(cfg)
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
		return errors.New("SparkPost client init failed")
	}
	return nil
}
