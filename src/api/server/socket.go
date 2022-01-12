package server

type Socket interface {
	IsRunning() bool
	Start(host string, port uint16) error
	Stop() error
	OnConnection() (Client, error)
}
