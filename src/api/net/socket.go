package net

import (
	"net"
)

type Socket interface {
	IsRunning() bool
	Start(host string, port uint16) error
	OnConnection() (net.Conn, error)
	Close() error
}
