# STiNE iCal Formatter

ICS files acquired from the STiNE scheduler export tool seem to be corrupted and can't be imported into calendar clients.

This tool fixes the format of these files so they can be successfully imported.

Additionally, it detects recurring events and merges them together using iCal `RRULE` and `EXDATE`.

## Run locally

### Prerequisites

- [Go](https://go.dev)
- [Git](https://git-scm.com/)
- [Node.js](https://nodejs.org) with npm
- [Make](https://www.gnu.org/software/make/)

```bash
git clone https://github.com/juho05/stine-ical-formatter
cd stine-ical-formatter
make init
make run
```

Open http://localhost:8080. Set the `PORT` environment variable to use a different port.

`make init` downloads the Go modules, installs the npm packages and builds `web/static/css/tailwind.css`. It only needs to be run once.

### Development

```bash
make tailwind-watch
make go-watch
```

Run both in separate terminals. The first rebuilds the stylesheet when a template changes, the second restarts the server when Go files, templates or static files change.

### Build

```bash
make
```

Writes the binary to `bin/stine-ical-formatter`.

## Deploy with Docker

Create `docker-compose.yml`:
```yaml
services:
  stine-ical-formatter:
    image: ghcr.io/juho05/stine-ical-formatter
    restart: unless-stopped
    environment:
      # set to true if stine-ical-formatter is behind a reverse proxy that sets the X-Forwarded-For header
      RATE_LIMIT_X_FORWARDED_FOR: false
    ports:
      - "8080:8080"
```

Run `docker compose up -d` in the same directory.
