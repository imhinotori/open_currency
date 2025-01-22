package configuration

type Configuration struct {
	General
	HTTP   HTTP
	Google Google
}

type General struct {
	DebugEnabled bool
}

type HTTP struct {
	Address string
	Port    int
	SSL     bool
	SSLCert string
	SSLKey  string
}

type Google struct {
	CredentialsFile string
	SheetId         string
}

func Load() (*Configuration, error) {
	cfg := &Configuration{
		General: General{
			DebugEnabled: GetEnvBool("DEBUG_ENABLED"),
		},
		HTTP: HTTP{
			Address: GetEnv("HTTP_ADDRESS", "0.0.0.0"),
			Port:    GetEnvInt("HTTP_PORT", 8080),
			SSL:     GetEnvBool("HTTP_SSL"),
			SSLCert: GetEnv("HTTP_SSL_CERT", "/etc/ssl/certs/open_currency.crt"),
			SSLKey:  GetEnv("HTTP_SSL_KEY", "/etc/ssl/private/open_currency.key"),
		},
		Google: Google{
			CredentialsFile: GetEnv("GOOGLE_CREDENTIALS_FILE", "/etc/google-service-account.json"),
			SheetId:         GetEnv("GOOGLE_SHEET_ID", ""),
		},
	}

	return cfg, nil
}
