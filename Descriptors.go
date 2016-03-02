package main

// generic representation of a descriptor.
type descriptor interface {
	Tag() int
	Len() int
}

// parse a byte slice of descriptors into a slice of
// descriptor structs. ignores unknown types.
func ParseDescriptors(data []byte) []descriptor {
	descriptorConstructors := map[byte]func([]byte) descriptor{
		0x48: NewServiceDescriptor,
		// new descriptor parsers register here
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

	Dt   byte   `json:"descriptor_tag"`
	Dl   byte   `json:"descriptor_length"`
	St   byte   `json:"service_type"`
	Spnl byte   `json:"service_provider_name_length"`
	Spn  string `json:"service_provider_name"`
	Snl  byte   `json:"service_name_length"`
	Sn   string `json:"service_name"`
}

func NewServiceDescriptor(data []byte) descriptor {
	sd := ServiceDescriptor{}
	sd.Dt = data[0]
	sd.Dl = data[1]
	sd.St = data[2]
	sd.Spnl = data[3]
	sd.Spn = string(data[4 : 4+sd.Spnl])
	sd.Snl = data[4+sd.Spnl]
	sd.Sn = string(data[5+sd.Spnl : 5+sd.Spnl+sd.Snl])
	return sd
}

func (d ServiceDescriptor) Tag() int {
	return int(d.Dt)
}
func (d ServiceDescriptor) Len() int {
	return int(d.Dl)
}
