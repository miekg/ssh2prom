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

- ssh_failed_total{}
- ssh_failed_root_total{}

That can be scraped by prometheus.

No semantic checks are done, this is purely text manipulation with some basic zone file syntax
understanding.

# OPTIONS

`-a` **ADDR**
: Start the prometheus server on *ADDR*

`-d`
: Enable debugging, show the logs and parsed users

`-n`
: Dry run, do everything except export the metrics

`-u` **UNIT**
: Use unit **UNIT** instead of "ssh"

# AUTHOR

Miek Gieben <miek@miek.nl>.
