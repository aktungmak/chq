package main

import (
	"errors"
	"fmt"
)

// represents a SCTE35 splice_info_section
// see SCTE35-2013 table 7-1
type Scte35SpliceInfo struct {
	Tid         byte  `json:"table_id"`
	Ssi         bool  `json:"section_syntax_indicator"`
	Pi          bool  `json:"private_indicator"`
	Sl          int   `json:"section_length"`
	Pv          byte  `json:"protocol_version"`
	Ep          bool  `json:"encrypted_packet"`
	Ea          int   `json:"encryption_algorithm"`
	Ptsa        int64 `json:"pts_adjustment"`
	Cwi         byte  `json:"cw_index"`
	Tier        int   `json:"tier"`
	Scl         int   `json:"splice_command_length"`
	Sct         byte  `json:"splice_command_type"`
	Sc          Scte35SpliceCommand
	Dll         int `json:"descriptor_loop_length"`
	Descriptors []descriptor
	Ecrc        uint32 `json:"E_CRC_32"`
	Crc         uint32 `json:"CRC_32"`
}

func NewScte35SpliceInfo(data []byte) (*Scte35SpliceInfo, error) {
	sis := &Scte35SpliceInfo{}
	if len(data) < 8 {
		return sis, errors.New("SCTE35 data too short to parse!")
	}
	sis.Tid = data[0]
	sis.Ssi = data[1]&128 != 0
	sis.Pi = data[1]&64 != 0
	sis.Sl = (int(data[1]&15) << 8) + int(data[2])

	if len(data) < sis.Sl+3 {
		return sis, errors.New("SCTE35 data length and section_length field mismatch!")
	}

	sis.Pv = data[3]
	sis.Ep = data[4]&128 != 0
	sis.Ea = int(data[4]&126) >> 1
	sis.Ptsa = (int64(data[4]&126) << 32) +
		(int64(data[5]) << 24) +
		(int64(data[6]) << 16) +
		(int64(data[7]) << 8) +
		(int64(data[8]))
	sis.Cwi = data[9]
	sis.Tier = (int(data[10]) << 4) + (int(data[11]&240) >> 4)
	sis.Scl = (int(data[11]&15) << 8) + int(data[12]) // this is unreliable as it is usually 0xfff
	sis.Sct = data[13]

	switch sis.Sct {
	case 0x00:
		// splice_null
		fmt.Printf("splice_null\n")
	case 0x04:
		// splice_schedule
		fmt.Printf("splice_schedule\n")
	case 0x05:
		// splice_insert
		fmt.Printf("splice_insert\n")
	case 0x06:
		// time_signal
		fmt.Printf("time_signal\n")
	case 0x07:
		// bandwidth_reservation
		fmt.Printf("bandwidth_reservation\n")
	case 0xff:
		// private_command
		fmt.Printf("private_command\n")
	default:
		// reserved
	}

	return sis, nil
}

// this represents all splice command types. need a better way of doing this...
type Scte35SpliceCommand interface{}

type Scte35SpliceNull struct {
}

type Scte35SpliceSchedule struct {
	Sc byte
	[]*Scte35SpliceInsert
}

type Scte35SpliceInsert struct {

}

type Scte35TimeSignal struct {
}

type Scte35BandwidthReservation struct {
}

type Scte35PrivateCommand struct {
}

type Scte35SpliceNull struct {
}

type Scte35SpliceSchedule struct {
}

type Scte35SpliceInsert struct {
}

type Scte35TimeSignal struct {
}

type Scte35BandwidthReservation struct {
}

type Scte35PrivateCommand struct {
}

func NewScte35SpliceNull(data []byte) (*Scte35SpliceNull, error) {
	return &Scte35SpliceNull{}, nil
}
func NewScte35SpliceSchedule(data []byte) (*Scte35SpliceSchedule, error) {
	ss := &Scte35SpliceSchedule{}
	if len(data) < 1 {
		return ss, errors.New("SCTE35 splice_schedule command data is too short to parse!")
	}
	ss.Sc = data[0]
	for i := 1, i < len(data); i +=  {
		si := NewScte35SpliceInsert(data[i:])
	}

}
func NewScte35SpliceInsert(data []byte) (*Scte35SpliceInsert, error) {

}
func NewScte35TimeSignal(data []byte) (*Scte35TimeSignal, error) {

}
func NewScte35BandwidthReservation(data []byte) (*Scte35BandwidthReservation, error) {

}
func NewScte35PrivateCommand(data []byte) (*Scte35PrivateCommand, error) {

}
func NewScte35SpliceNull(data []byte) (*Scte35SpliceNull, error) {

}
func NewScte35SpliceSchedule(data []byte) (*Scte35SpliceSchedule, error) {

}
func NewScte35SpliceInsert(data []byte) (*Scte35SpliceInsert, error) {

}
func NewScte35TimeSignal(data []byte) (*Scte35TimeSignal, error) {

}
func NewScte35BandwidthReservation(data []byte) (*Scte35BandwidthReservation, error) {

}
func NewScte35PrivateCommand(data []byte) (*Scte35PrivateCommand, error) {

}
