package client

type Feature interface {
	Name() string
	Min() int
	Max() int
	Step() int
	Exists() bool
	Update(string) error
	OnSet() (chan string, error)
	OnSetFunc(func(string)) error
}
