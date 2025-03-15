package formatter

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/juho05/log"
)

func Format(files []io.Reader) ([]byte, error) {
	calendar := &ics.Calendar{
		Components:         []ics.Component{},
		CalendarProperties: []ics.CalendarProperty{},
	}
	start := time.Now()
	for i, f := range files {
		cal, err := ics.ParseCalendar(newReader(f))
		if err != nil {
			return nil, fmt.Errorf("invalid file content: %w", err)
		}
		if i == 0 {
			calendar.CalendarProperties = append(calendar.CalendarProperties, cal.CalendarProperties...)
		}
		components := make([]ics.Component, 0, len(cal.Components))
		for _, c := range cal.Components {
			if _, ok := c.(*ics.VTimezone); !ok || i == 0 {
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

func combineEvents(cal *ics.Calendar) error {
	events := cal.Events()

	distinct := make(map[string][]*ics.VEvent, 10)
	for _, e := range events {
		startAtProp := strings.Split(e.GetProperty(ics.ComponentPropertyDtStart).Value, ":")
		e.SetProperty(ics.ComponentPropertyDtStart, startAtProp[len(startAtProp)-1])
		endAtProp := strings.Split(e.GetProperty(ics.ComponentPropertyDtEnd).Value, ":")
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
		key := strings.Join([]string{
			e.GetProperty(ics.ComponentPropertySummary).Value,
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
			excluded := int(start.Sub(previous).Hours()/24)/7 - 1
			for x := range excluded {
				exclude := previous.AddDate(0, 0, 7*(x+1))
				d[0].AddExdate(exclude.Format("20060102T150405"), &ics.KeyValues{
					Key:   "TZID",
					Value: []string{"CampusNetZeit"},
				})
			}
			previous = start

			cal.RemoveEvent(e.Id())
		}

		if len(d) > 1 {
			d[0].AddRrule(fmt.Sprintf("FREQ=WEEKLY;UNTIL=%s", d[len(d)-1].GetProperty(ics.ComponentPropertyDtStart).Value))
		}
		d[0].GetProperty(ics.ComponentPropertyDtStart).ICalParameters["TZID"] = []string{"CampusNetZeit"}
		d[0].GetProperty(ics.ComponentPropertyDtEnd).ICalParameters["TZID"] = []string{"CampusNetZeit"}
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
