## ClouDNS Backup

A basic tool to export all master zones in a ClouDNS account to Bind zone files.

## Status?

It's working in the most basic form, later we might support restore also.

## Usage?

Run `cloudns-backup backup --target /path/to/backup` for authentication pass `--auth-id` and `--password` which can 
also be set using `CLOUDNS_AUTH_ID` and `CLOUDNS_AUTH_PASSWORD` environment variables.

A docker container is available at [ghcr.io/ripienaar/cloudns-backup](https://ghcr.io/ripienaar/cloudns-backup)

I use this from a Gitea action, to create daily backups of my zones in git, as below - should be similar for GitHub.

```yaml
name: CI

on:
  repository_dispatch: {}
  workflow_dispatch: {}
  push:
    branches:
      - main

  schedule:
    - cron: '@daily'

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Code Checkout
        id: checkout
        uses: actions/checkout@v4

      - name: Install Dependencies
        id: dependencies
        env:
          DEBIAN_FRONTEND: noninteractive

        run: |
          apt-get update || true
          apt-get -y install rake || true
          apt-get -y install docker.io || true

      - name: Setup Directories
        id: setup
        run: |
          mkdir -p zones

      - name: Get Zones
        id: download
        run: |
          RUNNER="zones_runner_${GITHUB_RUN_NUMBER}"

          git checkout "${GITHUB_REF:11}"

          docker run --name ${RUNNER} ghcr.io/ripienaar/cloudns-backup:latest \
            backup \
            --target /zones \
            --auth-id "${{ secrets.CLOUDNS_AUTH_ID }}" \
            --password "${{ secrets.CLOUDNS_AUTH_PASSWORD }}"
          docker cp ${RUNNER}:/zones `pwd`
          docker rm ${RUNNER}

      - name: Push Zones
        id: push
        run: |
          if [ -z "$(git status --porcelain)" ]; then
            git status
            echo "Nothing to commit"
          else
            echo ">>> Differences that will be stored"
            echo
            git diff
            echo

            git config --global user.name 'Actions Backup'
            git config --global user.email 'root@example.net'
            git remote set-url origin https://x-access-token:${{ secrets.GITHUB_TOKEN }}@git.example.net/$GITHUB_REPOSITORY
            git add zones
            git commit -am "ClouDNS Backup Job ${GITHUB_RUN_NUMBER}"
            git push
          fi
```

## Contact?

R.I. Pienaar / rip@devco.net / [devco.net](https://www.devco.net/)
