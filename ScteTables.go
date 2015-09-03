package main

import (
	"errors"
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

type Scte35SpliceCommand interface{}

func NewScte35SpliceInfo(data []byte) (*Scte35SpliceInfo, error) {
	sis := &Scte35SpliceInfo{}
	if len(data) < 8 {
		return sis, errors.New("SCTE35 data too short to parse!")
	}
	sis.Tid = data[0]
	sis.Ssi = data[1]&128 != 0
	sis.Pi = data[1]&64 != 0
	sis.Sl = (int(data[1]&15) << 8) + int(data[2])
	sis.Pv = data[3]
	sis.Ep = data[4]&128 != 0
	sis.Ea = int(data[4]&126) >> 1
	sis.Ptsa = (int64(data[4]&126) << 32) +
		(int64(data[5]) << 24) +
		(int64(data[6]) << 16) +
		(int64(data[7]) << 8) +
		(int64(data[8]))

	return sis, nil
}
