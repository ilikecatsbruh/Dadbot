# DadBot

A Go desktop deal finder for second-hand CDs, vinyl, DVDs, Blu-rays, cassettes and music collections in Athens and across Greece.

## Current foundation

- Wails v2 desktop shell with a dark glass-style dashboard
- Go scan engine with a 20-goroutine worker pool
- Starter Greek/English search targets for Athens and Greece
- Target validation, scan cancellation, status reporting and bounded in-memory logs

## Planned modules

- Source adapters for permitted listing APIs, feeds and user-provided links
- SQLite listing cache and duplicate detection
- Targets manager with location, radius, category and price filters
- Listings browser with deal scoring and shipping details
- Searchable cache and log viewers
- Optional notifications

## Access and rate limits

DadBot will not bypass logins, CAPTCHAs, paywalls, access controls or platform rate limits. Facebook Marketplace support will use user-provided listing links or other authorized access methods. Each source adapter must enforce its own request limits.

## Development

Requirements: Go 1.23+ and the Wails v2 prerequisites for your operating system.

```sh
go mod tidy
go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.2
wails dev
```

The current scan engine is a foundation only: it processes configured targets but does not yet fetch marketplace listings.
