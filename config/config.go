package config

// Config Application configuration structure
type Config struct {
	DB *DBConfig
}

// DBConfig structure
type DBConfig struct {
	Dialect  string
	Hostname string
	Port     uint16
	Username string
	Password string
	Name     string
	Charset  string
	Prefix   string
}

// GetConfig get Application configuration
func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "postgres",
			Hostname: "localhost",
			Port:     5432,
			Username: "shorty",
			Password: "root",
			Name:     "ruian",
			Charset:  "utf8",
			Prefix:   "view_address_",
		},
	}
}
