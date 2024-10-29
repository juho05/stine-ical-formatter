# STiNE iCal Formatter

ICS files acquired from the STiNE scheduler export tool seem to be corrupted and can't be imported into calendar clients.

This tool fixes the format of these files so they can be successfully imported.

Additionally, it detects recurring events and converts them to actual iCal recurrences.

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