package auth

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
var ErrSmthWentWrong = errors.New("something went wrong")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidUsername = errors.New("invalid username")
