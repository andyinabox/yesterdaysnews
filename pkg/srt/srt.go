package srt

import (
	"fmt"
	"sync"
	"time"
)

const TimeFormat = "15:04:05,000"
const subTmpl = "%d\n%s --> %s\n%s\n\n"

type Caption struct {
	Start    time.Time
	Duration time.Duration
	Text     string
}

func (c Caption) End() time.Time {
	return c.Start.Add(c.Duration)
}

type SRT struct {
	Captions []Caption
	mu       sync.Mutex
}

func New() *SRT {
	return &SRT{
		Captions: []Caption{},
	}
}

func (s *SRT) Add(start time.Time, duration time.Duration, text string) {
	s.mu.Lock()
	s.Captions = append(s.Captions, Caption{
		Start:    start,
		Duration: duration,
		Text:     text,
	})
	s.mu.Unlock()
}

func (s *SRT) AddToEnd(duration time.Duration, text string) {

	var start time.Time

	if len(s.Captions) != 0 {
		start = s.Captions[len(s.Captions)-1].End()
	}

	s.Add(start, duration, text)

}

func (s *SRT) String() (out string) {

	var start, end string
	for i, c := range s.Captions {
		start = c.Start.Format(TimeFormat)
		end = c.End().Format(TimeFormat)
		out = out + fmt.Sprintf(subTmpl, i, start, end, c.Text)
	}

	return
}
