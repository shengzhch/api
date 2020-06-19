package base

type Packet interface {
	String() string
	Bytes() []byte
	Empty() bool
	Type() int
}


