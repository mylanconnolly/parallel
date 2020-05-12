package main

import "time"

// Ctx is used to represent command context.
type Ctx struct {
	Cmd   string
	Input string
	Start time.Time
	Time  time.Time
}
