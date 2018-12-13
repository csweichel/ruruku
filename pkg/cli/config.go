package cli

type Config struct {
	Host    string
	TLSCert string `yaml:"tlsCert"`
	Timeout uint32
	Token   string
}
