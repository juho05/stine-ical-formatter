package formatter

import (
	"errors"
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
)

const inputTimezone = "CampusNetZeit"

var (
	inputStandard = map[string]string{
		"DTSTART":      "16011028T030000",
		"RRULE":        "FREQ=YEARLY;BYDAY=-1SU;BYMONTH=10",
		"TZOFFSETFROM": "+0200",
		"TZOFFSETTO":   "+0100",
	}
	inputDaylight = map[string]string{
		"DTSTART":      "16010325T020000",
		"RRULE":        "FREQ=YEARLY;BYDAY=-1SU;BYMONTH=3",
		"TZOFFSETFROM": "+0100",
		"TZOFFSETTO":   "+0200",
	}
	unsupportedEventProperties = []ics.ComponentProperty{
		ics.ComponentPropertyRrule,
		ics.ComponentPropertyRdate,
		ics.ComponentPropertyExdate,
		ics.ComponentPropertyExrule,
		ics.ComponentProperty(ics.PropertyRecurrenceId),
		ics.ComponentProperty(ics.PropertyDuration),
	}
)

// Event times are relabeled as Europe/Berlin without conversion, which is only correct for the known STiNE format.
func validate(cal *ics.Calendar) error {
	timezones := cal.Timezones()
	for _, tz := range timezones {
		if err := validateTimezone(tz); err != nil {
			return fmt.Errorf("timezone: %w", err)
		}
	}
	for _, e := range cal.Events() {
		if len(timezones) == 0 {
			return fmt.Errorf("event %s: file does not define timezone %s", e.Id(), inputTimezone)
		}
		if err := validateEvent(e); err != nil {
			return fmt.Errorf("event %s: %w", e.Id(), err)
		}
	}
	return nil
}

func validateTimezone(tz *ics.VTimezone) error {
	if !hasExactProperties(tz.Properties, map[string]string{"TZID": inputTimezone}) {
		return fmt.Errorf("expected only TZID:%s", inputTimezone)
	}
	var standard, daylight int
	for _, c := range tz.Components {
		switch c := c.(type) {
		case *ics.Standard:
			standard++
			if !hasExactProperties(c.Properties, inputStandard) {
				return errors.New("unexpected STANDARD definition")
			}
		case *ics.Daylight:
			daylight++
			if !hasExactProperties(c.Properties, inputDaylight) {
				return errors.New("unexpected DAYLIGHT definition")
			}
		default:
			return fmt.Errorf("unexpected component %T", c)
		}
	}
	if standard != 1 || daylight != 1 {
		return errors.New("expected exactly one STANDARD and one DAYLIGHT definition")
	}
	return nil
}

func validateEvent(e *ics.VEvent) error {
	for _, name := range []ics.ComponentProperty{ics.ComponentPropertyDtStart, ics.ComponentPropertyDtEnd} {
		prop := e.GetProperty(name)
		if prop == nil {
			return fmt.Errorf("no %s", name)
		}
		tzid := prop.ICalParameters["TZID"]
		if len(prop.ICalParameters) != 1 || len(tzid) != 1 || tzid[0] != inputTimezone {
			return fmt.Errorf("%s must have TZID=%s as its only parameter", name, inputTimezone)
		}
		if _, err := time.Parse(localTimeFormat, prop.Value); err != nil {
			return fmt.Errorf("%s is not a local date time: %w", name, err)
		}
	}
	for _, name := range unsupportedEventProperties {
		if e.GetProperty(name) != nil {
			return fmt.Errorf("unsupported property %s", name)
		}
	}
	return nil
}

func hasExactProperties(props []ics.IANAProperty, want map[string]string) bool {
	if len(props) != len(want) {
		return false
	}
	seen := make(map[string]bool, len(props))
	for _, p := range props {
		value, ok := want[p.IANAToken]
		if !ok || seen[p.IANAToken] || value != p.Value || len(p.ICalParameters) != 0 {
			return false
		}
		seen[p.IANAToken] = true
	}
	return true
}
