package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

func main() {
	if err := run(); err != nil {
		cmd := path.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "%s: %s\n", cmd, err)
	}
}

func run() error {
	var (
		smtpHostArg = arg{"smtp-host", 'm'}
		smtpHost    string
		smtpUserArg = arg{"smtp-user", 'u'}
		smtpUser    string
		smtpPassArg = arg{"smtp-pass", 'p'}
		smtpPass    string
		mailRcptArg = arg{"mail-rcpt", 'r'}
		mailRcpt    string
		noteSiteArg = arg{"note-site", 's'}
		noteSite    string
	)

	fs := flag.NewFlagSet("main", flag.ContinueOnError)
	fs.StringVar(&smtpHost, smtpHostArg.name, smtpHost, "smtp host")
	fs.StringVar(&smtpUser, smtpUserArg.name, smtpUser, "smtp user")
	fs.StringVar(&smtpPass, smtpPassArg.name, smtpPass, "smtp pass")
	fs.StringVar(&mailRcpt, mailRcptArg.name, mailRcpt, "mail recipient")
	fs.StringVar(&noteSite, noteSiteArg.name, noteSite, "site to notice")

	args, err := preprocessArgs(
		fs, os.Args,
		smtpHostArg, smtpUserArg, smtpPassArg,
		mailRcptArg,
		noteSiteArg,
	)
	if err != nil {
		return err
	}
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	m := newMailSend(smtpHost, smtpUser, smtpPass, mailRcpt)
	o := newSiteObservation(noteSite)

	go handleSending(o.dsc, o.esc, m)

	return o.run()
}

type arg struct {
	name string
	alt  rune
}

type mailSend struct {
	host string
	user string
	pass string
	rcpt string
}

func newMailSend(host, user, pass, rcpt string) *mailSend {
	return &mailSend{
		host: host,
		user: user,
		pass: pass,
		rcpt: rcpt,
	}
}

func (s *mailSend) send() error {
	return nil
}

type data struct{}

type siteObservation struct {
	site string
	dsc  chan data
	esc  chan error
}

func newSiteObservation(site string) *siteObservation {
	return &siteObservation{
		site: site,
	}
}

func (o *siteObservation) run() error {
	return nil
}

type sender interface {
	send() error
}

func handleSending(dc chan data, ec chan error, s sender) {}

func preprocessArgs(fs *flag.FlagSet, osArgs []string, args ...arg) ([]string, error) {
	ret := osArgs[:1]
	for _, a := range osArgs[1:] {
		if a[0] != '-' || (len(a) > 1 && a[:2] == "--") {
			ret = append(ret, a)
			continue
		}

		f := fs.Lookup(a[1:])
		if f != nil {
			ret = append(ret, a)
			continue
		}

		for _, sub := range a[1:] {
			aname := findArgNameByAlt(args, sub)
			if aname == "" {
				ret = append(ret, "-"+string(sub))
				continue
			}
			ret = append(ret, "--"+aname)
		}
	}

	return ret, nil
}

func findArgNameByAlt(args []arg, alt rune) string {
	for _, a := range args {
		if a.alt != 0 && a.alt == alt {
			return a.name
		}
	}

	return ""
}
