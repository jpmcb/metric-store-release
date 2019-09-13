package app

import (
	"log"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

type CAPI struct {
	ExternalAddr string `env:"CAPI_ADDR_EXTERNAL, required, report"`
	CAPath       string `env:"CAPI_CA_PATH,     required, report"`
	CommonName   string `env:"CAPI_COMMON_NAME, required, report"`
}

type UAA struct {
	ClientID     string `env:"UAA_CLIENT_ID,     required"`
	ClientSecret string `env:"UAA_CLIENT_SECRET, required"`
	Addr         string `env:"UAA_ADDR,          required, report"`
	CAPath       string `env:"UAA_CA_PATH,       required, report"`
}

type Config struct {
	MetricStoreAddr  string `env:"METRIC_STORE_ADDR, required, report"`
	Addr             string `env:"ADDR, required, report"`
	InternalIP       string `env:"INTERNAL_IP, report"`
	HealthPort       int    `env:"HEALTH_PORT, report"`
	CertPath         string `env:"EXTERNAL_CERT, required, report"`
	KeyPath          string `env:"EXTERNAL_KEY, required, report"`
	SkipCertVerify   bool   `env:"SKIP_CERT_VERIFY, report"`
	ProxyCAPath      string `env:"PROXY_CA_PATH, required, report"`
	SecurityEventLog string `env:"SECURITY_EVENT_LOG, report"`

	CAPI CAPI
	UAA  UAA

	LogLevel string `env:"LOG_LEVEL,                      report"`
}

func LoadConfig() *Config {
	cfg := &Config{
		LogLevel:        "info",
		SkipCertVerify:  false,
		Addr:            ":8083",
		InternalIP:      "0.0.0.0",
		HealthPort:      6065,
		MetricStoreAddr: "localhost:8080",
	}

	if err := envstruct.Load(cfg); err != nil {
		log.Fatalf("failed to load config from environment: %s", err)
	}

	_ = envstruct.WriteReport(cfg)

	return cfg
}