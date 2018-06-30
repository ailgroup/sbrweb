package app

import (
	"errors"
	"fmt"
	"time"
)

const (
	GuestGTE      = 4
	GuestLTE      = 1
	HotelIDsLTE   = 5
	timeShortForm = "01-02"
)

var (
	ErrHotelIDNullOrZero    = errors.New("invalid hotel_id: query not present or value not defined")
	ErrGuestCountNullOrZero = errors.New("invalid guest_count: query not present or value is 0")

	ErrSearchCriterion = errors.New("invalid search criterion: verify only one kind of criterion (hotel_id, lat_lng, city_code)")
)

var (
	ErrGuestLTEMsg    = "invalid guest_count '%d': must be less than or equal '%d'"
	ErrGuestGTEMsg    = "invalid guest_count '%d': must be greater than or equal '%d'"
	ErrHotelIDsLTEMsg = "invalid hotel_ids '%d': must be less than or equal '%d'"
	ErrArriveFmtMsg   = "invalid arrive '%v': format (MM-DD '%s'). %v"
	ErrDepartFmtMsg   = "invalid depart '%v': format (MM-DD '%s'). %v"
	ErrStayRangeMsg   = "invalid range: arrive '%v' must be before depart '%v'"
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
func Lt(current, compare int) bool {
	return current < compare
}
func Gt(current, compare int) bool {
	return current > compare
}
func startBeforeEnd(start, end time.Time) bool {
	return start.Before(end)
}
func ValidArriveDepart(ts string) (time.Time, bool, error) {
	t, err := time.Parse(timeShortForm, ts)
	if err != nil {
		return t, false, err
	}
	return t, true, nil
}
