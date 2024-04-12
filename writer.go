package main

import (
	"bytes"

	"go.science.ru.nl/log"
)

// metricsWriter gets the line from the journald, and parses lines like:
//
//	apr 11 17:31:13 sshd[577691]: Failed password for root from 61.177.172.136 port 13804 ssh2
//
// Where 'root' (or any other user) is extracted. The text after 'root' or 'invalid user XXXX' is checked
// for colons. If found the connection is assumed to be coming in over IPv6.
type metricsWriter struct{}

var (
	InvalidUser = []byte("Failed password for invalid user ") // xxxx from 61.177.172.136 port 13804 ssh2
	InvalidRoot = []byte("Failed password for root")          // from 61.177.172.136 port 13804 ssh2
	ValidUser   = []byte("session opened for user ")          // xxxx(uid=
)

var (
	Space = []byte(" ")
	Colon = []byte(":")
	Brace = []byte("(")
)

func (mw metricsWriter) Write(p []byte) (int, error) {
	if len(p) < 3 {
		return len(p), nil
	}
	p = p[:len(p)-1]
	family := "1"
	i := bytes.Index(p, InvalidUser)
	if i > 0 {
		log.Debugf("%s", p)
		space := bytes.Index(p[i+len(InvalidUser):], Space)
		if space != 0 {
			start := i + len(InvalidUser)
			end := start + space
			user := string(p[start:end])
			colon := bytes.Index(p[end:], Colon)
			if colon > 0 {
				family = "2"
			}
			log.Debugf("User: %q (fam=%s)", user, family)
			if !*flgDry {
				failedUserLogins.WithLabelValues(family).Inc()
			}
			return len(p), nil
		}
	}
	i = bytes.Index(p, InvalidRoot)
	if i > 0 {
		log.Debugf("%s", p)
		end := i + len(InvalidRoot)
		colon := bytes.Index(p[end:], Colon)
		if colon > 0 {
			family = "2"
		}
		log.Debugf("User: %q (fam=%s)", "root", family)
		if !*flgDry {
			failedRootLogins.WithLabelValues(family).Inc()
			failedUserLogins.WithLabelValues(family).Inc() // also inc the total
		}
		return len(p), nil
	}
	i = bytes.Index(p, ValidUser)
	if i > 0 {
		log.Debugf("%s", p)
		brace := bytes.Index(p[i+len(ValidUser):], Brace)
		if brace != 0 {
			start := i + len(ValidUser)
			end := start + brace
			user := string(p[start:end])
			log.Debugf("Valid user: %q", user)
			if !*flgDry {
				userLogins.Inc()
			}
			return len(p), nil
		}
	}
	return len(p), nil
}
