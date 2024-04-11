# ssh2prom

Annoyingly simple daemon that tails the ssh journal and exports failed root login attempts.

Exports the following metrics:

- ssh_failed_user_total{}
- ssh_failed_root_total{}

The user_total also includes the root user, the root_total is the subset the counts root only.

Every occurence of `Failed password for .... from .... port ... ssh2` is counted.

It tails the log using the 'fuck you' (-fu) option of journalctl.
