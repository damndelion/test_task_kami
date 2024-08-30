package customrErrors

import (
	"errors"
)

var ErrAlreadyBooked = errors.New("time slot already booked, please try different time")
var ErrEndBeforeStart = errors.New("start time must be before end time")
