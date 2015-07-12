package main

type TsPacket struct {
	Header          TsPacketHeader
	AdaptationField AdaptationField
	Payload         []byte
}

func NewTsPacket(data []byte) TsPacket {
	// TODO parse data and populate packet
	return TsPacket{}
}

type TsPacketHeader struct{}

type AdaptationField struct{}
