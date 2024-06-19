package config

import (
	"context"
	"errors"

	"github.com/bluesky-social/indigo/atproto/syntax"
)

// FeedGeneratorConfig represents configurable information about a generator
// server
type FeedGeneratorConfig struct {
	ServiceDID   syntax.ATURI
	PublisherDID syntax.ATURI
}

type configContextKey struct{}

var config configContextKey

// ErrConfigNotFound is returned by GetConfig if the config has not been set
// via WithConfig
var ErrConfigNotFound = errors.New("config not set")

// WithConfig inserts configuration for a feed server onto a context.Context
// object. The parent context must NOT be nil. If an existing context does not
// exist, use context.TODO() or context.Background()
func WithConfig(parent context.Context, cfg FeedGeneratorConfig) context.Context {
	if parent == nil {
		panic("parent must not be nil")
	}

	return context.WithValue(parent, config, cfg)
}

// GetConfig will get the FeedGeneratorConfig set by WithConfig. If the passed
// context is nil, GetConfig panics. If the config is not found, ErrConfigNotFound
// is returned
func GetConfig(parent context.Context) (FeedGeneratorConfig, error) {
	if parent == nil {
		panic("parent must not be nil")
	}

	fgc, ok := parent.Value(config).(FeedGeneratorConfig)
	if !ok {
		return FeedGeneratorConfig{}, errors.New("config not set")
	}

	return fgc, nil
}
