package main

import (
	"errors"
	"fmt"
)

// represents a sCTE35 splice_info_section
// see sCTE35-2013 table 7-1
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
	*scte35SpliceNull           `json:"splice_null,omitempty"`
	*scte35SpliceSchedule       `json:"splice_schedule,omitempty"`
	*scte35SpliceInsert         `json:"splice_insert,omitempty"`
	*scte35TimeSignal           `json:"time_signal,omitempty"`
	*scte35BandwidthReservation `json:"bandwidth_reservation,omitempty"`
	*scte35PrivateCommand       `json:"private_command,omitempty"`
	Dll                         int `json:"descriptor_loop_length"`
	Descriptors                 []descriptor
	Ecrc                        uint32 `json:"E_CRC_32"`
	Crc                         uint32 `json:"CRC_32"`
}

func NewScte35SpliceInfo(data []byte) (*Scte35SpliceInfo, error) {
	var err error
	sis := &Scte35SpliceInfo{}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Malformed SCTE35 splice_info_section")
		}
	}()

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
		sis.scte35SpliceNull, _ = newScte35SpliceNull(data[14:])
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
		return sis, errors.New("Unsupported SCTE35 reserved splice command type!")
	}

	end := sis.Sl + 3
	if sis.Ep {
		sis.Ecrc = uint32(data[end-4]) + uint32(data[end-5])<<8 + uint32(data[end-6])<<16 + uint32(data[end-7])<<24
	}
	sis.Crc = uint32(data[end]) + uint32(data[end-1])<<8 + uint32(data[end-2])<<16 + uint32(data[end-3])<<24

	return sis, nil
}

type scte35SpliceNull struct{}

func newScte35SpliceNull(data []byte) (*scte35SpliceNull, error) {
	return &scte35SpliceNull{}, nil
}

type scte35SpliceSchedule struct {
	Sc   byte
	Sins []*scte35SpliceInsert
}

func newScte35SpliceSchedule(data []byte) (*scte35SpliceSchedule, error) {
	ss := &scte35SpliceSchedule{}
	if len(data) < 1 {
		return ss, errors.New("SCTE35 splice_schedule command data is too short to parse!")
	}
	ss.Sc = data[0]
	for i := 1; i < len(data); i += 1 {
		// si, err := newscte35SpliceInsert(data[i:])
	}
	return ss, nil

}

type scte35SpliceInsert struct {
	Seid        uint32 `json:"splice_event_id"`
	Seci        bool   `json:"splice_event_cancel_indicator"`
	Ooni        bool   `json:"out_of_network_indicator"`
	Psf         bool   `json:"program_splice_flag"`
	Df          bool   `json:"duration_flag"`
	Sif         bool   `json:"splice_immediate_flag"`
	*spliceTime `json:"splice_time,omitempty"`
	Cc          int `json:"component_count"`
	// TODO component splices
	*breakDuration `json:"break_duration,omitempty"`
	Upid           int  `json:"unique_program_id"`
	An             byte `json:"avail_num"`
	Ae             byte `json:"avails_expected"`
}

func newScte35SpliceInsert(data []byte) (*scte35SpliceInsert, error) {
	si := &scte35SpliceInsert{}
	si.Seid = (uint32(data[0]) << 24) + (uint32(data[1]) << 16) + (uint32(data[2]) << 8) + uint32(data[3])
	si.Seci = data[4]&128 != 0
	si.Ooni = data[5]&128 != 0
	si.Psf = data[5]&64 != 0
	si.Df = data[5]&32 != 0
	si.Sif = data[5]&16 != 0

	// keep track of option field length
	i := 6

	if si.Psf && si.Sif {
		st, err := newSpliceTime(data[i:])
		if err != nil {
			return si, err
		}
		if st.Tsf {
			i += 5
		} else {
			i += 1
		}
	}

	if si.Psf {
		si.Cc = int(data[i])
		cs := make([]*elementaryPidData, 0)
		for j := 0; j < si.CC; j++ {
			epd := *elementaryPidData{}
			epd.Ct = data[i]
			i += 1
			if !si.Sif {
				// epd.spliceTime =
			}
			cs = append(cs, epd)
		}
	}

	if si.Df {
		bd, err := newBreakDuration(data[i:])
		if err != nil {
			return si, err
		}
		i += 5
	}

	si.Upid = (int(data[i]) << 8) + int(data[i+1])
	si.An = data[i+2]
	si.Ae = data[i+3]
	return si, nil
}

type scte35TimeSignal struct {
	*spliceTime
}

func newScte35TimeSignal(data []byte) (*scte35TimeSignal, error) {
	return &scte35TimeSignal{}, nil
}

type scte35BandwidthReservation struct{}

func newScte35BandwidthReservation(data []byte) (*scte35BandwidthReservation, error) {
	return &scte35BandwidthReservation{}, nil
}

type scte35PrivateCommand struct {
	Id uint32 `json:"identifier"`
	Pb []byte `json:"private_byte"`
}

func newScte35PrivateCommand(data []byte) (*scte35PrivateCommand, error) {
	return &scte35PrivateCommand{}, nil
}

type spliceTime struct {
	Tsf  bool  `json:"time_specified_flag"`
	Ptst int64 `json:"pts_time"`
}

func newSpliceTime(data []byte) (*spliceTime, error) {
	return &spliceTime{}, nil
}

type breakDuration struct {
	Ar bool  `json:"auto_return"`
	D  int64 `json:"duration"`
}

func newBreakDuration(data []byte) (*breakDuration, error) {
	return &breakDuration{}, nil
}

type elementaryPidData struct {
	Ct          byte `json:"component_tag"`
	*spliceTime `json:"splice_time,omitempty"`
}
