# AFK Observer

The **AFK Observer** is a TimelyGator component that detects user activity/inactivity (AFK status) by monitoring global keyboard and mouse input. It sends heartbeats to the TimelyGator server when a user becomes AFK or returns.

---

## Overview

- **Platform**: Linux only (via [gohook](https://github.com/robotn/gohook))
- **Input Sources**: Global mouse and keyboard listeners
- **Heartbeat Mechanism**: Sends periodic status events with `afk` or `not-afk`
- **Command-line Configurable**: Uses `cobra` for flexible CLI interface

---

## CLI Usage

```bash
cd server/observers/afk-observer

./go run ./main \
  --host localhost \
  --port 8080 \
  --timeout 300s \
  --poll-time 10s \
  --verbose
```

| Flag         | Description                              | Default   |
|--------------|------------------------------------------|-----------|
| `--host`     | TimelyGator server host                  | `localhost` |
| `--port`     | Server port                              | `8080`    |
| `--timeout`  | Seconds of inactivity before AFK         | `180s`    |
| `--poll-time`| Frequency of polling user input          | `5s`      |
| `--verbose`  | Enable verbose logging                   | `false`   |
| `--testing`  | Enable testing mode (no real uploads)    | `false`   |

---

## How It Works

1. **Startup**: Initializes the TimelyGator client, checks platform compatibility (Linux-only).
2. **Event Loop**: Captures keyboard and mouse input events using `gohook`.
3. **Activity Tracking**:
   - MouseListener tracks:
     - Movement (deltaX/Y)
     - Clicks
     - Scrolls
   - KeyboardListener tracks:
     - Key presses
   - These listeners update a shared `lastActivity` timestamp.
4. **Heartbeat Loop**:
   - Runs every `poll-time`.
   - Checks time since last input.
   - Sends heartbeats when state changes (AFK → active or vice versa), or periodically to confirm status.

---

## Bucket Format

The observer uses a single bucket per device:
```
afkstatus
```

Event example:
```json
{
  "timestamp": "2025-04-21T16:03:12Z",
  "duration": 0,
  "data": {
    "status": "not-afk"
  }
}
```

---

## Components

### `AFKWatcher`
Located in `afk.go`, manages:
- Client startup
- Bucket creation
- Heartbeat generation logic

### `listener` Package
Implements:
- `MouseListener` — motion, click, scroll tracking
- `KeyboardListener` — key press tracking
- `baseEventFactory` — thread-safe shared structure for tracking state
- `StartAllListeners()` — OS-level event router using `gohook`

### `unix.go`
- Singleton `LastInputUnix` monitors event timestamps
- Provides `SecondsSinceLastInput()` for polling loop

---

## Development & Testing

- Run with `--verbose` to log AFK state transitions
- Simulate user inactivity by avoiding all input
- Uses `os.Getppid()` to exit if parent process dies
- `testing` flag disables real server writes (useful for dry-run/testing)

---

## Limitations

- Currently supports **Linux only**
- No explicit debounce; frequent input toggles AFK quickly
- Scroll events may vary across devices/environments

---

*This documentation is part of the TimelyGator observability suite.*
