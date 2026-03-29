package config

type Config struct {
	App   AppConfig
	DB    DBConfig
	Cache CacheConfig
}

type AppConfig struct {
	Port         string
	AllowOrigins string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type CacheConfig struct {
	Host     string
	Port     string
	Password string
}
