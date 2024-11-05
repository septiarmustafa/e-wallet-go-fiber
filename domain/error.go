package domain

import "errors"

var ErrAuthFailed = errors.New("Error authentication failed")
var ErrUsernameTaken = errors.New("Error username already taken")
var ErrOtpInvalid = errors.New("Error OTP invalid")
var ErrAccountNotFound = errors.New("Error account not found")
var ErrInquiryNotFound = errors.New("Error inquiry not found")
var ErrInsufficientBalance = errors.New("Error insufficient balance")
