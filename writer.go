package main

import (
	"bytes"

	"go.science.ru.nl/log"
)

// metricsWriter gets the line from the journald, and parses lines like:
//
//	apr 11 17:31:13 sshd[577691]: Failed password for root ...'
//
// Where 'root' (or any other user) is extracted.
type metricsWriter struct{}

// Failed password for invalid user usuario
// Failed password for root

var (
	InvalidUser = []byte("Failed password for invalid user ")
	InvalidRoot = []byte("Failed password for root")
)

func (mw metricsWriter) Write(p []byte) (int, error) {
	if len(p) < 3 {
		return len(p), nil
	}
	p = p[:len(p)-1]
	log.Debugf("%s", p)

	i := bytes.Index(p, InvalidUser)
	if i > 0 {
		space := bytes.Index(p[i+len(InvalidUser):], []byte(" "))
		if space != 0 {
			start := i + len(InvalidUser)
			end := start + space
			user := string(p[start:end])
			log.Debugf("User: %q", user)
			if !*flgDry {
				failedUserLogins.Inc()
			}
			return len(p), nil
		}
	}
	i = bytes.Index(p, InvalidRoot)
	if i > 0 {
		log.Debugf("User: %q", "root")
		if !*flgDry {
			failedRootLogins.Inc()
			failedUserLogins.Inc() // also inc the total
		}
	}
	return len(p), nil
}
