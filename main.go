package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/emersion/go-smtp"
)

var (
	listenADDR    = flag.String("l", "0.0.0.0", "listen address: default is 0.0.0.0 for all interfaces")
	listenPort    = flag.String("p", "1025", "listen port: default is 1025")
	tlsPort       = flag.String("tls", "1026", "listen port for TLS: default is 1026")
	tlsCert       = flag.String("cert", "cert.pem", "path to tls cert: default is cert.pem")
	tlsKey        = flag.String("key", "cert.key", "path to tls key: default is cert.key")
	tlsInsecure   = flag.Bool("k", true, "allow insecure tls [ true | false ] : default true")
	serverAddress string
	tlsAddress    string
	tlsEnabled    = false
	certKeyPair   tls.Certificate
)

func init() {
	flag.Parse()
	serverAddress = fmt.Sprintf("%s:%s", *listenADDR, *listenPort)
	err := validateTLSConfig()
	if err == nil {
		tlsEnabled = true
		tlsAddress = fmt.Sprintf("%s:%s", *listenADDR, *tlsPort)
	} else {
		log.Println(err)
	}
}

func validateTLSConfig() error {
	_, certErr := os.Stat(*tlsCert)
	_, keyErr := os.Stat(*tlsKey)
	if certErr == nil && keyErr == nil {
		pair, err := tls.LoadX509KeyPair(*tlsCert, *tlsKey)
		if err != nil {
			return fmt.Errorf("failed to load cert and key pair: %s", err)
		}
		certKeyPair = pair
		_, err = x509.ParseCertificate(pair.Certificate[0])
		if err != nil {
			return fmt.Errorf("failed to parse cert from pair: %s", err)
		}
	} else {
		return fmt.Errorf("cert and key:\ncertError: %s\nkeyError: %s", certErr, keyErr)
	}
	return nil
}

type backend struct{}

func (bkd *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{}, nil
}

type session struct{}

func (s *session) AuthPlain(username, password string) error {
	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	return nil
}

func (s *session) Data(r io.Reader) error {
	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}

func startTLS() {
	log.Println("Starting SMTPS server at", tlsAddress)
	s := smtp.NewServer(&backend{})
	s.Addr = tlsAddress
	s.Domain = "localhost"
	s.AllowInsecureAuth = *tlsInsecure
	s.Debug = os.Stdout
	s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{certKeyPair}}
	log.Fatal(s.ListenAndServeTLS())
}

func startHTTP() {
	log.Println("Starting SMTP server at", serverAddress)
	s := smtp.NewServer(&backend{})
	s.Addr = serverAddress
	s.Domain = "localhost"
	s.AllowInsecureAuth = true
	s.Debug = os.Stdout
	log.Fatal(s.ListenAndServe())
}

func main() {

	if tlsEnabled {
		go startTLS()
	}
	startHTTP()
}
