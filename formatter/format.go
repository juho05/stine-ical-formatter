package formatter

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strings"
	"time"
	_ "time/tzdata"

	ics "github.com/arran4/golang-ical"
	"github.com/juho05/log"
)

const (
	timezone        = "Europe/Berlin"
	localTimeFormat = "20060102T150405"
	utcTimeFormat   = "20060102T150405Z"
)

func Format(files []io.Reader) ([]byte, error) {
	calendar := &ics.Calendar{
		Components:         []ics.Component{newTimezone()},
		CalendarProperties: []ics.CalendarProperty{},
	}
	start := time.Now()
	for i, f := range files {
		cal, err := ics.ParseCalendar(newReader(f))
		if err != nil {
			return nil, fmt.Errorf("invalid file content: %w", err)
		}
		err = validate(cal)
		if err != nil {
			return nil, fmt.Errorf("unexpected file content: %w", err)
		}
		if i == 0 {
			calendar.CalendarProperties = append(calendar.CalendarProperties, cal.CalendarProperties...)
		}
		components := make([]ics.Component, 0, len(cal.Components))
		for _, c := range cal.Components {
			if _, ok := c.(*ics.VTimezone); !ok {
				components = append(components, c)
			}
		}
		calendar.Components = append(calendar.Components, components...)
	}
	err := combineEvents(calendar)
	if err != nil {
		return nil, fmt.Errorf("combine events: %w", err)
	}

	data := []byte(calendar.Serialize())
	log.Tracef("formatted %d files in %s resulting in %d bytes", len(files), time.Since(start).String(), len(data))
	return data, nil
}

func newTimezone() *ics.VTimezone {
	tz := ics.NewTimezone(timezone)
	standard := &ics.Standard{}
	standard.AddProperty(ics.ComponentPropertyDtStart, "19701025T030000")
	standard.AddProperty(ics.ComponentPropertyRrule, "FREQ=YEARLY;BYDAY=-1SU;BYMONTH=10")
	standard.AddProperty(ics.ComponentProperty(ics.PropertyTzoffsetfrom), "+0200")
	standard.AddProperty(ics.ComponentProperty(ics.PropertyTzoffsetto), "+0100")
	daylight := &ics.Daylight{}
	daylight.AddProperty(ics.ComponentPropertyDtStart, "19700329T020000")
	daylight.AddProperty(ics.ComponentPropertyRrule, "FREQ=YEARLY;BYDAY=-1SU;BYMONTH=3")
	daylight.AddProperty(ics.ComponentProperty(ics.PropertyTzoffsetfrom), "+0100")
	daylight.AddProperty(ics.ComponentProperty(ics.PropertyTzoffsetto), "+0200")
	tz.Components = append(tz.Components, standard, daylight)
	return tz
}

func combineEvents(cal *ics.Calendar) error {
	events := cal.Events()
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("load timezone: %w", err)
	}

	distinct := make(map[string][]*ics.VEvent, 10)
	for _, e := range events {
		dtStart := e.GetProperty(ics.ComponentPropertyDtStart)
		if dtStart == nil {
			return fmt.Errorf("event %s has no DTSTART", e.Id())
		}
		dtEnd := e.GetProperty(ics.ComponentPropertyDtEnd)
		if dtEnd == nil {
			return fmt.Errorf("event %s has no DTEND", e.Id())
		}
		startAtProp := strings.Split(dtStart.Value, ":")
		e.SetProperty(ics.ComponentPropertyDtStart, startAtProp[len(startAtProp)-1])
		endAtProp := strings.Split(dtEnd.Value, ":")
		e.SetProperty(ics.ComponentPropertyDtEnd, endAtProp[len(endAtProp)-1])

		startAt, err := e.GetStartAt()
		if err != nil {
			return fmt.Errorf("get start at: %w", err)
		}
		endAt, err := e.GetEndAt()
		if err != nil {
			return fmt.Errorf("get end at: %w", err)
		}
		startH, startM, startS := startAt.Clock()
		endH, endM, endS := endAt.Clock()
		var summary string
		if s := e.GetProperty(ics.ComponentPropertySummary); s != nil {
			summary = s.Value
		}
		var location string
		if l := e.GetProperty(ics.ComponentPropertyLocation); l != nil {
			location = l.Value
		}
		key := strings.Join([]string{
			summary,
			location,
			startAt.Weekday().String(),
			fmt.Sprintf("%d:%d:%d", startH, startM, startS),
			endAt.Weekday().String(),
			fmt.Sprintf("%d:%d:%d", endH, endM, endS),
		}, "\t")

		if distinct[key] == nil {
			distinct[key] = make([]*ics.VEvent, 0, 1)
		}
		distinct[key] = append(distinct[key], e)
	}

	for _, d := range distinct {
		slices.SortFunc(d, func(a *ics.VEvent, b *ics.VEvent) int {
			aStartAt := mustGetStartAt(a)
			bStartAt := mustGetStartAt(b)
			return aStartAt.Compare(bStartAt)
		})

		previous := mustGetStartAt(d[0])
		for i, e := range d {
			if i == 0 {
				continue
			}
			start := mustGetStartAt(e)
			// a clock change makes the gap an hour shorter or longer
			excluded := int(math.Round(start.Sub(previous).Hours()/24))/7 - 1
			for x := range excluded {
				exclude := previous.AddDate(0, 0, 7*(x+1))
				d[0].AddExdate(exclude.Format(localTimeFormat), &ics.KeyValues{
					Key:   "TZID",
					Value: []string{timezone},
				})
			}
			previous = start

			cal.RemoveEvent(e.Id())
		}

		if len(d) > 1 {
			until, err := time.ParseInLocation(localTimeFormat, d[len(d)-1].GetProperty(ics.ComponentPropertyDtStart).Value, loc)
			if err != nil {
				return fmt.Errorf("parse last start: %w", err)
			}
			d[0].AddRrule(fmt.Sprintf("FREQ=WEEKLY;UNTIL=%s", until.UTC().Format(utcTimeFormat)))
		}
		d[0].GetProperty(ics.ComponentPropertyDtStart).ICalParameters["TZID"] = []string{timezone}
		d[0].GetProperty(ics.ComponentPropertyDtEnd).ICalParameters["TZID"] = []string{timezone}
	}

	return nil
}

func mustGetStartAt(e *ics.VEvent) time.Time {
	start, err := e.GetStartAt()
	if err != nil {
		panic(err)
	}
	return start
}
