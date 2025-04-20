package models

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Redis struct {
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	SQLite struct {
		Path            string `mapstructure:"path"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 秒
	} `mapstructure:"sqlite"`

	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`

	RAG struct {
		ServiceAddr string `mapstructure:"service_addr"`
	} `mapstructure:"rag"`
}
