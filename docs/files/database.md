# TimelyGator Database

This document describes the internal schema, models, and API for the TimelyGator database system, which stores time-based event data. It uses **GORM** over **SQLite** as the backend.

---

## Overview

TimelyGator stores data in two primary models:

- **Bucket**: Metadata container that groups events.
- **Event**: A timestamped JSON object representing an observation or action.

These models are defined in `models.go`, and the operational interface is built in `database.go`.

---

## Models

### `Event`
Represents a single atomic record (e.g., a window switch, browser tab, keystroke).

```go
ID        uint           // Primary key
BucketID  string         // Foreign key to Bucket
Timestamp time.Time      // UTC timestamp of event
Duration  float64        // Duration in seconds
Data      datatypes.JSON // Payload (e.g., {app, title, url, incognito})
```

Utility methods:
- `ToJSONDict()` – expands the event into a merged JSON map.
- `DataEqualEvent(other)` – compares two events by JSON contents.
- `DataEqualJSON(json)` – compares against a raw JSON blob.

### `Bucket`
Defines a collection of time-ordered events and related metadata.

```go
ID       string         // Primary key
Name     *string        // Optional label
Type     string         // Application type (e.g., currentwindow)
Client   string         // Originating client name
Hostname string         // Device hostname
Created  time.Time      // When the bucket was created
Data     datatypes.JSON // Optional metadata
```

---

## Initialization

The function `InitDB(cfg Config) (*Datastore, error)` initializes the database:

- Locates the `data/` directory.
- Opens or creates a SQLite database file.
- Calls `AutoMigrate` for `Bucket` and `Event`.

Returns a `*Datastore`, the primary handle.

---

## Core API

### `Datastore` (global manager)

```go
type Datastore struct {
    db *gorm.DB
}
```

#### Methods
- `DB() *gorm.DB` – exposes raw GORM handle
- `Buckets() map[string]map[string]interface{}` – returns metadata for all buckets
- `CreateBucket(...)` – inserts a new bucket
- `UpdateBucket(bucketID, updates)` – modifies fields in a bucket
- `DeleteBucket(bucketID)` – removes the bucket
- `GetBucket(bucketID)` – returns a `*Bucket` object if found

### `Bucket` (logical event group)

Returned by `GetBucket()` or `NewBucket()`, provides full access to a bucket’s events.

```go
type Bucket struct {
    ds       *Datastore
    bucketID string
}
```

#### Bucket API
- `Metadata()` – returns the bucket’s current fields
- `Get(limit, start, end)` – fetches N events, optionally filtered by time range
- `GetByID(eventID)` – fetches an event by primary key
- `Insert(event[s])` – insert single or batch of events
- `Delete(eventID)` – remove an event
- `Replace(eventID, event)` – update an event in-place
- `ReplaceLast(event)` – update the most recent event
- `GetEventCount(start, end)` – return the total number of events

---

## Event Insertion

Use the helper:
```go
NewEvent(id, timestamp, duration, data)
```
This handles:
- Accepting `time.Time`, `string`, or `nil` for timestamp
- Accepting `time.Duration` or `float64` for duration
- Marshaling map data into `datatypes.JSON`

Then insert:
```go
bucket.Insert(event)
```

---

## JSON Utilities

JSON fields are stored using `gorm.io/datatypes.JSON`. Several helper functions and methods exist to ensure round-trip safety, shallow equality, and proper merging with metadata.

---

## Notes

- All timestamps are stored in **UTC**, truncated to millisecond precision.
- The `created` field in `Bucket` is the bucket creation time, not the first event.
- Events are **not automatically expired** — cleanup policies are managed externally.
- The system currently uses SQLite, but GORM allows switching to Postgres or MySQL.

---

*This document is part of the TimelyGator backend system documentation.*

