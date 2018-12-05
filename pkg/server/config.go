package server

type Config struct {
	GRPC struct {
		Enabled bool
		Port    uint32
	}
	UI struct {
		Enabled bool
		Port    uint32
	}
	DB struct {
		Filename string
	}
}
