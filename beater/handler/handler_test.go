package handler

import (
	ec "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/moln/cdcbeat/config"
	"os"
	"testing"
)
import (
	"github.com/elastic/beats/v7/libbeat/publisher/pipeline"
)

func TestHandler(t *testing.T) {
	client, _ := pipeline.NewNilPipeline().Connect()
	os.Remove("position.yml")

	cfg := ec.NewConfig()
	cfg.Unpack("position.yml")
	h := NewEventHandler(
		client,
		&config.Registry{Path: "position.yml", Permission: os.ModePerm},
		cfg,
		logp.NewLogger("testing"),
	)

	pos := h.ReadPos()
	if pos != nil {
		t.Error("pos not empty ")
	}
	h.OnPosSynced(
		nil,
		mysql.Position{
			Name: "test",
			Pos:  111111,
		},
		nil,
		true,
	)
	h.OnPosSynced(
		nil,
		mysql.Position{
			Name: "test",
			Pos:  222,
		},
		nil,
		true,
	)
	pos = h.ReadPos()
	if pos.Name != "test" && pos.Pos != 222 {
		t.Error("pos read error ", pos)
	}

	os.Remove("position.yml")
}
