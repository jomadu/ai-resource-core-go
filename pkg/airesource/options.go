package airesource

import "time"

type LoadOptions struct {
	MaxFileSize     int64
	MaxArraySize    int
	MaxNestingDepth int
	Timeout         time.Duration
}

func DefaultLoadOptions() LoadOptions {
	return LoadOptions{
		MaxFileSize:     10 * 1024 * 1024, // 10MB
		MaxArraySize:    10000,
		MaxNestingDepth: 100,
		Timeout:         30 * time.Second,
	}
}

type LoadOption func(*LoadOptions)

func WithMaxFileSize(size int64) LoadOption {
	return func(o *LoadOptions) {
		o.MaxFileSize = size
	}
}

func WithMaxArraySize(size int) LoadOption {
	return func(o *LoadOptions) {
		o.MaxArraySize = size
	}
}

func WithMaxNestingDepth(depth int) LoadOption {
	return func(o *LoadOptions) {
		o.MaxNestingDepth = depth
	}
}

func WithTimeout(timeout time.Duration) LoadOption {
	return func(o *LoadOptions) {
		o.Timeout = timeout
	}
}
