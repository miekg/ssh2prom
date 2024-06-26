%%%
title = "ssh2prom 8"
area = "System Administration"
workgroup = "Prometheus"
%%%

# NAME

ssh2prom - export failed logins to prometheus

# SYNOPSIS

**ssh2prom** [**OPTIONS**]...

# DESCRIPTION

**Ssh2prom** parsed the journald of ssh and extract failed login attemps. It exports two metrics

- ssh_failed_total{family="1|2"}: all failed logins
- ssh_failed_root_total{family="1|2"}: failed logins for root only
- ssh_sucess_total{}: all successful logins

If family is 1 it is an IPv4 connection, for 2 it is coming over IPv6.

# OPTIONS

`-a` **ADDR**
: Start the prometheus server on *ADDR*

`-d`
: Enable debugging, show the logs and parsed users and address family

`-n`
: Dry run, do everything except exporting the metrics

`-u` **UNIT**
: Use unit **UNIT** instead of "ssh"

# AUTHOR

Miek Gieben <miek@miek.nl>.
