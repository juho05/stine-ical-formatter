# STiNE iCal Formatter

ICS files acquired from the STiNE scheduler export tool seem to be corrupted and can't be imported into calendar clients.

This tool fixes the format of these files so they can be successfully imported.

## Run locally

### Prerequisites

- [Go](https://go.dev)
- [Git](https://git-scm.com/)

```bash
git clone https://github.com/juho05/stine-ical-formatter
cd stine-ical-formatter
go run .
```

Open http://localhost:8080.

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