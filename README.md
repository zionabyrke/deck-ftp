# DECK
## Deploy Every Commit, Kontinuous

A small CLI for git commits and deploys changed files to hosts that have no built-in CI/CD, currently applied on InfinityFree using FTP

No more dragging files into FileZilla by hand

## Why

Free/shared hosts like InfinityFree give you FTP and no CI, no deploy hooks, no webhooks. DECK fills that gap: commit on `main`, DECK diffs your local files against what's already live, and only uploads what changed

See [`ARCHITECTURE.md`](./ARCHITECTURE.md) for the UMLs

## Status

Prototype v1. FTP only. Single host, single target. Local git hook trigger (no cloud CI yet — see roadmap in the architecture doc).

## Install

```bash
git clone https://github.com/zionabyrke/deck-ftp
cd deck-ftp
go build .
```

## Setup

**1. Config**: `deck-ftp init` creates a starter `deck.yaml`:

```yaml
local_dir: site
remote_dir: htdocs
```

**2. Credentials**: create a `.env` file (never commit this):

```
DECK_FTP_HOST=ftpupload.net:21
DECK_FTP_USER=epiz_xxxxx
DECK_FTP_PASS=yourpass
```

## Usage

```bash
deck-ftp diff   # show what would change, no network calls
deck-ftp push   # deploy the diff
deck-ftp status # config + pending changes at a glance
deck-ftp install-hook # auto-run push on every commit to main
```

Once `install-hook` is set up, just commit normally:

```bash
git add .
git commit -m "update homepage"
# DECK diffs and deploys automatically
```

## How it decides what to upload

Every file in `local_dir` is hashed (sha256). The hash set is compared against the last successful deploy's manifest (`.deck/manifest.json`, local, gitignored). Only new, changed, or removed files touch the network

## License

[MIT LICENSE](./LICENSE)

## Contributing

Adapters beyond FTP (SFTP, cPanel API, S3) are the most useful thing to add next. Please see the `Adapter` interface in the architecture doc
