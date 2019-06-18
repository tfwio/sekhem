package config

import (
	"fmt"

	"tfw.io/Go/fsindex/util"
)

// Server info for JSON i/o.
type Server struct {
	Port string `json:"port"`
	Host string `json:"host"`
	// default=false unless `os.Args[1] == "tls"` or specified
	// in `[data/]config.json`.
	TLS bool   `json:"tls"`
	Crt string `json:"crt,omitempty"`
	Key string `json:"key,omitempty"`
}

func (s *Server) info() {
	println("> Server")
	println(fmt.Sprintf("--> Host = %s", s.Host))
	println(fmt.Sprintf("--> Port = %s", s.Port))
	println(fmt.Sprintf("--> TLS  = %v", s.TLS))
	println(fmt.Sprintf("--> Key  = %s", s.Key))
	println(fmt.Sprintf("--> Crt  = %s", s.Crt))
}
func (s *Server) hasKey() bool {
	return util.FileExists(s.Key)
}
func (s *Server) hasCert() bool {
	return util.FileExists(s.Crt)
}

func (s *Server) initServerConfig() {
	s.Host = constServerDefaultHost
	s.Port = constServerDefaultPort
	s.TLS = UseTLS
	s.Crt = constServerTLSCertDefault
	s.Key = constServerTLSKeyDefault
}
