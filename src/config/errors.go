package config

import (
	"errors"
)

var (
	ErrStartChargeVoltageMoreStop     = errors.New("start charge voltage should be less than stop charge voltage")
	ErrStartChargeVoltage             = errors.New("start charge voltage should be less than 8.4V")
	ErrStopChargeVoltage              = errors.New("stop charge voltage should not be more than 8.4V")
	ErrStartChargeVoltageLessShutdown = errors.New("start charge voltage should be more than shutdown woltage")
	ErrShutdownVoltage                = errors.New("shutdown voltage level should be more than 6.1V")
)
