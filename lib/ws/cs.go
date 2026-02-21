package ws

import (
	"orglang/go-engine/lib/kv"
)

func newExchangeCS(loader kv.Loader) (exchangeCS, error) {
	dto := new(exchangeCS)
	loadingErr := loader.Load("exchange", dto)
	if loadingErr != nil {
		return exchangeCS{}, loadingErr
	}
	validateErr := dto.Validate()
	if validateErr != nil {
		return exchangeCS{}, validateErr
	}
	return *dto, nil
}

type exchangeCS struct {
	Protocol protocolCS `mapstructure:"protocol"`
	Server   serverCS   `mapstructure:"server"`
}

type protocolCS struct {
	Modes []protoModeCS `mapstructure:"modes"`
	HTTP  httpCS        `mapstructure:"http"`
}

type serverCS struct {
	Mode serverModeCS `mapstructure:"mode"`
	Echo echoCS       `mapstructure:"echo"`
}

type httpCS struct {
	Port int `mapstructure:"port"`
}

type echoCS struct{}

type protoModeCS string

const (
	httpProto protoModeCS = "http"
)

type serverModeCS string

const (
	echoServer serverModeCS = "echo"
)
