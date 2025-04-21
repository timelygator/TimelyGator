package afkobserver

import (
	"log"
	"sync"
	"time"

	"timelygator/server/observers/afk-observer/listener"
)

// LastInputUnix tracks the last time an input event was detected.
type LastInputUnix struct {
	mouseListener    *listener.MouseListener
	keyboardListener *listener.KeyboardListener
	lastActivity     time.Time
	mu               sync.Mutex
}

// NewLastInputUnix creates and initializes a new LastInputUnix instance.
func NewLastInputUnix() *LastInputUnix {
	// Initialize mouse and keyboard listeners
	mouseListener := listener.NewMouseListener()
	keyboardListener := listener.NewKeyboardListener()

	// Start capturing input events
	listener.StartAllListeners(keyboardListener, mouseListener)

	return &LastInputUnix{
		mouseListener:    mouseListener,
		keyboardListener: keyboardListener,
		lastActivity:     time.Now(),
	}
}

// SecondsSinceLastInput returns the number of seconds since the last input event.
func (li *LastInputUnix) SecondsSinceLastInput() float64 {
	li.mu.Lock()
	defer li.mu.Unlock()

	now := time.Now()
	if li.mouseListener.HasNewEvent() || li.keyboardListener.HasNewEvent() {
		log.Println("[LastInputUnix] New input event detected.")
		li.lastActivity = now

		// Get/clear events
		li.mouseListener.NextEvent()
		li.keyboardListener.NextEvent()
	}
	return now.Sub(li.lastActivity).Seconds()
}

var lastInputInstance *LastInputUnix
var once sync.Once

// SecondsSinceLastInput is a singleton-style function to get seconds since the last input.
func SecondsSinceLastInput() float64 {
	once.Do(func() {
		lastInputInstance = NewLastInputUnix()
	})
	return lastInputInstance.SecondsSinceLastInput()
}

// func Execute() {
// 	for {
// 		time.Sleep(1 * time.Second)
// 		elapsed := SecondsSinceLastInput()
// 		log.Printf("Seconds since last input: %.2f\n", elapsed)
// 	}
// }
