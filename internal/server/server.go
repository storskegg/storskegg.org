package server

type Server interface {
	Serve() error
}

type server struct {
	config *Config
}

func New(config *Config) Server {
	return &server{
		config: config,
	}
}

func (s *server) Serve() error {
	return nil
}
