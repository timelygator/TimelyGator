package main

import (
	"testing"
	"time"

	"gorm.io/datatypes"
)

func TestParseDateFlag(t *testing.T) {
	validDateStr := "2024-04-01"
	invalidDateStr := "2024-13-01"

	timeParsed, err := parseDateFlag(validDateStr)
	if err != nil {
		t.Errorf("parseDateFlag failed on valid input: %v", err)
	}
	if timeParsed.Format("2006-01-02") != validDateStr {
		t.Errorf("Expected %s, got %s", validDateStr, timeParsed.Format("2006-01-02"))
	}

	_, err = parseDateFlag(invalidDateStr)
	if err == nil {
		t.Errorf("Expected error for invalid date string, got none")
	}
}

func TestSameDay(t *testing.T) {
	date1 := time.Date(2024, 3, 31, 12, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC)
	date3 := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	if !sameDay(date1, date2) {
		t.Errorf("sameDay failed, expected true for same dates")
	}
	if sameDay(date1, date3) {
		t.Errorf("sameDay failed, expected false for different dates")
	}
}

func TestWeightedChoice(t *testing.T) {
	items := []sampleData{
		{App: "test1", Weight: 0},
		{App: "test2", Weight: 1},
		{App: "test3", Weight: 99},
	}

	counts := map[string]int{"test1": 0, "test2": 0, "test3": 0}
	iterations := 1000
	for i := 0; i < iterations; i++ {
		choice := weightedChoice(items)
		counts[choice.App]++
	}

	if counts["test1"] != 0 {
		t.Errorf("Expected 'test1' to never be chosen, but got %d occurrences", counts["test1"])
	}
	if counts["test3"] <= counts["test2"] {
		t.Errorf("Expected 'test3' to be chosen significantly more than 'test2', got test3: %d, test2: %d", counts["test3"], counts["test2"])
	}
}

func TestPickDuration(t *testing.T) {
	minutes := 10.0
	maxSecs := 120.0

	for i := 0; i < 100; i++ {
		dur := pickDuration(minutes, maxSecs)
		if dur < 5*time.Minute || dur > 20*time.Minute {
			t.Errorf("Duration out of expected range (5-20 mins): got %v", dur)
		}
	}

	// Test with zero minutes, random duration
	for i := 0; i < 100; i++ {
		dur := pickDuration(0, maxSecs)
		if dur < 5*time.Second || dur > 120*time.Second {
			t.Errorf("Random duration out of expected range (5-120 secs): got %v", dur)
		}
	}
}

func TestGetString(t *testing.T) {
	dataJSON := datatypes.JSON(`{"app": "testapp", "title": "testtitle"}`)
	app := getString(dataJSON, "app")
	title := getString(dataJSON, "title")
	missing := getString(dataJSON, "missing")

	if app != "testapp" {
		t.Errorf("Expected 'testapp', got %s", app)
	}

	if title != "testtitle" {
		t.Errorf("Expected 'testtitle', got %s", title)
	}

	if missing != "" {
		t.Errorf("Expected empty string for missing key, got %s", missing)
	}
}
