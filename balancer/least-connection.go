package balancer

import "net/url"

type Server struct {
	url   *url.URL
	count uint64
}

type ServerPool []*Server

func intitServers() (ServerPool, error) {
	servers := []string{
		"http://localhost:8000",
		"http://localhost:8001",
		"http://localhost:8002",
	}

	pool := make(ServerPool, 0, len(servers))

	for _, server := range servers {
		parsedServer, err := url.Parse(server)

		if err != nil {
			return nil, err
		}

		hydratedUrl := &Server{parsedServer, 0}

		pool = append(pool, hydratedUrl)
	}

	return pool, nil
}

func (pool ServerPool) NextServer() *Server {

	if len(pool) == 0 {
		return nil
	}

	freeServer := pool[0]

	for _, server := range pool[1:] {
		if server.count < freeServer.count {
			freeServer = server
		}
	}

	return freeServer
}
