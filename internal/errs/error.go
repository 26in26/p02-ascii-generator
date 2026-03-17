package errs

import "errors"

var (
	InvalidDimensions    = errors.New("invalid dimensions")
	InvalidAspectRatio   = errors.New("invalid aspect ratio")
	BufferNotInitialized = errors.New("buffer not initialized")
)
