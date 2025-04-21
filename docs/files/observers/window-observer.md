# Window Observer

The **Window Observer** is a TimelyGator component that monitors the active window on the user’s machine and sends periodic heartbeats to the TimelyGator server with metadata about the current application, window title, URL, and incognito state.

## Overview

- **Cross-platform**: Supports Linux (X11), macOS, and Windows.
- **Polling-based**: Checks the active window at a configurable interval.
- **Anonymization**: Optionally exclude or regex-mask window titles.
- **Swift helper (macOS)**: Fallback to a native Swift helper for improved performance.

## Prerequisites

- Go 1.17+ installed.
- [TimelyGator server](../backend.md) running and accessible.
  - To run go to `cd server` and run `go run .`
- Platform-specific dependencies:
  - **Linux**: X11 development libraries (`libx11-dev`, `xprop`).
  - **macOS**: JXA or AppleScript support.
  - **Windows**: Windows API access via `golang.org/x/sys/windows`.

## Configuration

The observer reads default settings from environment variables or a `.env` file (via [godotenv](https://github.com/joho/godotenv)). You can also override via CLI flags.

### Environment Variables / `.env`

- `HOST` – TimelyGator server host (default: `localhost`).
- `PORT` – TimelyGator server port (default: `8080`).
- `TESTING` – Enable testing mode (default: `false`).
- `VERBOSE` – Enable verbose logging (default: `false`).
- `EXCLUDE_TITLE` – Always replace titles with `excluded` (default: `false`).
- `EXCLUDE_TITLES` – Comma-separated regex list to anonymize only matching titles.
- `POLL_TIME` – Polling interval in seconds (default: `1`).
- `STRATEGY` – macOS strategy: `jxa`, `applescript`, or `swift`.

## CLI Usage

```bash
cd server/observers/window-observer

go run ./main
  --host YOUR_HOST \
  --port YOUR_PORT \
  --testing=false \
  --verbose \
  --exclude-title=false \
  --exclude-titles "^Secret.*,Private.*" \
  --poll-time 2.5 \
  --strategy jxa
```

| Flag             | Description                                             | Default            |
|------------------|---------------------------------------------------------|--------------------|
| `--host`         | TimelyGator server host                                 | `cfg.Host`         |
| `--port`         | TimelyGator server port                                 | `cfg.Port`         |
| `--testing`      | Testing mode: skips real server write                   | `cfg.Testing`      |
| `--verbose`      | Turn on verbose (microseconds) logging                  | `cfg.Verbose`      |
| `--exclude-title`| Replace *every* title with `excluded`                   | `cfg.ExcludeTitle` |
| `--exclude-titles`| Regex list for titles to anonymize                     | `cfg.ExcludeTitles` |
| `--poll-time`    | Polling interval (seconds)                              | `cfg.PollTime`     |
| `--strategy`     | macOS only: JXA, AppleScript, or Swift helper selection | `cfg.Strategy`     |

## How It Works

1. **Startup**: Load `.env`, parse flags, ensure macOS permissions if needed.
2. **Connect**: Create a TimelyGator client, wait for server readiness, create a `currentwindow` bucket.
3. **Swift Helper** *(macOS / `swift` strategy)*:
   - Launches a bundled Swift binary for continuous monitoring.
   - Sends its own heartbeats via the helper and exits Go loop.
4. **Polling Loop**:
   - Every `poll-time` seconds, call `lib.GetCurrentWindow(strategy)`.
   - Optionally anonymize titles via `--exclude-title` or `--exclude-titles` regexes.
   - Marshal window data and send a heartbeat with `timestamp`, `duration=0`, and the JSON payload.
   - Continue until interrupted by SIGINT / SIGTERM.

## Development & Testing

- **Unit Tests**: See `main_test.go` and `lib_common_test.go`. Run:

  ```bash
  # in window-observer directory
  go test ./lib
  go test
  ```

- **Platform Stubs**: `lib_darwin_stub.go`, `lib_windows_stub.go` allow compilation on non-target OS.
- **Adding Mocks**: Override `getCurrentWindow` in tests to simulate various window states.

## Troubleshooting

- **Permission Denied (macOS)**: Ensure Accessibility permissions are granted to the observer binary. This is a very common issue on macOS.
- **X11 Connection Errors (Linux)**: Confirm `$DISPLAY` is set and X server is running.
- **Empty Window Info**: If title or app is blank, falls back to `"unknown"`.

---

*This documentation is part of the TimelyGator observability suite.*
