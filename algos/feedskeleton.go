package algos

import (
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/syntax"
)

// FeedInput is a collection of data that will be available to the feed
// algorithm
type FeedInput struct {
	VerifiedUser syntax.ATURI
	FeedName     string
	Limit        int
	Cursor       string
}

// BlueskyFeed represents an implementation of a feed Algorithm. A feed
// implementation has a short name (<= 16 characters), and an GenerateFeed
// method that includes the logic for the algorithm
type BlueskyFeed interface {
	ShortName() string // must be <= 16 characters, per protocol spec
	GenerateFeed(FeedInput) (bsky.FeedGetFeedSkeleton_Output, error)
}
