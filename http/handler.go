package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/jghiloni/go-bsky-feed-generator/algos"
	"github.com/jghiloni/go-bsky-feed-generator/config"
)

var feedNSID = syntax.NSID("app.bsky.feed.generator")

// ServeFeeds will add the handler created by FeedHandler to a list of
// *http.ServeMux objects. If none are set, http.DefaultServeMux is used
// (which is equivalent to calling net/http.Handle())
func ServeFeeds(ctx context.Context, feeds []algos.BlueskyFeed, muxes ...*http.ServeMux) {
	if len(muxes) == 0 {
		muxes = []*http.ServeMux{http.DefaultServeMux}
	}

	for _, mux := range muxes {
		mux.Handle("/xrpc/app.bsky.feed.getFeedSkeleton", FeedHandler(ctx, feeds))
	}
}

// FeedHandler takes a context that has been created by github.com/jghiloni/go-bsky-feed-generator/config.WithConfig
// and a slice of feed implementations, and creates a standard net/http Handler
// for them. It does basic validation,
// including:
//
//  1. Ensuring that the feed AT URI Authority matches the configuration's
//     service DID
//  2. Ensuring that the feed URI collection is app.bsky.feed.generator
//
// The limit and cursor parameters are NOT validated (aside from ensuring that
// the limit parameter is an integer), and that is left to the individual
// algorithm implementation
func FeedHandler(ctx context.Context, feeds []algos.BlueskyFeed) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feedURI, err := syntax.ParseATURI(r.FormValue("feed"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		cfg, err := config.GetConfig(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		if feedURI.Authority() != cfg.ServiceDID.Authority() ||
			feedURI.Collection() != feedNSID {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		input := algos.FeedInput{
			FeedName: string(feedURI.RecordKey()),
			Cursor:   r.FormValue("cursor"),
		}

		input.Limit, _ = strconv.Atoi(r.FormValue("limit"))
		feed := findFeed(feeds, feedURI.RecordKey().String())
		if feed == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		output, err := feed.GenerateFeed(input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(output)
	})
}

func findFeed(feeds []algos.BlueskyFeed, feedName string) algos.BlueskyFeed {
	for _, feed := range feeds {
		if feed.ShortName() == feedName {
			return feed
		}
	}

	return nil
}
