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

type RequestHandler interface {
	Hello(c *Conn, protoVersion uint32) uint32
	Authenticate(c *Conn, UUID, key string) (uint32, uint16, string, error)
	GetRFAddr(c *Conn, id uint32) (uint16, error)
	JoinRequest(c *Conn, UUID string, typeId uint32) (JoinRequestStatus, uint32, uint16, error)
}

func (c *Conn) HandleRequests(rh RequestHandler) {
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
			
			id, addr, key, err := rh.Authenticate(c, v.UUID, v.Key)
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
		}
	}
}
