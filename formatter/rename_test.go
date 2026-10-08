package formatter

import (
	"slices"
	"strings"
	"testing"
)

const renameCalendar = "BEGIN:VCALENDAR\r\n" +
	"VERSION:2.0\r\n" +
	"BEGIN:VEVENT\r\n" +
	"UID:1\r\n" +
	"SUMMARY:Vorlesung B\\, Teil 1\r\n" +
	"END:VEVENT\r\n" +
	"BEGIN:VEVENT\r\n" +
	"UID:2\r\n" +
	"SUMMARY:Seminar A\r\n" +
	"END:VEVENT\r\n" +
	"BEGIN:VEVENT\r\n" +
	"UID:3\r\n" +
	"SUMMARY:Vorlesung B\\, Teil 1\r\n" +
	"END:VEVENT\r\n" +
	"END:VCALENDAR\r\n"

func TestNames(t *testing.T) {
	names, err := Names([]byte(renameCalendar))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Seminar A", "Vorlesung B, Teil 1"}
	if !slices.Equal(names, want) {
		t.Fatalf("got %q, want %q", names, want)
	}
}

func TestRename(t *testing.T) {
	out, err := Rename([]byte(renameCalendar), []string{"  ", "ML; VL\r\n"})
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if strings.Count(got, "SUMMARY:ML\\; VL\r\n") != 2 {
		t.Errorf("expected both events to be renamed:\n%s", got)
	}
	if strings.Count(got, "SUMMARY:Seminar A\r\n") != 1 {
		t.Errorf("expected blank name to keep the original:\n%s", got)
	}
}

func TestRenameWithoutChanges(t *testing.T) {
	out, err := Rename([]byte(renameCalendar), []string{"", ""})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != renameCalendar {
		t.Errorf("expected unchanged calendar:\n%s", out)
	}
}

func TestRenameWrongCount(t *testing.T) {
	if _, err := Rename([]byte(renameCalendar), []string{"A"}); err == nil {
		t.Fatal("expected error")
	}
}
