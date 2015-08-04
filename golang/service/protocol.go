package service

import (
	"encoding/gob"
)

/* These are the messages that can be passed between client and server. All messages are asynchronos */

type ServerMessage struct {
	Data interface{}
}

/*
	Start
*/
type HelloMsg struct {
	ProtocolVersion uint32
}

/*
	Authentication connects a gateway to our service
*/
type AuthenticateMsg struct {
	TypeId uint32
	UUID   string
	Key    string
}

type AuthenticationFailedMsg struct {
	Reason string
}

type AuthenticationSucceededMsg struct {
	DeviceId uint32
	NewKey   string
	RFAddr   uint16
}

/*
	A Join request is a sent when a wireless device wants to join our setup

	Client sends JoinRequest. Server sends JoinRequestApproved, Declined or Pending.
*/
type JoinRequestMsg struct {
	TypeId uint32
	UUID   string
}

type JoinRequestApprovedMsg struct {
	UUID     string
	DeviceId uint32
	RFAddr   uint16
}

type JoinRequestDeclinedMsg struct {
	UUID string
}

type JoinRequestPendingMsg struct {
	UUID string
}

/*
	Time
*/
type GetTimeMsg struct {
	MyTime uint32
}

type TimeMsg struct {
	ServerTime uint32
}

/*
	Wireless mgt
*/
type GetRFAddrMsg struct {
	DeviceId uint32
}

type RFAddrMsg struct {
	DeviceId uint32
	RFAddr   uint16
}

/*
	Wired devices establish a parent <-> child
*/
type ConnectDeviceMsg struct {
	ParentDeviceId uint32
	TypeId         uint32
	UUID           string
	Port           uint16
}

type ConnectDeviceApprovedMsg struct {
	ParentDeviceId uint32
	UUID           string
	DeviceId       uint32
}

type ConnectDeviceDeclinedMsg struct {
	ParentDeviceId uint32
	UUID           string
}

type MoveDeviceMsg struct {
	ParentDeviceId uint32
	DeviceId       uint32
	Port           uint16
}

type DisconnectDeviceMsg struct {
	ParentDeviceId uint32
	DeviceId       uint32
}

type GetDevices struct {
	ParentDeviceId uint32
	TypeId         *uint32
	Virtual        *bool
}

type DeviceInfo struct {
	ParentDeviceId uint32
	DeviceId       uint32
	TypeId         uint32
	UUID           string
}

/* Type information */
type GetTypeMsg struct {
	TypeId uint32
}

type TypeMsg struct {
	TypeId          uint32
	IsVirtual       bool
	RefreshInterval uint32
	Src             string
	Masks           []uint32
	Intervals       []uint8
}

type NoTypeMsg struct {
	TypeId uint32
}

/* Sensor data */
type RawSensorData struct {
	Sensors uint32
	Values 	[]uint16
}

type SensorDataMsg struct {
	DeviceId  uint32
	TimeStamp uint32
	Values    map[string]float32
	RawData	  *RawSensorData
}

type SensorDataSyncedMsg struct {
}

/* Register all types */
func init() {
	gob.RegisterName("service.HelloMsg", HelloMsg{})
	gob.RegisterName("service.AuthenticateMsg", AuthenticateMsg{})
	gob.RegisterName("service.AuthenticationFailedMsg", AuthenticationFailedMsg{})
	gob.RegisterName("service.AuthenticationSucceededMsg", AuthenticationSucceededMsg{})
	gob.RegisterName("service.ConnectDeviceApprovedMsg", ConnectDeviceApprovedMsg{})
	gob.RegisterName("service.ConnectDeviceDeclinedMsg", ConnectDeviceDeclinedMsg{})
	gob.RegisterName("service.ConnectDeviceMsg", ConnectDeviceMsg{})
	gob.RegisterName("service.DisconnectDeviceMsg", DisconnectDeviceMsg{})
	gob.RegisterName("service.GetRFAddrMsg", GetRFAddrMsg{})
	gob.RegisterName("service.GetTimeMsg", GetTimeMsg{})
	gob.RegisterName("service.GetTypeMsg", GetTypeMsg{})
	gob.RegisterName("service.HelloMsg", HelloMsg{})
	gob.RegisterName("service.JoinRequestApprovedMsg", JoinRequestApprovedMsg{})
	gob.RegisterName("service.JoinRequestDeclinedMsg", JoinRequestDeclinedMsg{})
	gob.RegisterName("service.JoinRequestMsg", JoinRequestMsg{})
	gob.RegisterName("service.JoinRequestPendingMsg", JoinRequestPendingMsg{})
	gob.RegisterName("service.MoveDeviceMsg", MoveDeviceMsg{})
	gob.RegisterName("service.NoTypeMsg", NoTypeMsg{})
	gob.RegisterName("service.RFAddrMsg", RFAddrMsg{})
	gob.RegisterName("service.RawSensorData", RawSensorData{})
	gob.RegisterName("service.SensorDataMsg", SensorDataMsg{})
	gob.RegisterName("service.SensorDataSyncedMsg", SensorDataSyncedMsg{})
	gob.RegisterName("service.ServerMessage", ServerMessage{})
	gob.RegisterName("service.TimeMsg", TimeMsg{})
	gob.RegisterName("service.TypeMsg", TypeMsg{})
}
