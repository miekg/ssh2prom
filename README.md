# ssh2prom

Annoyingly simple daemon that tails the ssh journal and exports failed root login attempts.

Exports the following metrics:

- `ssh_failed_total{family="1|2"}`: all failed attempts (including root)
- `ssh_failed_root_total{family="1|2"}`: failed attempts where the user is root
- `ssh_sucess_total{}`: all successful logins

The total also includes the root user, the root_total is the subset that counts root only.
If family is 1 it's a connection over IPv4, with 2 it's IPv6.

Every occurence of `Failed password for .... from .... port ... ssh2` is counted as a failed
attempt.
Every `session opened for user ` is counted as a success.

It tails the log using the 'fuck you' (-fu) option of journalctl.
