package entity

import (
	"time"
)

const (
	TimeFormat = "2006-01-02"
)

type UserDetails struct {
	Id               int
	UserId           int
	Phone            string
	Gender           int
	TypeOfDisability int
	Address          string
	Birthdate        Birthdate
	Image            string
	Description      string
}

type Birthdate interface {
	GetDOB(year, month, day int) time.Time
}

func GetDOB(year, month, day int) time.Time {
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dob
}
