package server

type Config struct {
	GRPC struct {
		Enabled bool
		Port    uint32
	}
	UI struct {
		Enabled bool
		Port    uint32
		TLS     bool
		Key     string
		Cert    string
	}
	DB struct {
		Filename string
	}
	TLS struct {
		Enabled bool

		// Key is the path to the private key file
		Key string
		// Cert is the path to the certificate file
		Cert string
	}
}
