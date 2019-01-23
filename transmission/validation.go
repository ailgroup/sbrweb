package transmission

import (
	"errors"
	"fmt"
	"time"
)

const (
	GuestMax      = 4
	GuestMin      = 1
	HotelIDsMax   = 5
	timeShortForm = "2006-01-02"
)

var (
	ErrBaseParamsNull       = errors.New("invalid params: arrive, depart, guest_count not present")
	ErrArriveNull           = errors.New("invalid arrive: query not present or value not defined")
	ErrDepartNull           = errors.New("invalid depart: query not present or value not defined")
	ErrHotelIDNullOrZero    = errors.New("invalid hotel_ids: query not present or value not defined")
	ErrGuestCountNullOrZero = errors.New("invalid guest_count: query not present or value is 0")

	ErrSearchCriterion = errors.New("invalid search criterion: verify only one kind of criterion (hotel_id, lat_lng, city_code)")
	ErrGuestMaxMsg     = "invalid guest_count '%d': must be less than or equal '%d'"
	ErrGuestMinMsg     = "invalid guest_count '%d': must be greater than or equal '%d'"
	ErrHotelIDsMaxsg   = "invalid hotel_ids '%d': must be less than or equal '%d'"
	ErrArriveFmtMsg    = "invalid arrive '%s': format (YYYY-MM-DD '%s'). %s"
	ErrStayInPastMsg   = "invalid arrive '%s': cannot be before today '%s'"
	ErrDepartFmtMsg    = "invalid depart '%s': format (YYYY-MM-DD '%s'). %s"
	ErrStayRangeMsg    = "invalid range: depart '%s' must be after arrive '%s'"
)

func ErrLtGt(msgFormat string, given, expect int) error {
	return fmt.Errorf(msgFormat, given, expect)
}
func ErrStayFormat(msgFormat, errMsg, given, expect string) error {
	return fmt.Errorf(msgFormat, given, expect, errMsg)
}
func ErrStayRange(msgFormat, given, expect string) error {
	return fmt.Errorf(msgFormat, given, expect)
}
func ErrArriveInPast(msgFormat, given string, now time.Time) error {
	return fmt.Errorf(msgFormat, given, now)
}
func Lt(current, compare int) bool {
	return current < compare
}
func Gt(current, compare int) bool {
	return current > compare
}
func DepartBeforeArrive(d, a time.Time) bool {
	return d.Before(a)
}
func ArriveNotInPast(a, n time.Time) bool {
	return a.Before(n)
}

// StayFormat validates format of incoming params, returning boolean for conditional checking, error for creating an informative validation error message, and parsed time in the app time zone location in order to check other before/after validations.
func StayFormat(stay string, loc *time.Location) (time.Time, bool, error) {
	t, err := time.ParseInLocation(timeShortForm, stay, loc)
	if err != nil {
		return t, false, err
	}
	return t, true, nil
}

// BeginOfDay sets the beginning of the day for time
func BeginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
