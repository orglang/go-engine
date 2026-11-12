package kv

import (
	"log/slog"
	"reflect"

	"github.com/spf13/viper"
)

func newViperDriver(logger *slog.Logger) (*viperDriver, error) {
	driver := viper.New()
	driver.AddConfigPath(".")
	driver.SetConfigType("yaml")
	driver.SetConfigName("reference")
	if err := driver.ReadInConfig(); err != nil {
		return nil, err
	}
	driver.SetConfigName("application")
	if err := driver.MergeInConfig(); err != nil {
		logger.Warn("falling back to default config")
	}
	name := slog.String("name", reflect.TypeFor[viperDriver]().Name())
	return &viperDriver{driver, logger.With(name)}, nil
}

type viperDriver struct {
	viper *viper.Viper
	log   *slog.Logger
}

func (d *viperDriver) Load(key string, dst any) error {
	err := d.viper.UnmarshalKey(key, dst)
	if err != nil {
		d.log.Error("load failed", slog.String("key", key), slog.Any("reason", err))
		return err
	}
	d.log.Info("load succeed", slog.String("key", key), slog.Any("val", dst))
	return nil
}
