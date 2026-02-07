package property

type ExchangeCS struct {
	Protocol ProtocolCS `mapstructure:"protocol"`
	Server   ServerCS   `mapstructure:"server"`
	Client   ClientCS   `mapstructure:"client"`
}

type ProtocolCS struct {
	Modes []ProtoModeCS `mapstructure:"modes"`
	HTTP  httpCS        `mapstructure:"http"`
}

type ServerCS struct {
	Mode ServerModeCS `mapstructure:"mode"`
	Echo EchoCS       `mapstructure:"echo"`
}

type ClientCS struct {
	Mode  ClientModeCS `mapstructure:"mode"`
	Resty RestyCS      `mapstructure:"resty"`
}

type httpCS struct {
	Port uint16 `mapstructure:"port"`
}

type EchoCS struct{}
type RestyCS struct{}

type ProtoModeCS string

const (
	HTTPProto ProtoModeCS = "http"
)

type ServerModeCS string

const (
	EchoServer ServerModeCS = "echo"
)

type ClientModeCS string

const (
	RestyClient ClientModeCS = "resty"
)
