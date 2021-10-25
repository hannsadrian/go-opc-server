package dao

type Driver interface {
	GetInstance() interface{}
	Setup() error
	Display(channel int, pixels [][]int)
	Unregister()
}