package query

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"

	"github.com/golangminecraft/minecraft-server/src/api"
	log "github.com/golangminecraft/minecraft-server/src/api/logger"
	"github.com/golangminecraft/minecraft-server/src/util"
)

type Server struct {
	isRunning bool
	server    api.Server
	listener  *net.UDPConn
}

type KVSection struct {
	Key   string
	Value string
}

func NewServer() api.QueryServer {
	return &Server{
		isRunning: false,
		server:    nil,
		listener:  nil,
	}
}

func (s *Server) Initialize(server api.Server) error {
	s.server = server

	return nil
}

func (s *Server) Start(host string, port uint16) error {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return err
	}

	listener, err := net.ListenUDP("udp4", addr)

	if err != nil {
		return err
	}

	s.listener = listener
	s.isRunning = true

	go s.AcceptConnections()

	return nil
}

func (s *Server) AcceptConnections() {
	for s.isRunning {
		if err := s.HandleConnection(); err != nil {
			log.Error("query", err)
		}
	}
}

func (s *Server) HandleConnection() error {
	var addr *net.UDPAddr
	var err error
	var sessionID int32

	// Handshake Request
	{
		data := make([]byte, 7)
		oob := make([]byte, 1024)

		_, _, _, addr, err = s.listener.ReadMsgUDP(data, oob)

		if err != nil {
			return err
		}

		if data[2] != 0x09 {
			return fmt.Errorf("unexpected packet type: 0x%02X", data[2])
		}

		sessionID = int32(binary.BigEndian.Uint32(data[3:7])) & 0x0F0F0F0F
	}

	challengeToken := rand.Int31()

	// Handshake Response
	{
		challengeTokenString := strconv.FormatInt(int64(challengeToken), 10)

		buf := &bytes.Buffer{}

		if err := buf.WriteByte(0x09); err != nil {
			return err
		}

		if err := binary.Write(buf, binary.BigEndian, sessionID); err != nil {
			return err
		}

		if _, err := buf.Write([]byte(challengeTokenString)); err != nil {
			return err
		}

		if err := buf.WriteByte(0x00); err != nil {
			return err
		}

		if _, err := s.listener.WriteToUDP(buf.Bytes(), addr); err != nil {
			return err
		}
	}

	// Stat
	{
		data := make([]byte, 15)
		oob := make([]byte, 1024)

		n, _, _, addr, err := s.listener.ReadMsgUDP(data, oob)

		if err != nil {
			return err
		}

		if data[2] != 0x00 {
			return fmt.Errorf("unexpected packet type: 0x%02X", data[2])
		}

		if n > 11 {
			// Full Query

			if err := s.SendFullQuery(addr, sessionID); err != nil {
				return err
			}

			log.Infof("query", "Received full query request from %s\n", addr.String())
		} else {
			// Basic Query

			if err := s.SendBasicQuery(addr, sessionID); err != nil {
				return err
			}

			log.Infof("query", "Received basic query request from %s\n", addr.String())
		}
	}

	return nil
}

func (s Server) SendBasicQuery(addr *net.UDPAddr, sessionID int32) error {
	buf := &bytes.Buffer{}

	if err := buf.WriteByte(0x00); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, sessionID); err != nil {
		return err
	}

	if _, err := buf.Write(util.StringNT(s.server.MOTD().Format())); err != nil {
		return err
	}

	if _, err := buf.Write(util.StringNT("SMP")); err != nil {
		return err
	}

	if _, err := buf.Write(util.StringNT("world")); err != nil {
		return err
	}

	if _, err := buf.Write(util.StringNT(strconv.FormatInt(int64(s.server.OnlinePlayers()), 10))); err != nil {
		return err
	}

	if _, err := buf.Write(util.StringNT(strconv.FormatInt(int64(s.server.MaxPlayers()), 10))); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.LittleEndian, s.server.Port()); err != nil {
		return err
	}

	if _, err := buf.Write(util.StringNT(s.server.Host())); err != nil {
		return err
	}

	_, err := s.listener.WriteToUDP(buf.Bytes(), addr)

	return err
}

func (s Server) SendFullQuery(addr *net.UDPAddr, sessionID int32) error {
	buf := &bytes.Buffer{}

	if err := buf.WriteByte(0x00); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, sessionID); err != nil {
		return err
	}

	if _, err := buf.Write([]byte{0x73, 0x70, 0x6C, 0x69, 0x74, 0x6E, 0x75, 0x6D, 0x00, 0x80, 0x00}); err != nil {
		return err
	}

	kvSection := []KVSection{
		{Key: "hostname", Value: s.server.MOTD().Format()},
		{Key: "game type", Value: "SMP"},
		{Key: "game_id", Value: "MINECRAFT"},
		{Key: "version", Value: "1.18.1"},
		{Key: "plugins", Value: "GoLangMinecraft Server"},
		{Key: "map", Value: "world"},
		{Key: "numplayers", Value: strconv.FormatInt(int64(s.server.OnlinePlayers()), 10)},
		{Key: "maxplayers", Value: strconv.FormatInt(int64(s.server.MaxPlayers()), 10)},
		{Key: "hostport", Value: strconv.FormatInt(int64(s.server.Port()), 10)},
		{Key: "hostip", Value: s.server.Host()},
	}

	for _, value := range kvSection {
		if _, err := buf.Write([]byte(value.Key)); err != nil {
			return err
		}

		if err := buf.WriteByte(0x00); err != nil {
			return err
		}

		if _, err := buf.Write([]byte(value.Value)); err != nil {
			return err
		}

		if err := buf.WriteByte(0x00); err != nil {
			return err
		}
	}

	if err := buf.WriteByte(0x00); err != nil {
		return err
	}

	if _, err := buf.Write([]byte{0x01, 0x70, 0x6C, 0x61, 0x79, 0x65, 0x72, 0x5F, 0x00, 0x00}); err != nil {
		return err
	}

	for _, player := range s.server.Players() {
		if _, err := buf.Write([]byte(player.Username())); err != nil {
			return err
		}

		if err := buf.WriteByte(0x00); err != nil {
			return err
		}
	}

	if err := buf.WriteByte(0x00); err != nil {
		return err
	}

	_, err := s.listener.WriteToUDP(buf.Bytes(), addr)

	return err
}

func (s *Server) Close() error {
	s.isRunning = false

	return s.listener.Close()
}
