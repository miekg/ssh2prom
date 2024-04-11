package main

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"time"

	"go.science.ru.nl/log"
)

const journalctl = "journalctl"

func journalReader(unitName string) (io.ReadCloser, func() error, error) {
	cancel := func() error { return nil } // initialize as noop

	args := []string{"-u", unitName, "--no-hostname"} // only works with -o short-xxx options.
	args = append(args, "-f")                         // tail

	cmd := exec.Command(journalctl, args...)
	p, err := cmd.StdoutPipe()
	if err != nil {
		return nil, cancel, err
	}

	if err := cmd.Start(); err != nil {
		return nil, cancel, err
	}

	cancel = func() error {
		go func() {
			if err := cmd.Wait(); err != nil {
				log.Debugf("wait for %q failed: %s", journalctl, err)
			}
		}()
		return cmd.Process.Kill()
	}

	return p, cancel, nil
}

var ErrExpired = errors.New("timeout expired")

// journalFollow synchronously follows the io.Reader, writing each new journal entry to writer. The
// follow will continue until a single time.Time is received on the until channel (or it's closed).
func journalFollow(until <-chan time.Time, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	bufch := make(chan []byte)
	errch := make(chan error)

	go func() {
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				errch <- err
				return
			}
			bufch <- scanner.Bytes()
		}
		// When the context is Done() the 'until' channel is closed, this kicks in the defers in the GetContainerLogsHandler method.
		// this cleans up the journalctl, and closes all file descripters. Scan() then stops with an error (before any reads,
		// hence the above if err .. .isn't triggered). In the end this go-routine exits.
		// the error here is "read |0: file already closed".
	}()

	for {
		select {
		case <-until:
			return ErrExpired

		case err := <-errch:
			return err

		case buf := <-bufch:
			if _, err := w.Write(buf); err != nil {
				return err
			}
			// re-add the newline
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
}
