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

