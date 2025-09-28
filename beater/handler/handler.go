package handler

import (
	"fmt"
	"github.com/elastic/beats/v7/libbeat/beat"
	ec "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/elastic/elastic-agent-libs/paths"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/moln/cdcbeat/config"
	"io"
	"os"
	"regexp"
	"sigs.k8s.io/yaml"
	"sync"
	"time"
)
import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/replication"
)

type BeatEventHandler struct {
	canal.DummyEventHandler
	client      beat.Client
	logger      *logp.Logger
	fileFactory func() *os.File
}

func (h *BeatEventHandler) getFile() *os.File {
	return h.fileFactory()
}

func (h *BeatEventHandler) OnRow(e *canal.RowsEvent) error {
	//var rows []mapstr.M
	//for _, row := range e.Rows {
	//	kv := mapstr.M{}
	//	cols := []string{}
	//	for i, col := range e.Table.Columns {
	//		kv[col.Name] = row[i]
	//		cols = append(cols, col.Name)
	//	}
	//	rows = append(rows, kv)
	//}
	cols := []string{}
	for _, col := range e.Table.Columns {
		cols = append(cols, col.Name)
	}

	fields := mapstr.M{
		"action": e.Action,
		"table": mapstr.M{
			"name":    e.Table.Name,
			"schema":  e.Table.Schema,
			"pk":      e.Table.PKColumns,
			"columns": cols,
		},
		"rows": e.Rows,
	}
	ts := time.Now()
	if e.Header != nil {
		fields["server_id"] = e.Header.ServerID
		fields["pos"] = e.Header.LogPos
		ts = time.Unix(int64(e.Header.Timestamp), 0)
	}

	event := beat.Event{
		Timestamp: ts,
		Fields:    fields,
	}
	h.client.Publish(event)
	return nil
}

func (h *BeatEventHandler) OnTableChanged(header *replication.EventHeader, schema, table string) error {
	h.logger.Debug("OnTableChanged: ", schema, table)
	return nil
}

func (h *BeatEventHandler) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	h.logger.Debug("OnDDL: ", nextPos)
	return nil
}

func (h *BeatEventHandler) String() string {
	return "BeatEventHandler"
}

func (h *BeatEventHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	h.logger.Debug("OnPosSynced: ", pos)
	text, _ := yaml.Marshal(pos)
	file := h.getFile()
	_, err := file.WriteAt(text, 0)
	file.Truncate(int64(len(text)))
	if err != nil {
		return err
	}
	return nil
}

func (h *BeatEventHandler) Close() error {
	return h.getFile().Close()
}

func (h *BeatEventHandler) ReadPos() *mysql.Position {
	file := h.getFile()
	file.Seek(0, 0)
	data, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error(err)
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	pos := &mysql.Position{}
	err = yaml.Unmarshal(data, pos)
	if err != nil {
		h.logger.Error(err)
		return nil
	}
	return pos
}

func NewEventHandler(bc beat.Client, registry *config.Registry, cfg *ec.C, logger *logp.Logger) *BeatEventHandler {
	fileFactory := sync.OnceValue(func() *os.File {
		addr, err := cfg.String("addr", 0)
		path := paths.Resolve(
			paths.Data,
			fmt.Sprintf(
				registry.Path+"%[2]s",
				regexp.MustCompile(`[.\[\]:]`).ReplaceAllString(addr, "_"),
				"",
			),
		)
		logger.Debug("position save path: ", path)
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, registry.Permission)
		if err != nil {
			logger.Error(err)
			panic(err)
		}

		return file
	})

	return &BeatEventHandler{
		client:      bc,
		logger:      logger,
		fileFactory: fileFactory,
	}
}
