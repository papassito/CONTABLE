package config

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Version     string
}

type DatabaseConfig struct {
	Driver string
	Host   string
	Port   int
	User   string
	Pass   string
	Name   string
	SSL    bool
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

func LoadConfig() (*Config, error) {
	// TODO: Cargar variables de entorno / archivo config
	return &Config{
		App: AppConfig{
			Name:        "Contable Fix by KLIK",
			Environment: "development",
			Version:     "1.0.0",
		},
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  15,
			WriteTimeout: 15,
		},
	}, nil
}
