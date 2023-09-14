package models

const (
	RegisterUserType = "registerUser"
	LoginUserType    = "loginUser"
)

type MQMessage struct {
	Type string
	Body interface{}
}
