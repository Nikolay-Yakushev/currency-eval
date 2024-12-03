package models

import "time"

type Currency struct {
	Name  string
	Value float64
	Date  time.Time
}
