package dao

// Driver allows go-opc-server to output opc instructions
// to different kinds of consumers/domains
type Driver interface {
	GetInstance() interface{}
	Setup(interface{}) error
	Display(channel int, pixels [][]int)
	Unregister()
}