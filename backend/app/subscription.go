package app

import (
	"time"
)

type Subscription struct {
	ID         int       `json:"id"`
	Product    Product   `json:"product"`
	StartDate  time.Time `json:"startDate"`
	Status     string    `json:"subStatus"`
	NextStatus string    `json:"nextStatus"`
}
