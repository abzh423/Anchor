package api

type QueryServer interface {
	Initialize(Server) error
	Start(host string, port uint16) error
	Close() error
}
