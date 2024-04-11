package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	failedRootLogins = promauto.NewCounter(prometheus.GaugeOpts{
		Name: "ssh_failed_root_total",
		Help: "Counter of failed root logins.",
	})
	failedUserLogins = promauto.NewCounter(prometheus.GaugeOpts{
		Name: "ssh_failed_total",
		Help: "Counter of total failed logins.",
	})
)
