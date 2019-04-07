package server

type Feature interface {
	Name() string
	Min() int
	Max() int
	Step() int
	Exists() bool
	Set(string) error
	OnUpdate() (chan string, error)
	OnUpdateFunc(func(string)) error
	Value() string
}
