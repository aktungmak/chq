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

	ret := make([]descriptor, 0)
	i := 0
	for {
		ctr, ok := descriptorConstructors[data[i]]
		if !ok {
			break
		}
		desc := ctr(data[0:])
		ret = append(ret, desc)
		i += desc.Len()
	}
	return ret
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
