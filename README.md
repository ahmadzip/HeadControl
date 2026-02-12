# HeadControl

Minimal Headscale admin dashboard, built with Go and HTMX.

*[Baca dalam Bahasa Indonesia](README.id.md) — translated from Indonesian by an AI, if something sounds weird, now you know why.*

This project was made for personal use. I chose Go and HTMX as the foundation
because I needed something truly lightweight and fast no heavy JavaScript
frameworks, no complex build pipelines, just a single binary that serves
everything. If you find it useful or want to modify it for your own setup,
feel free to do so.

*Disclaimer: I am not particularly good at Go, so if you read the code and wonder
"why is it done this way?", the answer is most likely "because that is what worked."*

---

## Stack

- **Go** — HTTP server, API client, template rendering, SQLite storage
- **HTMX** — partial page updates without writing JavaScript
- **Lucide Icons** — icon set loaded via CDN
- **SQLite** — local configuration storage

## Features

- Connect to any Headscale instance via API key
- Dashboard with node/user statistics
- User management (create, rename, delete)
- Node management (rename, expire, delete, tags, routes)
- Node detail view
- Multiple color themes
- Responsive layout (desktop, tablet, mobile)

---

## Requirements

- Go 1.21 or later
- GCC (required by go-sqlite3, see note below)
- A running Headscale server with an API key

### GCC on Windows

The `go-sqlite3` driver requires CGO. On Windows, install
[MSYS2](https://www.msys2.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/),
then make sure `gcc` is available in your PATH.

---

## Setup

Clone the repository:

```
git clone https://github.com/ahmadzip/headcontrol.git
cd headcontrol
```

Install dependencies:

```
go mod tidy
```

Build:

```
go build .
```

This produces `headcontrol.exe` on Windows or `headcontrol` on Linux/macOS.

Run:

```
./headcontrol
```

The server starts on `http://localhost:8080` by default.

### Command-line flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `8080` | Server port |
| `-db` | `headcontrol.db` | SQLite database path |

Example:

```
./headcontrol -port 3000 -db /data/headcontrol.db
```

---

## First Run

1. Open `http://localhost:8080` in your browser.
2. You will be redirected to the setup page.
3. Enter your Headscale server URL (e.g. `https://headscale.example.com`).
4. Enter your API key (created with `headscale apikeys create`).
5. Click "Test Connection" to verify.
6. Click "Save" to proceed to the dashboard.

---

## Development

For hot-reload during development, install [Air](https://github.com/air-verse/air):

```
go install github.com/air-verse/air@latest
```

Then run:

```
air
```

Air watches for file changes and rebuilds automatically.
The configuration is in `.air.toml`.

---

## Project Structure

```
headcontrol/
  main.go                          entrypoint, route registration
  internal/
    handler/
      handler.go                   core struct, template engine, middleware
      helpers.go                   render helpers, time formatting
      setup.go                     setup page handlers
      dashboard.go                 dashboard page handlers
      users.go                     user management handlers
      nodes.go                     node management handlers
      settings.go                  settings page handlers
    headscale/
      client.go                    headscale API client
    model/
      models.go                    data structures
    store/
      store.go                     SQLite storage layer
  templates/
    layout/layout.html             base layout with sidebar
    pages/                         full page templates
    partials/                      HTMX partial templates
  static/
    css/app.css                    main stylesheet
    css/theme/                     color theme files
    js/app.js                      client-side logic
```

---

## Themes

HeadControl ships with 16 built-in themes. Switch between them using
the theme selector in the top navigation bar. Your selection is stored
in localStorage.

---

## Feedback

Built during break time. Suggestions, criticism, or pull requests are welcome.

---

## License

Do whatever you want with it.
