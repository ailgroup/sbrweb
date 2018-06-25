package app

import (
	"errors"
	"net/http"
)

type InputValidator interface {
	Validate(r *http.Request) error
}

var (
	ErrInvalidArrive          = errors.New("invalid arrive")
	ErrInvalidDepart          = errors.New("invalid depart")
	ErrInvalidGuestCount      = errors.New("invalid guest_count")
	ErrInvalidSearchCriterion = errors.New("invalid search criterion, verify only one kind of criterion (hotel_id, lat_lng, city_code) was sent")
	ErrInvalidHotelID         = errors.New("invalid hotel_id")
	ErrInvalidCityCode        = errors.New("invalid city_code")
)
