package domain

import "time"

type Date struct {
	Date string
}

func (d Date) UTCDate() (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	locIST, _ := time.LoadLocation("Asia/Kolkata")
	locUTC, _ := time.LoadLocation("UTC")
	sd, err := time.ParseInLocation(layout, d.Date, locIST)
	if err != nil {
		return time.Time{}, err
	}
	return sd.In(locUTC), nil
}

func (d Date) ISTDate() (string, error) {
	layout := "2006-01-02 15:04:05"
	locIST, _ := time.LoadLocation("Asia/Kolkata")
	sd, err := time.ParseInLocation(layout, d.Date, locIST)
	if err != nil {
		return "", err
	}
	return sd.Format(layout), nil
}
