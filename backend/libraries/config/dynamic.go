package config

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/jace-ys/retune"
	"github.com/jace-ys/retune/encoder/yaml"
	"github.com/jace-ys/retune/loader"
	"github.com/jace-ys/retune/reader"
	"github.com/jace-ys/retune/reader/json"
	"github.com/jace-ys/retune/source/file"
	"go.uber.org/zap"
)

var (
	cfg retune.Config
	log *zap.SugaredLogger
)

func MustLoadDynamic(logger *zap.SugaredLogger, filepath string) {
	log = logger

	opts := []retune.Option{
		retune.WithReader(json.NewReader(reader.WithEncoder(yaml.NewEncoder()))),
	}

	var err error
	cfg, err = retune.NewConfig(opts...)
	if err != nil {
		panic(err)
	}

	if err := cfg.Load(file.NewSource(file.WithPath(filepath))); err != nil {
		panic(err)
	}
}

func Stop() error {
	return cfg.Close()
}

type Value interface {
	reader.Value
}

func Get(path ...string) Value {
	return cfg.Get(path...)
}

func Watch(ctx context.Context, callback func(val Value), path ...string) error {
	wlog := watchLogger(log, path...)

	w, err := cfg.Watch(path...)
	if err != nil {
		return fmt.Errorf("error creating watcher for config: %w", err)
	}
	wlog.Debugw("created config watcher")

	go func() {
		<-ctx.Done()
		w.Stop()
		wlog.Debugw("stopped config watcher")
	}()

	go func() {
		for {
			v, err := w.Next()
			if err != nil {
				if !errors.Is(err, loader.ErrWatcherStopped) {
					wlog.Errorw("error watching config", "error", err)
				}
				return
			}

			if string(v.Bytes()) != "null" {
				wlog.Debugw("watched config updated", "value", v)
				callback(v)
			}
		}
	}()

	return nil
}

func watchLogger(logger *zap.SugaredLogger, path ...string) *zap.SugaredLogger {
	if _, file, line, ok := runtime.Caller(2); ok {
		source := fmt.Sprintf("%s/%s:%d", filepath.Base(filepath.Dir(file)), filepath.Base(file), line)
		return logger.With("watcher.source", source, "path", path)
	}
	return logger.With("path", path)
}
