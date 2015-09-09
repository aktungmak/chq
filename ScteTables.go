package main

import (
	"errors"
	"fmt"
)

// represents a SCTE35 splice_info_section
// see SCTE35-2013 table 7-1
type Scte35SpliceInfo struct {
	Tid                         byte  `json:"table_id"`
	Ssi                         bool  `json:"section_syntax_indicator"`
	Pi                          bool  `json:"private_indicator"`
	Sl                          int   `json:"section_length"`
	Pv                          byte  `json:"protocol_version"`
	Ep                          bool  `json:"encrypted_packet"`
	Ea                          int   `json:"encryption_algorithm"`
	Ptsa                        int64 `json:"pts_adjustment"`
	Cwi                         byte  `json:"cw_index"`
	Tier                        int   `json:"tier"`
	Scl                         int   `json:"splice_command_length"`
	Sct                         byte  `json:"splice_command_type"`
	*Scte35SpliceNull           `json:"splice_null,omitempty"`
	*Scte35SpliceSchedule       `json:"splice_schedule,omitempty"`
	*Scte35SpliceInsert         `json:"splice_insert,omitempty"`
	*Scte35TimeSignal           `json:"time_signal,omitempty"`
	*Scte35BandwidthReservation `json:"bandwidth_reservation,omitempty"`
	*Scte35PrivateCommand       `json:"private_command,omitempty"`
	Dll                         int `json:"descriptor_loop_length"`
	Descriptors                 []descriptor
	Ecrc                        uint32 `json:"E_CRC_32"`
	Crc                         uint32 `json:"CRC_32"`
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
		sis.Scte35SpliceNull, _ = NewScte35SpliceNull(data[14:])
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
		return sis, errors.New("SCTE35 reserved splice command type!")
	}

	end := len(data) - 1
	if sis.Ep {
		sis.Ecrc = uint32(data[end-4]) + uint32(data[end-5])<<8 + uint32(data[end-6])<<16 + uint32(data[end-7])<<24
	}
	sis.Crc = uint32(data[end]) + uint32(data[end-1])<<8 + uint32(data[end-2])<<16 + uint32(data[end-3])<<24

	return sis, nil
}

type Scte35SpliceNull struct{}

type Scte35SpliceSchedule struct {
	Sc   byte
	Sins []*Scte35SpliceInsert
}

type Scte35SpliceInsert struct {
	Seid uint32 `json:"splice_event_id"`
	Seci bool   `json:"splice_event_cancel_indicator"`
	Ooni bool   `json:"out_of_network_indicator"`
	Psf  bool   `json:"program_splice_flag"`
	Df   bool   `json:"duration_flag"`
	Sif  bool   `json:"splice_immediate_flag"`
	*spliceTime
	Cc int `json:"component_count"`
	// TODO component splices
	*breakDuration
	Upid int  `json:"unique_program_id"`
	An   byte `json:"avail_num"`
	Ae   byte `json:"avails_expected"`
}

type Scte35TimeSignal struct {
	*spliceTime
}

type Scte35BandwidthReservation struct{}

type Scte35PrivateCommand struct {
	Id uint32 `json:"identifier"`
	Pb []byte `json:"private_byte"`
}

type spliceTime struct {
	Tsf  bool  `json:"time_specified_flag"`
	Ptst int64 `json:"pts_time"`
}

type breakDuration struct {
	Ar bool  `json:"auto_return"`
	D  int64 `json:"duration"`
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
	for i := 1; i < len(data); i += 1 {
		// si, err := NewScte35SpliceInsert(data[i:])
	}
	return ss, nil

}
func NewScte35SpliceInsert(data []byte) (*Scte35SpliceInsert, error) {
	si := &Scte35SpliceInsert{}

	return si, nil

}
func NewScte35TimeSignal(data []byte) (*Scte35TimeSignal, error) {
	return &Scte35TimeSignal{}, nil
}
func NewScte35BandwidthReservation(data []byte) (*Scte35BandwidthReservation, error) {
	return &Scte35BandwidthReservation{}, nil
}
func NewScte35PrivateCommand(data []byte) (*Scte35PrivateCommand, error) {
	return &Scte35PrivateCommand{}, nil
}
