
%%%
title = "prom2ssh 8"
area = "System Administration"
workgroup = "Prometheus"
%%%

# NAME

prom2ssh - export failed logins to prometheus

# SYNOPSIS

**dnsfmt** [**OPTIONS**]...

# DESCRIPTION

**Prom2ssh** parsed the journald of ssh and extract failed login attemps. It exports two metrics

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
