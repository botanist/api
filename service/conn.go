package service

import (
	"encoding/gob"
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

func (s *Conn) Hello(ver uint32) error {
	return s.Write(HelloMsg{ProtocolVersion:ver})
}

func (s *Conn) Authenticate(typeId uint32, uuid string, key string) error {
	return s.Write(AuthenticateMsg{TypeId: typeId, UUID: uuid, Key: key})
}

func (s *Conn) AuthenticationFailed(reason string) error {
	return s.Write(AuthenticationFailedMsg{Reason:reason})
}

func (s *Conn) AuthenticationSucceeded(deviceId uint32, rfAddr uint16, newKey string) error {
	return s.Write(AuthenticationSucceededMsg{DeviceId: deviceId, RFAddr: rfAddr, NewKey: newKey})
}

func (s *Conn) GetRFAddr(id uint32) error {
	return s.Write(GetRFAddrMsg{DeviceId: id})
}

func (s *Conn) SendRFAddr(id uint32, addr uint16) error {
	return s.Write(RFAddrMsg{DeviceId: id, RFAddr: addr})
}

/* Types */
func (s *Conn) GetType(id uint32) error {
	return s.Write(GetTypeMsg{TypeId: id})
}

func (s *Conn) SendNoType(id uint32) error {
	return s.Write(NoTypeMsg{TypeId: id})
}

func (s *Conn) SendType(id uint32, src string, masks []uint32, intervals []uint8) error {
	return s.Write(TypeMsg{
		TypeId:    id,
		Src:       src,
		Masks:     masks,
		Intervals: intervals,
	})
}

/* JOINING */
func (s *Conn) RequestJoin(UUID string, typeId uint32) error {
	return s.Write(JoinRequestMsg{
		TypeId: typeId,
		UUID:   UUID,
	})
}

func (s *Conn) JoinRequestPending(UUID string) error {
	return s.Write(JoinRequestPendingMsg{UUID: UUID})
}

func (s *Conn) JoinRequestDeclined(UUID string) error {
	return s.Write(JoinRequestDeclinedMsg{UUID: UUID})
}

func (s *Conn) JoinRequestApproved(UUID string, id uint32, addr uint16) error {
	return s.Write(JoinRequestApprovedMsg{UUID: UUID, DeviceId: id, RFAddr: addr})
}

/* WIRED SENSORS */
func (s *Conn) ConnectDeviceDeclined(UUID string, parentId uint32) error {
	return s.Write(ConnectDeviceDeclinedMsg{UUID: UUID, ParentDeviceId: parentId})
}

func (s *Conn) ConnectDeviceApproved(UUID string, parentId, id uint32) error {
	return s.Write(ConnectDeviceApprovedMsg{UUID: UUID, ParentDeviceId: parentId, DeviceId: id})
}

/* Utilities */

func (s *Conn) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func NewConn(conn net.Conn) (*Conn, error) {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	return &Conn{conn: conn, encoder: enc, decoder: dec}, nil
}
