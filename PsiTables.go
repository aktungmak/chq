package main

import (
	"errors"
	"fmt"
)

// generic representation of a descriptor. TODO implement it!
type descriptor struct{}

// Represents a MPEG TS Program Association Table
// see ISO 13818-1 section 2.4.4.3
type Pat struct {
	Tid  byte `json:"table_id"`
	Ssi  bool `json:"section_syntax_indicator"`
	Sl   int  `json:"section_length"`
	Tsid int  `json:"transport_stream_id"`
	Vn   int  `json:"version_number"`
	Cni  bool `json:"current_next_indicator"`
	Sn   byte `json:"section_number"`
	Lsn  byte `json:"last_section_number"`
	Pgms []pgm
	Crc  uint32 `json:"CRC_32"`
}

type pgm struct {
	Pn    int   `json:"program_number"`
	Pmpid int16 `json:"program_map_PID"`
}

func NewPat(data []byte) (*Pat, error) {
	pat := &Pat{}
	if len(data) < 8 {
		return pat, errors.New("PAT data too short!")
	}
	pat.Tid = data[0]
	pat.Ssi = data[1]&128 != 0
	pat.Sl = (int(data[1]&15) << 8) + int(data[2])
	pat.Tsid = (int(data[3]) << 8) + int(data[4])
	pat.Vn = int(data[5]&62) >> 1
	pat.Cni = data[5]&1 != 0
	pat.Sn = data[6]
	pat.Lsn = data[7]

	if len(data) < pat.Sl+3 {
		return pat, errors.New("PAT data length and section_length field mismatch!")
	}

	pat.Pgms = make([]pgm, 0)
	i := 8
	for ; i < pat.Sl-4; i += 4 {
		pn := (int(data[i]) << 8) + int(data[i+1])
		pid := ((int16(data[i+2]) & 31) << 8) + int16(data[i+3])
		pat.Pgms = append(pat.Pgms, pgm{pn, pid})
	}

	pat.Crc = uint32(data[i+3]) + uint32(data[i+2])<<8 + uint32(data[i+1])<<16 + uint32(data[i])<<24

	if crc := CalculateCrc32(data[:i]); crc != pat.Crc {
		return pat, errors.New(fmt.Sprintf("CRC error! Calculated %x but PAT says %x", crc, pat.Crc))
	}

	return pat, nil
}

// Represents a Programme Map Table
// See ISO 13818-1 table 2-28
type Pmt struct {
	Tid         byte  `json:table_id`
	Ssi         bool  `json:section_syntax_indicator`
	Sl          int   `json:section_length`
	Pn          int   `json:program_number`
	Vn          int   `json:version_number`
	Cni         bool  `json:current_next_indicator`
	Sn          byte  `json:section_number`
	Lsn         byte  `json:last_section_number`
	Pcrpid      int16 `json:PCR_PID`
	Pil         int   `json:program_info_length`
	Descriptors []descriptor
	Pels        []pel
	Crc         uint32 `json:"CRC_32"`
}
type pel struct {
	St          byte  `json:stream_type`
	Pid         int16 `json:elementary_PID`
	Eil         int   `json:ES_info_length`
	Descriptors []descriptor
}

func NewPmt(data []byte) (*Pmt, error) {
	pmt := &Pmt{}
	return pmt, nil
}
