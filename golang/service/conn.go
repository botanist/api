package service

import (
	"encoding/gob"
	"errors"
	"net"
	"sync"
)

type Conn struct {
	conn    net.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
	sync.RWMutex
}

func (c *Conn) Close() {
	c.conn.Close()
}

func (c *Conn) Read() (interface{}, error) {
	var m ServerMessage

	err := c.decoder.Decode(&m)
	if err != nil {
		return nil, err
	}

	return m.Data, nil
}

func (c *Conn) Write(data interface{}) error {
	c.Lock()
	defer c.Unlock()

	msg := ServerMessage{data}
	err := c.encoder.Encode(msg)
	return err
}

func (c *Conn) Hello(ver uint32) error {
	return c.Write(HelloMsg{ProtocolVersion: ver})
}

func (c *Conn) Authenticate(typeId uint32, uuid string, key string) error {
	return c.Write(AuthenticateMsg{TypeId: typeId, UUID: uuid, Key: key})
}

func (c *Conn) AuthenticationFailed(reason string) error {
	return c.Write(AuthenticationFailedMsg{Reason: reason})
}

func (c *Conn) AuthenticationSucceeded(deviceId uint32, rfAddr uint16, newKey string) error {
	return c.Write(AuthenticationSucceededMsg{DeviceId: deviceId, RFAddr: rfAddr, NewKey: newKey})
}

func (c *Conn) GetRFAddr(id uint32) error {
	return c.Write(GetRFAddrMsg{DeviceId: id})
}

func (c *Conn) RFAddr(id uint32, addr uint16) error {
	return c.Write(RFAddrMsg{DeviceId: id, RFAddr: addr})
}

/* Types */
func (c *Conn) GetType(id uint32) error {
	return c.Write(GetTypeMsg{TypeId: id})
}

func (c *Conn) SendNoType(id uint32) error {
	return c.Write(NoTypeMsg{TypeId: id})
}

func (c *Conn) Type(id uint32, isVirtual bool, ttl uint32, src string, masks []uint32, intervals []uint8) error {
	return c.Write(TypeMsg{
		TypeId:          id,
		IsVirtual:       isVirtual,
		RefreshInterval: ttl,
		Src:             src,
		Masks:           masks,
		Intervals:       intervals,
	})
}

/* JOINING */
func (c *Conn) RequestJoin(UUID string, typeId uint32) error {
	return c.Write(JoinRequestMsg{
		TypeId: typeId,
		UUID:   UUID,
	})
}

func (c *Conn) JoinRequestPending(UUID string) error {
	return c.Write(JoinRequestPendingMsg{UUID: UUID})
}

func (c *Conn) JoinRequestDeclined(UUID string) error {
	return c.Write(JoinRequestDeclinedMsg{UUID: UUID})
}

func (c *Conn) JoinRequestApproved(UUID string, id uint32, addr uint16) error {
	return c.Write(JoinRequestApprovedMsg{UUID: UUID, DeviceId: id, RFAddr: addr})
}

/* WIRED SENSORS */
func (c *Conn) ConnectDevice(UUID string, typeId uint32, port uint16, parentId uint32) error {
	return c.Write(ConnectDeviceMsg{UUID: UUID, TypeId: typeId, ParentDeviceId: parentId, Port: port})
}

func (c *Conn) MoveDevice(deviceId, parentId uint32, port uint16) error {
	return c.Write(MoveDeviceMsg{DeviceId: deviceId, ParentDeviceId: parentId, Port: port})
	
}

func (c *Conn) DisconnectDevice(deviceId, parentId uint32) error {
	return c.Write(DisconnectDeviceMsg{DeviceId: deviceId, ParentDeviceId: parentId})
}

func (c *Conn) ConnectDeviceDeclined(UUID string, parentId uint32) error {
	return c.Write(ConnectDeviceDeclinedMsg{UUID: UUID, ParentDeviceId: parentId})
}

func (c *Conn) ConnectDeviceApproved(UUID string, parentId, id uint32) error {
	return c.Write(ConnectDeviceApprovedMsg{UUID: UUID, ParentDeviceId: parentId, DeviceId: id})
}

/* Utilities */

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func NewConn(conn net.Conn) (*Conn, error) {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	return &Conn{conn: conn, encoder: enc, decoder: dec}, nil
}

var UnauthorizedErr error = errors.New("Authentication required")
