package listener

import (
	"reflect"
	"testing"
	"time"

	hook "github.com/robotn/gohook"
)

func TestKeyboardListener_OnPressAndNextEvent(t *testing.T) {
	kl := NewKeyboardListener()
	// Initially, event data should have "presses" set to 0.
	eventData := kl.NextEvent()
	if presses, ok := eventData["presses"]; ok && presses.(int) != 0 {
		t.Errorf("Expected initial presses to be 0, got %v", presses)
	}
	// Simulate key presses.
	kl.onPress('a')
	kl.onPress('b')
	if !kl.HasNewEvent() {
		t.Errorf("Expected new event flag to be true after onPress")
	}
	eventData = kl.NextEvent()
	if presses, ok := eventData["presses"]; !ok || presses.(int) != 2 {
		t.Errorf("Expected presses to be 2, got %v", eventData["presses"])
	}
	if kl.HasNewEvent() {
		t.Errorf("Expected new event flag to be false after NextEvent")
	}
}

func TestMouseListener_OnMove_OnClick_OnScroll(t *testing.T) {
	ml := NewMouseListener()
	// Check initial values.
	eventData := ml.NextEvent()
	if clicks, ok := eventData["clicks"].(int); !ok || clicks != 0 {
		t.Errorf("Expected initial clicks to be 0, got %v", eventData["clicks"])
	}
	// Simulate mouse movement.
	ml.onMove(10, 10) // sets initial position
	ml.onMove(20, 30)
	eventData = ml.NextEvent()
	// Expected: deltaX = |20-10| = 10, deltaY = |30-10| = 20.
	if dx, ok := eventData["deltaX"].(int); !ok || dx != 10 {
		t.Errorf("Expected deltaX to be 10, got %v", eventData["deltaX"])
	}
	if dy, ok := eventData["deltaY"].(int); !ok || dy != 20 {
		t.Errorf("Expected deltaY to be 20, got %v", eventData["deltaY"])
	}

	// Simulate mouse click.
	ml.onClick(20, 30, 1, true)
	eventData = ml.NextEvent()
	if clicks, ok := eventData["clicks"].(int); !ok || clicks != 1 {
		t.Errorf("Expected clicks to be 1, got %v", eventData["clicks"])
	}

	// Simulate mouse scroll.
	ml.onScroll(0, 0, 5, 3)
	eventData = ml.NextEvent()
	if scrollX, ok := eventData["scrollX"].(int); !ok || scrollX != 5 {
		t.Errorf("Expected scrollX to be 5, got %v", eventData["scrollX"])
	}
	if scrollY, ok := eventData["scrollY"].(int); !ok || scrollY != 3 {
		t.Errorf("Expected scrollY to be 3, got %v", eventData["scrollY"])
	}
}

func TestMouseListener_OnMove_Accumulation(t *testing.T) {
	ml := NewMouseListener()
	// First move initializes position.
	ml.onMove(0, 0)
	// Then move to (10,10) and then to (20,5).
	ml.onMove(10, 10)
	ml.onMove(20, 5)
	eventData := ml.NextEvent()
	// Calculations:
	// First move: initial position set, no delta.
	// Second move: dx = |10-0| = 10, dy = |10-0| = 10.
	// Third move: dx += |20-10| = 10, dy += |5-10| = 5.
	// So, deltaX = 20, deltaY = 15.
	if dx, ok := eventData["deltaX"].(int); !ok || dx != 20 {
		t.Errorf("Expected deltaX to be 20, got %v", eventData["deltaX"])
	}
	if dy, ok := eventData["deltaY"].(int); !ok || dy != 15 {
		t.Errorf("Expected deltaY to be 15, got %v", eventData["deltaY"])
	}
}

func TestStartAllListeners_NilListeners(t *testing.T) {
	// Call StartAllListeners with nil listeners. It should not panic.
	go StartAllListeners(nil, nil)
	// Allow a short duration for the event loop to start.
	time.Sleep(100 * time.Millisecond)
	hook.End() // End the event loop.
}

func TestBaseEventFactory_NextEventReset(t *testing.T) {
	// Create a baseEventFactory instance and set some dummy data.
	bef := newBaseEventFactory()
	bef.eventData["test"] = "value"
	bef.newEvent = true

	data := bef.NextEvent()
	// Ensure that the returned data matches what was set.
	if !reflect.DeepEqual(data, map[string]interface{}{"test": "value"}) {
		t.Errorf("Expected event data to be {test: value}, got %v", data)
	}
	// After NextEvent, the internal map should be reset.
	if len(bef.eventData) != 0 {
		t.Errorf("Expected event data to be reset, got %v", bef.eventData)
	}
	if bef.newEvent {
		t.Errorf("Expected newEvent flag to be false after NextEvent")
	}
}
