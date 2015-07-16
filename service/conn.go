package service

import (
	"encoding/gob"
	"log"
	"net"
	"sync"
)

type Conn struct {
	conn    net.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
	sync.RWMutex
}

func (s *Conn) Close() {
	s.conn.Close()
}

func (s *Conn) Read() (interface{}, error) {
	var m ServerMessage

	err := s.decoder.Decode(&m)
	if err != nil {
		return nil, err
	}

	return m.Data, nil
}

func (s *Conn) Write(data interface{}) error {
	s.Lock()
	defer s.Unlock()

	msg := ServerMessage{data}
	err := s.encoder.Encode(msg)
	return err
}

func (s *Conn) Authenticate(typeId uint32, uuid string, key string) error {
	return s.Write(AuthenticateMsg{TypeId: typeId, UUID: uuid, Key: key})
}

func (s *Conn) GetType(id uint32) error {
	return s.Write(GetTypeMsg{TypeId: id})
}

func (s *Conn) SendType(id uint32, src string, masks []uint32, intervals []uint8) error {
	return s.Write(TypeMsg{
		TypeId:    id,
		Src:       src,
		Masks:     masks,
		Intervals: intervals,
	})
}

func (s *Conn) RequestJoin(UUID string, typeId uint32) error {
	return s.Write(JoinRequestMsg{
		TypeId: typeId,
		UUID:   UUID,
	})
}

func NewConn(conn net.Conn) (*Conn, error) {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	return &Conn{conn: conn, encoder: enc, decoder: dec}, nil
}
