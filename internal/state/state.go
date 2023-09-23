package state

type Type int

type State interface {
	Type() Type
}
