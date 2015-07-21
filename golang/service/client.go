package service

import (
	"io"
	"log"
)

type Client interface {
	Hello(c *Conn, serverProtoVer uint32)

	AuthenticationFailed(c *Conn, reason string)
	AuthenticationSucceeded(c *Conn, deviceId uint32, rfAddr uint16, newKey string)

	JoinRequestApproved(c *Conn, UUID string, deviceId uint32, rfAddr uint16)
	JoinRequestPending(c *Conn, UUID string)
	JoinRequestDeclined(c *Conn, UUID string)

	OnClose()
}

func (c *Conn) ClientHandler(rh Client) {
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
			rh.Hello(c, v.ProtocolVersion)

		case AuthenticationSucceededMsg:
			v := msg.(AuthenticationSucceededMsg)
			rh.AuthenticationSucceeded(c, v.DeviceId, v.RFAddr, v.NewKey)

		case AuthenticationFailedMsg:
			v := msg.(AuthenticationFailedMsg)
			rh.AuthenticationFailed(c, v.Reason)

		case JoinRequestApprovedMsg:
			v := msg.(JoinRequestApprovedMsg)
			rh.JoinRequestApproved(c, v.UUID, v.DeviceId, v.RFAddr)

		case JoinRequestPendingMsg:
			v := msg.(JoinRequestPendingMsg)
			rh.JoinRequestPending(c, v.UUID)

		case JoinRequestDeclinedMsg:
			v := msg.(JoinRequestDeclinedMsg)
			rh.JoinRequestDeclined(c, v.UUID)
		}

	}

}
