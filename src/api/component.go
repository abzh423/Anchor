package api

type Component interface {
	Initialize(Server) error
	Start() error
	AddClient(Client)
	RemoveClient(Client)
}
