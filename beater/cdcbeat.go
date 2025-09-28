package beater

import (
	"fmt"
	"github.com/elastic/beats/v7/libbeat/beat"
	ec "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
	"github.com/elastic/go-ucfg"
	"github.com/go-mysql-org/go-mysql/canal"
	hd "github.com/moln/cdcbeat/beater/handler"
	"github.com/moln/cdcbeat/config"
	"github.com/siddontang/go-log/loggers"
	"net/url"
	"sync"
)

// cdcbeat configuration.
type cdcbeat struct {
	done        chan struct{}
	config      *config.Config
	client      beat.Client
	handler     *hd.BeatEventHandler
	connections []*canal.Canal
	logger      *logp.Logger
}

// New creates an instance of cdcbeat.
func New(b *beat.Beat, cfg *ec.C) (beat.Beater, error) {
	c := config.NewConfig()

	b.Info.Logger.Info("cdcbeat New.")
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &cdcbeat{
		done:        make(chan struct{}),
		config:      c,
		connections: make([]*canal.Canal, 0),
		logger:      b.Info.Logger.Named("cdcbeat"),
	}
	return bt, nil
}

// Run starts cdcbeat.
func (bt *cdcbeat) Run(b *beat.Beat) error {
	bt.logger.Info("cdcbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	logger := hd.NewLogpProxyLogger(bt.logger)

	wg := sync.WaitGroup{}
	for _, cfg := range bt.config.Dbs {
		cc, err := newCanalConfig(cfg, logger)
		if err != nil {
			return err
		}

		handler := hd.NewEventHandler(bt.client, bt.config.Registry, cfg, bt.logger)
		cn, err := canal.NewCanal(cc)
		cn.SetEventHandler(handler)
		bt.connections = append(bt.connections, cn)
		wg.Add(1)
		go func(cn *canal.Canal, handler *hd.BeatEventHandler) {
			defer wg.Done()
			pos := handler.ReadPos()
			defer handler.Close()

			var err error
			if pos == nil {
				bt.logger.Info("Read empty position from file.")
				err = cn.Run()
			} else {
				bt.logger.Info("Start position from file: ", pos)
				err = cn.RunFrom(*pos)
			}

			if err != nil {
				bt.logger.Error(err)
			}
		}(cn, handler)
	}
	wg.Wait()

	return nil
}

func newCanalConfig(cfg *ec.C, logger loggers.Advanced) (*canal.Config, error) {
	cc := canal.NewDefaultConfig()
	cc.Logger = logger

	opts := []ucfg.Option{
		ucfg.PathSep("."),
		ucfg.ResolveEnv,
		ucfg.VarExp,
		ucfg.StructTag("toml"),
	}
	cfg2 := ucfg.New()
	cfg2.Merge(cfg, opts...)
	err := cfg2.Unpack(cc, opts...)

	if err != nil {
		return nil, err
	}

	parse, err := url.Parse("db://" + cc.Addr)
	if err != nil {
		return nil, err
	}
	port := parse.Port()
	if port == "" {
		cc.Addr += ":3306"
	}

	if !cfg.HasField("dump") {
		cc.Dump.ExecutionPath = ""
	}
	return cc, nil
}

// Stop stops cdcbeat.
func (bt *cdcbeat) Stop() {
	bt.logger.Info("cdcbeat is stopping.")
	for _, cn := range bt.connections {
		cn.Close()
	}
	bt.client.Close()
	close(bt.done)
}
