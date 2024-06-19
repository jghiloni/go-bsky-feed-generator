package config_test

import (
	"context"

	"github.com/bluesky-social/indigo/atproto/syntax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jghiloni/go-bsky-feed-generator/config"
)

var _ = Describe("Config", func() {
	It("Persists the config in the context", func() {
		cfg := config.FeedGeneratorConfig{
			ServiceDID: syntax.ATURI("at://jaygles.bsky.social"),
		}

		ctx := config.WithConfig(context.Background(), cfg)
		Expect(ctx).NotTo(BeZero())
		Expect(ctx).NotTo(BeEquivalentTo(context.Background()))

		retrievedCfg, err := config.GetConfig(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedCfg).To(Equal(cfg))
	})
})
