package domain

import "errors"

var ErrAuthFailed = errors.New("error authentication failed")
var ErrUsernameTaken = errors.New("error username already taken")
var ErrOtpInvalid = errors.New("error OTP invalid")
var ErrAccountNotFound = errors.New("error account not found")
var ErrInquiryNotFound = errors.New("error inquiry not found")
var ErrInsufficientBalance = errors.New("error insufficient balance")
