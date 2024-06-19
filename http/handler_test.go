package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	nh "net/http"
	"net/http/httptest"

	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/syntax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jghiloni/go-bsky-feed-generator/algos"
	"github.com/jghiloni/go-bsky-feed-generator/config"
	"github.com/jghiloni/go-bsky-feed-generator/http"
)

type testFeed struct{}

func (t *testFeed) ShortName() string {
	return "testFeed"
}

func (t *testFeed) GenerateFeed(input algos.FeedInput) (bsky.FeedGetFeedSkeleton_Output, error) {
	return bsky.FeedGetFeedSkeleton_Output{
		Cursor: nil,
		Feed: []*bsky.FeedDefs_SkeletonFeedPost{
			{
				Post: "at://did.plc.directory/1234",
			},
		},
	}, nil
}

var _ = Describe("Handler", func() {
	var s *httptest.Server

	BeforeEach(func() {
		s = httptest.NewServer(http.FeedHandler(config.WithConfig(context.Background(), config.FeedGeneratorConfig{
			ServiceDID: syntax.ATURI("at://did:plc:e2fun4xcfwtcrqfdwhfnghxk"),
		}), []algos.BlueskyFeed{&testFeed{}}))
	})

	AfterEach(func() {
		s.Close()
	})

	It("Returns correctly", func() {
		resp, err := nh.Get(fmt.Sprintf("%s/xrpc/app.bsky.feed.getFeedSkeleton?feed=at://did:plc:e2fun4xcfwtcrqfdwhfnghxk/app.bsky.feed.generator/testFeed", s.URL))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(nh.StatusOK))

		var o bsky.FeedGetFeedSkeleton_Output
		err = json.NewDecoder(resp.Body).Decode(&o)
		Expect(err).NotTo(HaveOccurred())

		Expect(o.Cursor).To(BeNil())
		Expect(o.Feed).To(HaveLen(1))
		Expect(o.Feed[0].Post).To(Equal("at://did.plc.directory/1234"))
	})

	DescribeTable("Error conditions", func(url string, statusCode int) {
		resp, err := nh.Get(fmt.Sprintf("%s/xrpc/app.bsky.feed.getFeedSkeleton?feed=%s", s.URL, url))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(statusCode))
	},
		Entry("invalid feed uri", "a2://foo.bar", 500),
		Entry("invalid authority", "at://x/1234/5678", 500),
		Entry("mismatched service did", "at://did:plc:invalid", 404),
		Entry("invalid collection", "at://did:plc:e2fun4xcfwtcrqfdwhfnghxk/app.sky.invalid.generator/testFeed", 404),
		Entry("invalid record key", "at://did:plc:e2fun4xcfwtcrqfdwhfnghxk/app.bsky.feed.generator/invalid", 404),
	)
})
