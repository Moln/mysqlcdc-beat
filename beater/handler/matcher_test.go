package handler

import (
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/moln/cdcbeat/config"
	"strconv"
	"testing"
)

var logger = logp.NewLogger("testing")

func TestNotMatch(t *testing.T) {
	configs := [][]*config.MatchItemConfig{
		{
			{
				Database: "foo",
			},
		},
		{
			{
				Tables: []string{
					`foo_[a-z]+`,
				},
			},
		},
		{
			{
				Database: "test",
				Tables: []string{
					`foo_[a-z]+`,
				},
			},
		},
	}

	for i, cfg := range configs {
		t.Run("not match:"+strconv.Itoa(i), func(t *testing.T) {
			m := NewMatcher(cfg, logger)
			if m.match("test", "foo_123") {
				t.Error("match failed")
			}
		})
	}
}

func TestMatch(t *testing.T) {
	configs := [][]*config.MatchItemConfig{
		{},
		{
			{
				Database: "test",
			},
		},
		{
			{
				Database: "t.*",
			},
		},
		{
			{
				Tables: []string{
					`foo_123`,
				},
			},
		},
		{
			{
				Database: "test",
				Tables: []string{
					`foo_\d+`,
				},
			},
		},
	}

	for i, cfg := range configs {
		t.Run("match:"+strconv.Itoa(i), func(t *testing.T) {
			m := NewMatcher(cfg, logger)
			if !m.match("test", "foo_123") {
				t.Error("match failed")
			}
		})
	}
}
