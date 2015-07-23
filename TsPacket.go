package main

type TsPacket struct {
	Header          TsPacketHeader
	AdaptationField AdaptationField
	Payload         []byte
}

func NewTsPacket(data []byte) TsPacket {
	// TODO parse data and populate packet
	return TsPacket{
		Header:  NewTsPacketHeader(data),
		Payload: data[4:188],
	}

}

// Represents the 4-byte transport stream packet header
// See ISO 13818-1 Table 2-2
type TsPacketHeader struct {
	SyncByte byte
	Tei      bool
	Pusi     bool
	Tp       bool
	Pid      Pid
	Tsc      byte
	Afc      byte
	Cc       byte
}

// Constructor to create a new header struct
// expects data to start with a sync byte
func NewTsPacketHeader(data []byte) TsPacketHeader {
	return TsPacketHeader{
		SyncByte: data[0],
		Tei:      data[1]&128 != 0,
		Pusi:     data[1]&64 != 0,
		Tp:       data[1]&32 != 0,
		Pid:      ((Pid(data[1]) & 31) << 8) + Pid(data[2]),
		Tsc:      (data[3] & 192) >> 6,
		Afc:      (data[3] & 48) >> 4,
		Cc:       data[3] & 15,
	}
}

type AdaptationField struct{}

type Pid int16
