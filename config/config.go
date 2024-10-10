package config

import "github.com/caarlos0/env"

// ServerConfiguration OS environment
type ServerConfiguration struct {
	Port                    int     `env:"PORT" envDefault:"8080"`
	SSL                     int     `env:"SSL" envDefault:"0"`
	CertPath                string  `env:"CERT_PATH" envDefault:"./cert/cert.pem"`
	KeyPath                 string  `env:"KEY_PATH" envDefault:"./cert/key.pem"`
}

// Get
func Get() (configuration *ServerConfiguration, err error) {
	configuration = new(ServerConfiguration)
	err = env.Parse(configuration)
	return
}
