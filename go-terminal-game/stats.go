package main

import "time"

type stats struct {
	start  time.Time
	frames int
	fps    float64
}

func newStats() *stats {
	return &stats{
		start: time.Now(),
	}
}

func (s *stats) update() {
	s.frames++
	if s.frames == MAX_SAMPLES {
		s.fps = float64(s.frames) / time.Since(s.start).Seconds()
		s.frames = 0
	}
}
