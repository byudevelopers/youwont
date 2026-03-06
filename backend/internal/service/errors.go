package service

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrForbidden          = errors.New("forbidden")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInsufficientPoints = errors.New("insufficient points")
	ErrBetNotOpen         = errors.New("bet is not open")
	ErrAlreadyWagered     = errors.New("already wagered on this bet")
	ErrDeciderCannotWager = errors.New("decider cannot wager")
	ErrCannotSelfDecide   = errors.New("creator cannot be decider")
	ErrNoOpposingSide     = errors.New("cannot resolve with no opposing wagers")
)
