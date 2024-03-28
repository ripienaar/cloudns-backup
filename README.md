## ClouDNS Backup

A basic tool to export all master zones in a ClouDNS account to Bind zone files.

## Status?

It's working in the most basic form, later we might support restore also.

## Usage?

Run `cloudns-backup backup --target /path/to/backup` for authentication pass `--auth-id` and `--password` which can 
also be set using `CLOUDNS_AUTH_ID` and `CLOUDNS_AUTH_PASSWORD` environment variables.

A docker container is available at [ghcr.io/ripienaar/cloudns-backup](https://ghcr.io/ripienaar/cloudns-backup)

## Contact?

R.I. Pienaar / rip@devco.net / [devco.net](https://www.devco.net/)
