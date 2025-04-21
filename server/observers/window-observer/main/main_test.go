//go:build linux || darwin || windows
// +build linux darwin windows

package main

import (
    "encoding/json"
    "regexp"
    "testing"
    "time"

    "timelygator/server/database/models"
    "timelygator/server/observers/window-observer/lib"
)

// --- fake heartbeatSender ---
type fakeTG struct {
    bucket string
    data   interface{}
    pulse  float64
    called bool
}

func (f *fakeTG) Heartbeat(bucket string, data interface{}, pulse float64, wait bool, extra *float64) error {
    f.bucket = bucket
    f.data = data
    f.pulse = pulse
    f.called = true
    return nil
}

// stub getCurrentWindow
func mockWin(app, title string, url *string, inc *bool) func(string) (*lib.WindowInfo, error) {
    return func(_ string) (*lib.WindowInfo, error) {
        return &lib.WindowInfo{App: app, Title: title, URL: url, Incognito: inc}, nil
    }
}

func TestHandleHeartBeat_Basic(t *testing.T) {
    url := "https://foo"
    inc := true

    // override getCurrentWindow
    oldGW := getCurrentWindow
    getCurrentWindow = mockWin("MyApp", "MyTitle", &url, &inc)
    defer func() { getCurrentWindow = oldGW }()

    ft := &fakeTG{}
    poll := 3 * time.Second

    handleHeartBeat(ft, "buck", "any", false, nil, poll)

    if !ft.called {
        t.Fatal("expected Heartbeat to be called")
    }
    if ft.bucket != "buck" {
        t.Errorf("bucket = %q; want %q", ft.bucket, "buck")
    }

    // we know we passed a *models.Event in data
    ev, ok := ft.data.(*models.Event)
    if !ok {
        t.Fatalf("expected data to be *models.Event, got %T", ft.data)
    }
    var payload map[string]interface{}
    if err := json.Unmarshal(ev.Data, &payload); err != nil {
        t.Fatalf("unmarshal event.Data: %v", err)
    }
    if payload["app"] != "MyApp" || payload["title"] != "MyTitle" {
        t.Errorf("payload = %+v; want app=MyApp,title=MyTitle", payload)
    }

    wantPulse := poll.Seconds() + 1
    if ft.pulse != wantPulse {
        t.Errorf("pulse = %v; want %v", ft.pulse, wantPulse)
    }
}

func TestHandleHeartBeat_RegexAnonymize(t *testing.T) {
    oldGW := getCurrentWindow
    getCurrentWindow = mockWin("X", "Secret", nil, nil)
    defer func() { getCurrentWindow = oldGW }()

    ft := &fakeTG{}
    patterns := []*regexp.Regexp{regexp.MustCompile("(?i)secret")}
    handleHeartBeat(ft, "b", "any", false, patterns, time.Second)

    ev := ft.data.(*models.Event)
    var payload map[string]interface{}
    _ = json.Unmarshal(ev.Data, &payload)
    if payload["title"] != "excluded" {
        t.Errorf("title = %v; want excluded", payload["title"])
    }
}

func TestHandleHeartBeat_ForceExclude(t *testing.T) {
    oldGW := getCurrentWindow
    getCurrentWindow = mockWin("X", "Whatever", nil, nil)
    defer func() { getCurrentWindow = oldGW }()

    ft := &fakeTG{}
    handleHeartBeat(ft, "b", "any", true, nil, time.Second)

    ev := ft.data.(*models.Event)
    var payload map[string]interface{}
    _ = json.Unmarshal(ev.Data, &payload)
    if payload["title"] != "excluded" {
        t.Errorf("title = %v; want excluded", payload["title"])
    }
}

func TestInterrupt(t *testing.T) {
    ch := interrupt()
    if ch == nil {
        t.Fatal("interrupt returned nil")
    }
}

func TestHostPortPtr(t *testing.T) {
    s := "hello"
    if got := hostPtr(&s); got != "hello" {
        t.Errorf("hostPtr = %q; want hello", got)
    }
    i := 42
    if got := portPtr(&i); got != "42" {
        t.Errorf("portPtr = %q; want 42", got)
    }
}
