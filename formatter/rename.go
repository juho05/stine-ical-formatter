package formatter

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode"

	ics "github.com/arran4/golang-ical"
)

func Names(calendar []byte) ([]string, error) {
	cal, err := ics.ParseCalendar(bytes.NewReader(calendar))
	if err != nil {
		return nil, fmt.Errorf("parse calendar: %w", err)
	}
	return eventNames(cal), nil
}

func Rename(calendar []byte, names []string) ([]byte, error) {
	cal, err := ics.ParseCalendar(bytes.NewReader(calendar))
	if err != nil {
		return nil, fmt.Errorf("parse calendar: %w", err)
	}
	original := eventNames(cal)
	if len(names) != len(original) {
		return nil, fmt.Errorf("expected %d names, got %d", len(original), len(names))
	}
	renamed := make(map[string]string, len(names))
	for i, n := range names {
		n = cleanName(n)
		if n != "" {
			renamed[original[i]] = n
		}
	}
	for _, e := range cal.Events() {
		summary := e.GetProperty(ics.ComponentPropertySummary)
		if summary == nil {
			continue
		}
		if n, ok := renamed[summary.Value]; ok {
			summary.Value = n
		}
	}
	return []byte(cal.Serialize(ics.WithNewLineWindows)), nil
}

func eventNames(cal *ics.Calendar) []string {
	names := make([]string, 0, 10)
	for _, e := range cal.Events() {
		summary := e.GetProperty(ics.ComponentPropertySummary)
		if summary == nil || summary.Value == "" || slices.Contains(names, summary.Value) {
			continue
		}
		names = append(names, summary.Value)
	}
	slices.Sort(names)
	return names
}

func cleanName(name string) string {
	name = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name)
	return strings.TrimSpace(name)
}
