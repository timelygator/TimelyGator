package listener

import (
    hook "github.com/robotn/gohook"
)

// StartAllListeners begins a global event loop that routes OS-level input events
// to the provided keyboard and mouse listeners.
func StartAllListeners(kbListener *KeyboardListener, msListener *MouseListener) {
    eventChan := hook.Start()
    go func() {
        defer hook.End()

        for event := range eventChan {
            switch event.Kind {

            // Keyboard events
            case hook.KeyDown:
                if kbListener != nil {
                    kbListener.onPress(rune(event.Keychar))
                }
            case hook.KeyUp:
                if kbListener != nil {
                    kbListener.onRelease(rune(event.Keychar))
                }

            // Mouse events
            case hook.MouseDown:
                if msListener != nil {
                    msListener.onClick(int(event.X), int(event.Y), int(event.Button), true)
                }
            case hook.MouseUp:
                if msListener != nil {
                    msListener.onClick(int(event.X), int(event.Y), int(event.Button), false)
                }
            case hook.MouseMove, hook.MouseDrag:
                if msListener != nil {
                    msListener.onMove(int(event.X), int(event.Y))
                }
            case hook.MouseWheel:
                if msListener != nil {
                    // Often event.Button is the direction, event.Clicks is the magnitude
                    msListener.onScroll(int(event.X), int(event.Y), int(event.Button), int(event.Clicks))
                }
            }
        }
    }()
}

