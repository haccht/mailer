package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"gopkg.in/gomail.v2"
)

var opts struct {
	Address string   `short:"a" description:"smtp address" default:"localhost"`
	Port    int      `short:"p" description:"smtp port" default:"25"`
	Subject string   `short:"s" description:"subject" required:"true"`
	Sender  string   `short:"r" description:"sender" required:"true"`
	CCopy   []string `short:"c" description:"carbon copies"`
	BCCopy  []string `short:"b" description:"blind carbon copies"`
}

func main() {
	recipients, err := flags.Parse(&opts)
	if err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("To", recipients...)
	m.SetHeader("From", opts.Sender)
	m.SetHeader("Subject", opts.Subject)
	m.SetHeader("Cc", opts.CCopy...)
	m.SetHeader("Bcc", opts.BCCopy...)
	m.SetBody("text/plain", string(body))

	d := gomail.Dialer{Host: opts.Address, Port: opts.Port}
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
}
