package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.science.ru.nl/log"
)

var (
	flgUnit  = flag.String("u", "ssh", "name of the ssh unit")
	flgDry   = flag.Bool("n", false, "dry run only show parsed lines")
	flgDebug = flag.Bool("d", false, "enable debug, show logs and parsed users")
	flgAddr  = flag.String("a", ":9396", "address to run prometheus exporter on")
)

func main() {
	flag.Parse()
	if *flgDebug {
		log.D.Set()
	}

	r, cancel, err := journalReader(*flgUnit)
	if err != nil {
		log.Fatalf("failed to get systemd journal logs reader: %s", err)
	}
	defer r.Close()
	defer cancel()

	// Follow until interrupted.
	untilTime := make(chan time.Time, 1)
	errChan := make(chan error, 1)

	var w io.Writer = os.Stdout
	if !*flgDry {
		w = metricsWriter{}
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Starting Prometheus on address %q", *flgAddr)
	go func() {
		err := http.ListenAndServe(*flgAddr, nil)
		if err != nil {
			log.Fatalf("Could not start HTTP server: %s", err)
		}
	}()

	go func(w io.Writer, errChan chan error) {
		err := journalFollow(untilTime, r, w)
		errChan <- err
	}(w, errChan)

	// Stop following logs if request context is completed.
	select {
	case err := <-errChan:
		log.Fatal(err)
		// add context that we can close? (signal??)
		//case <-r.Context().Done():
		//	close(untilTime)
	}
	return
}
