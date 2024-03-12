package clickhouse

type Config struct {
	Server   string `envconfig:"server"`
	Port     string `envconfig:"port" default:"9000"`
	Database string `envconfig:"database" default:"default"`
	User     string `envconfig:"user"`
	Password string `envconfig:"password"`
}
