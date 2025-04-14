package listener

import (
    "sync"
)

// baseEventFactory provides the shared fields and methods for listeners.
type baseEventFactory struct {
    mu        sync.Mutex
    eventData map[string]interface{}
    newEvent  bool
}

func newBaseEventFactory() *baseEventFactory {
    return &baseEventFactory{
        eventData: make(map[string]interface{}),
        newEvent:  false,
    }
}

// NextEvent returns the current aggregated event data and resets
// the internal state so it can begin building a new event.
func (bef *baseEventFactory) NextEvent() map[string]interface{} {
    bef.mu.Lock()
    defer bef.mu.Unlock()

    data := bef.eventData
    bef.newEvent = false
    bef.resetData()
    return data
}

func (bef *baseEventFactory) HasNewEvent() bool {
    bef.mu.Lock()
    defer bef.mu.Unlock()

    return bef.newEvent
}

// resetData can be overridden by derived types to restore default fields.
func (bef *baseEventFactory) resetData() {
    bef.eventData = make(map[string]interface{})
}

// KeyboardListener listens for keyboard presses.
type KeyboardListener struct {
    *baseEventFactory
}

// NewKeyboardListener creates a KeyboardListener and sets default fields.
func NewKeyboardListener() *KeyboardListener {
    kl := &KeyboardListener{
        baseEventFactory: newBaseEventFactory(),
    }
    kl.resetData()
    return kl
}

// resetData implements baseEventFactory's default fields for keyboard events.
func (kl *KeyboardListener) resetData() {
    kl.baseEventFactory.resetData()
    kl.eventData["presses"] = 0
}

// onPress is called by the global event loop to record a key press.
func (kl *KeyboardListener) onPress(key rune) {
    kl.mu.Lock()
    defer kl.mu.Unlock()

    oldVal, ok := kl.eventData["presses"].(int)
    if !ok {
        oldVal = 0
    }
    kl.eventData["presses"] = oldVal + 1

    kl.newEvent = true
}


// onRelease is called by the global event loop when a key is released.
func (kl *KeyboardListener) onRelease(key rune) {
    // Leaving for now
}

// MouseListener listens for mouse movements, clicks, and scrolls.
type MouseListener struct {
    *baseEventFactory
    pos *position
}

type position struct {
    x int
    y int
}

// NewMouseListener creates a MouseListener with the default fields set.
func NewMouseListener() *MouseListener {
    ml := &MouseListener{
        baseEventFactory: newBaseEventFactory(),
        pos:              nil,
    }
    ml.resetData()
    return ml
}

// resetData implements baseEventFactory's default fields for mouse events.
func (ml *MouseListener) resetData() {
    ml.baseEventFactory.resetData()
    ml.eventData["clicks"] = 0
    ml.eventData["deltaX"] = 0
    ml.eventData["deltaY"] = 0
    ml.eventData["scrollX"] = 0
    ml.eventData["scrollY"] = 0
}

// onMove is called by the global event loop for mouse-move events.
func (ml *MouseListener) onMove(x, y int) {
    ml.mu.Lock()
    defer ml.mu.Unlock()

    if ml.pos == nil {
        ml.pos = &position{x: x, y: y}
    } else {
        dx := abs(ml.pos.x - x)
        dy := abs(ml.pos.y - y)

        oldDeltaX, ok := ml.eventData["deltaX"].(int)
        if !ok {
            oldDeltaX = 0
        }
        oldDeltaY, ok := ml.eventData["deltaY"].(int)
        if !ok {
            oldDeltaY = 0
        }

        ml.eventData["deltaX"] = oldDeltaX + dx
        ml.eventData["deltaY"] = oldDeltaY + dy

        ml.pos.x = x
        ml.pos.y = y
    }
    ml.newEvent = true
}

// onClick is called by the global event loop for mouse-click events.
func (ml *MouseListener) onClick(x, y int, button int, down bool) {
    if down {
        ml.mu.Lock()
        defer ml.mu.Unlock()

        // If "clicks" has never been set or was cleared, initialize it to 0
        oldVal, ok := ml.eventData["clicks"].(int)
        if !ok {
            oldVal = 0
        }
        ml.eventData["clicks"] = oldVal + 1

        ml.newEvent = true
    }
}

// onScroll is called by the global event loop for scroll wheel events.
func (ml *MouseListener) onScroll(x, y int, scrollX, scrollY int) {
    ml.mu.Lock()
    defer ml.mu.Unlock()

    oldScrollX, ok := ml.eventData["scrollX"].(int)
    if !ok {
        oldScrollX = 0
    }
    oldScrollY, ok := ml.eventData["scrollY"].(int)
    if !ok {
        oldScrollY = 0
    }

    ml.eventData["scrollX"] = oldScrollX + abs(scrollX)
    ml.eventData["scrollY"] = oldScrollY + abs(scrollY)
    ml.newEvent = true
}
func abs(i int) int {
    if i < 0 {
        return -i
    }
    return i
}
