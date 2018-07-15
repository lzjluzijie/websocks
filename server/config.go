package server

type Config struct {
	ListenAddr   string
	Pattern      string
	TLS          bool
	CertPath     string
	KeyPath      string
	ReverseProxy string
}
