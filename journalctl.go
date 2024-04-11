package provider

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

	// Handle all the options.
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
func journalFollow(until <-chan time.Time, reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
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
			if _, err := writer.Write(buf); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, "\n"); err != nil {
				return err
			}
		}
	}
}

/* use as
   logsReader, cancel, err := journalReader(namespace, pod, container, opts)
   if err != nil {
           return errors.Wrap(err, "failed to get systemd journal logs reader")
   }
   defer logsReader.Close()
   defer cancel()

   // ResponseWriter must be flushed after each write.
   if _, ok := w.(writeFlusher); !ok {
           log.Warn("HTTP response writer does not support flushes")
   }
   fw := flushOnWrite(w)

   if !opts.Follow {
           io.Copy(fw, logsReader)
           return nil
   }

   // If in follow mode, follow until interrupted.
   untilTime := make(chan time.Time, 1)
   errChan := make(chan error, 1)

   go func(w io.Writer, errChan chan error) {
           err := journalFollow(untilTime, logsReader, w)
           if err != nil && err != ErrExpired {
                   err = errors.Wrap(err, "failed to follow systemd journal logs")
           }
           errChan <- err
   }(fw, errChan)

   // Stop following logs if request context is completed.
   select {
   case err := <-errChan:
           return err
   case <-r.Context().Done():
           close(untilTime)
   }
   return nil
*/
