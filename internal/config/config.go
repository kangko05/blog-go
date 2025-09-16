package config

type Config struct {
	JwtSecret string
	DbPath    string
}

func Load() *Config {
	return &Config{
		JwtSecret: "TODO: TEMPORARY SECRET KEY",
		DbPath:    "dev.db",
	}
}
