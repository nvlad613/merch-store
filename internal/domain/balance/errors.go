package balance

import "errors"

var (
	NotEnoughCoinsError = errors.New("operation failed: not enough coins")
)
