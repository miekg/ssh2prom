package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	failedRootLogins = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ssh_failed_root_total",
		Help: "Counter of failed root logins.",
	}, []string{"family"})
	failedUserLogins = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ssh_failed_total",
		Help: "Counter of total failed logins.",
	}, []string{"family"})
	userLogins = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ssh_success_total",
		Help: "Counter of total successful logins.",
	})
)
