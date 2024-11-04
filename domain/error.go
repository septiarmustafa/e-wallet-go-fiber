package domain

import "errors"

var ErrAuthFailed = errors.New("Error authentication failed")
var ErrUsernameTaken = errors.New("Error username already taken")
var ErrOtpInvalid = errors.New("Error OTP invalid")
