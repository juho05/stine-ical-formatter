package formatter

import (
	"fmt"
	"io"
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
	data := []byte(calendar.Serialize())
	log.Tracef("formatted %d files in %s resulting in %d bytes", len(files), time.Since(start).String(), len(data))
	return data, nil
}
