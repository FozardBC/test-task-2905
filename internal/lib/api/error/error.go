package api

import "errors"

var (
	ErrReqBodyDecode  = errors.New("failed to decode request body")
	ErrEncodeResponse = errors.New("failed to encode response body")
)
