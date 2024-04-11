package main

import "bytes"

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
	InvalidRoot = []byte("Failed password for ")
)

func (mw metricsWriter) Write(p []byte) (int, error) {
	i := bytes.Index(p, InvalidUser)
	if i > 0 {
		space := bytes.Index(p[len(InvalidUser)+1:], []byte(" "))
		if space != 0 {
			start := len(InvalidUser) + 1
			end := start + space
			println(string(p[start:end]))
			return len(p), nil
		}
	}

	return len(p), nil
}
