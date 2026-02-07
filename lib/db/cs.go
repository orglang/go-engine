package db

import (
	"orglang/go-engine/lib/kv"
)

func newStorageCS(l kv.Loader) (storageCS, error) {
	dto := new(storageCS)
	loadingErr := l.Load("storage", dto)
	if loadingErr != nil {
		return storageCS{}, loadingErr
	}
	validationErr := dto.Validate()
	if validationErr != nil {
		return storageCS{}, validationErr
	}
	return *dto, nil
}

type storageCS struct {
	Protocol protocolCS `mapstructure:"protocol"`
	Driver   driverCS   `mapstructure:"driver"`
}

type protocolCS struct {
	Mode     protoModeCS `mapstructure:"mode"`
	Postgres postgresCS  `mapstructure:"postgres"`
}

type driverCS struct {
	Mode driverModeCS `mapstructure:"mode"`
	Pgx  pgxCS        `mapstructure:"pgx"`
}

type postgresCS struct {
	URL string `mapstructure:"url"`
}

type pgxCS struct {
	MaxConns uint8 `mapstructure:"max_conns"`
}

type protoModeCS string

const (
	postgresProto protoModeCS = "postgres"
)

type driverModeCS string

const (
	pgxDriver driverModeCS = "pgx"
)
