package service

import (
	"io"
	"log"
)

type JoinRequestStatus int

const (
	JOIN_REQUEST_DECLINED JoinRequestStatus = iota
	JOIN_REQUEST_PENDING
	JOIN_REQUEST_APPROVED
)

type Server interface {
	Hello(c *Conn, clientProtoVer uint32) (serverProtoVer uint32)

	Authenticate(c *Conn, uuid string, typeId uint32, key string) (deviceId uint32, rfAddr uint16, newKey string, err error)

	GetRFAddr(c *Conn, id uint32) (rfAddr uint16, err error)

	JoinRequest(c *Conn, uuid string, typeId uint32) (status JoinRequestStatus, deviceId uint32, rfAddr uint16, err error)

	ConnectDevice(c *Conn, parentId uint32, uuid string, typeId uint32, port uint16) (id uint32, err error)
	DisconnectDevice(c *Conn, parentId, id uint32) error
	MoveDevice(c *Conn, parentId, id uint32, port uint16) error

	GetType(c *Conn, typeId uint32) (ttl uint32, isVirtual bool, src string, masks []uint32, intervals []byte, err error)

	OnClose()
}

func (c *Conn) ServerHandler(rh Server) {
	defer func() {
		rh.OnClose()
		c.Close()
	}()

	for {
		msg, err := c.Read()
		if err != nil {
			if err == io.EOF {
				/* Exit from loop */
				break
			}

			log.Println(err)
			continue
		}

		switch msg.(type) {

		case HelloMsg:
			v := msg.(HelloMsg)
			pv := rh.Hello(c, v.ProtocolVersion)
			err := c.Hello(pv)
			if err != nil {
				log.Println(err)
			}

		case AuthenticateMsg:
			v := msg.(AuthenticateMsg)

			id, addr, key, err := rh.Authenticate(c, v.UUID, v.TypeId, v.Key)
			if err != nil {
				err = c.AuthenticationFailed(err.Error())
			} else {
				err = c.AuthenticationSucceeded(id, addr, key)
			}

			if err != nil {
				log.Println(err)
			}

		case GetRFAddrMsg:
			v := msg.(GetRFAddrMsg)

			addr, err := rh.GetRFAddr(c, v.DeviceId)
			if err == UnauthorizedErr {
				c.AuthenticationFailed(err.Error())
				continue
			}

			if addr != 0 {
				err := c.RFAddr(v.DeviceId, addr)
				if err != nil {
					log.Println(err)
				}
			}

		case JoinRequestMsg:
			v := msg.(JoinRequestMsg)

			status, id, addr, err := rh.JoinRequest(c, v.UUID, v.TypeId)
			if err == UnauthorizedErr {
				err = c.AuthenticationFailed(err.Error())
				continue
			}

			if status == JOIN_REQUEST_APPROVED {
				err = c.JoinRequestApproved(v.UUID, id, addr)
			} else if status == JOIN_REQUEST_PENDING {
				err = c.JoinRequestPending(v.UUID)
			} else if status == JOIN_REQUEST_DECLINED {
				err = c.JoinRequestDeclined(v.UUID)
			} else {
				err = c.JoinRequestPending(v.UUID)
			}

			if err != nil {
				log.Println(err)
			}

		case ConnectDeviceMsg:
			v := msg.(ConnectDeviceMsg)

			id, err := rh.ConnectDevice(c, v.ParentDeviceId, v.UUID, v.TypeId, v.Port)
			if err == UnauthorizedErr {
				err = c.AuthenticationFailed(err.Error())
				continue
			}

			if err != nil {
				log.Println(err)
				continue
			}

			if id > 0 {
				err = c.ConnectDeviceApproved(v.UUID, v.ParentDeviceId, id)
			} else {
				err = c.ConnectDeviceDeclined(v.UUID, v.ParentDeviceId)
			}

			if err != nil {
				log.Println(err)
			}

		case MoveDeviceMsg:
			v := msg.(MoveDeviceMsg)

			err = rh.MoveDevice(c, v.ParentDeviceId, v.DeviceId, v.Port)
			if err == UnauthorizedErr {
				err = c.AuthenticationFailed(err.Error())
				continue
			}

			if err != nil {
				log.Println(err)
				continue
			}

		case DisconnectDeviceMsg:
			v := msg.(DisconnectDeviceMsg)

			err = rh.DisconnectDevice(c, v.ParentDeviceId, v.DeviceId)
			if err == UnauthorizedErr {
				err = c.AuthenticationFailed(err.Error())
				continue
			}

			if err != nil {
				log.Println(err)
				continue
			}

		case GetTypeMsg:
			v := msg.(GetTypeMsg)

			ttl, virtual, src, masks, intervals, err := rh.GetType(c, v.TypeId)
			if err != nil {
				log.Println(err)
				continue
			}

			if ttl != 0 {
				c.Type(v.TypeId, virtual, ttl, src, masks, intervals)

			}
		}

	}
}
