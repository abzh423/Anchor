package api

type Socket interface {
	Start(string, uint16) error
	Close() error
	AcceptConnection() (Client, error)
}
