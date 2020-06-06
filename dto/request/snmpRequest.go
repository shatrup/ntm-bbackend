package request

import (
	"fmt"
	"ntm-backend/errs"
	"time"
)

type SnmpRequest struct {
	NasId				 string `json:"nas_id"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

func (s SnmpRequest) IsValidNasId() bool {
	if s.NasId != "" {
		return true
	} else {
		return false
	}
}

func (s SnmpRequest) IsValidStartEndDate() bool {
	layout := "2006-01-02 15:04:05"
	sd, err := time.Parse(layout, s.StartDate)
	if err != nil {
		return false
	}
	ed, err := time.Parse(layout, s.EndDate)
	if err != nil {
		return false
	}
	if sd.Before(ed) {
		return false
	}
	return true
}

func (s SnmpRequest) isValidDate(date string) bool {
	layout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Asia/Kolkata")
	ist, err := time.ParseInLocation(layout, date, loc)
	if ist.IsZero() {
		fmt.Printf("No date has been set, %s\n", ist)
		return false
	}
	if err != nil {
		return false
	}
	if !ist.Before(time.Now()) {
		return false
	}
	return true
}

func (s SnmpRequest) IsValidEndDate() bool {
	return s.isValidDate(s.EndDate)
}

func (s SnmpRequest) IsValidStartDate() bool {
	return s.isValidDate(s.StartDate)
}

func (s SnmpRequest) Validate() errs.Errs {
	errors := &errs.Errors{}
	if !s.IsValidStartDate() {
		errors.Add(errs.InvalidStartDate)
		return errors
	}
	if !s.IsValidEndDate() {
		errors.Add(errs.InvalidEndDate)
		return errors
	}
	if len(errors.Errors) == 0 {
		if s.IsValidStartEndDate() {
			errors.Add(errs.InvalidStartEndDate)
		}
	}
	if !s.IsValidNasId() {
		errors.Add(errs.NasIdIsEmpty)
	}
	if len(errors.Errors) == 0 {
		return nil
	}
	return errors
}
