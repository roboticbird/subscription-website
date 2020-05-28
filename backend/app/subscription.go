package app

import (
	"time"
)

type Subscription struct {
	ID        int       `json:"id"`
	Product   Product   `json:"product"`
	StartDate time.Time `json:"start-date"`
	Status    string    `json:"status"`
}
