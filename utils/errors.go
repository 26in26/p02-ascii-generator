package utils

import "errors"

var (
	ErrInvalidDimensions    = errors.New("invalid dimensions")
	ErrInvalidAspectRatio   = errors.New("invalid aspect ratio")
	ErrEmptyFrame           = errors.New("empty frame")
	ErrBufferNotInitialized = errors.New("buffer not initialized")
)
