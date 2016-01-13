package main

// generic representation of a descriptor.
type descriptor interface {
	Tag() int
	Len() int
}

func ParseDescriptors(data []byte) []descriptor {
	descriptorConstructors := map[byte]func([]byte) descriptor{
		0x42: NewServiceDescriptor,
	}
	// todo handle multiple descriptors
	ctr, ok := descriptorConstructors[data[0]]
	if ok {
		return ctr(data)
	} else {
		return nil
	}
}

type ServiceDescriptor struct {
	Dt   byte
	Dl   byte
	St   byte
	Spnl byte
	Spn  string
	Snl  byte
	Sn   string
}

func NewServiceDescriptor(data []byte) descriptor {
	sd := ServiceDescriptor{}

	return sd
}

func (d ServiceDescriptor) Tag() int {
	return int(d.Dt)
}
func (d ServiceDescriptor) Len() int {
	return int(d.Dl)
}
