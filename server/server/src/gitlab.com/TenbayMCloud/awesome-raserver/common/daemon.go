package common

type Daemon interface {
	Initialize() error
	Run()
	Name() string
}

