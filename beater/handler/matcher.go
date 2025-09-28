package handler

import (
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/moln/cdcbeat/config"
	"regexp"
)

type matchItem struct {
	database *regexp.Regexp
	tables   []*regexp.Regexp
}

type TableMatcher struct {
	matches []*matchItem
}

func (m *TableMatcher) match(db string, table string) bool {
	for _, item := range m.matches {
		if item.database != nil && !item.database.MatchString(db) {
			return false
		}

		if len(item.tables) > 0 {
			for _, r := range item.tables {
				if r.MatchString(table) {
					return true
				}
			}
			return false
		}
	}

	return true
}

func NewMatcher(cfg []*config.MatchItemConfig, logger *logp.Logger) *TableMatcher {

	var matches []*matchItem
	for _, itemCfg := range cfg {

		item := &matchItem{}
		if itemCfg.Database != "" {
			dbRe, err := regexp.Compile(itemCfg.Database)
			if err != nil {
				logger.Error(err)
				panic(err)
			}
			item.database = dbRe
		}

		for _, table := range itemCfg.Tables {
			tableRe, err := regexp.Compile(table)
			if err != nil {
				logger.Error(err)
				panic(err)
			}
			item.tables = append(item.tables, tableRe)
		}
		matches = append(matches, item)
	}

	return &TableMatcher{
		matches: matches,
	}
}
