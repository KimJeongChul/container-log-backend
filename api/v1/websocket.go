package v1

import (
	"time"
)

// Timetout Websocket
const (
	pingPeriod     = 30 * time.Second
	pongWait       = 30 * time.Second
	writeWait      = 30 * time.Second
	maxMessageSize = 512
)
