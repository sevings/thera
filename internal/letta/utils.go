package letta

import (
	"fmt"
	"strings"
	"time"
)

// Custom time unmarshaler to handle different time formats
type FlexTime struct {
	time.Time
}

func (ft *FlexTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = strings.Trim(s, "\"")

	formats := []string{
		time.RFC3339,           // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05",  // Without timezone
		"2006-01-02T15:04:05Z", // With Z
		time.RFC1123,
		time.RFC1123Z,
		time.UnixDate,
		time.RubyDate,
	}

	var err error
	for _, format := range formats {
		var t time.Time
		t, err = time.Parse(format, s)
		if err == nil {
			ft.Time = t
			return nil
		}
	}

	return fmt.Errorf("unable to parse time: %s", s)
}
