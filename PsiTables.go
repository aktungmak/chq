package main

import (
	"errors"
	"fmt"
)

// Represents a MPEG TS Program Association Table
// see ISO 13818-1 section 2.4.4.3
type Pat struct {
	Tid  byte `mpeg:"table_id"`
	Ssi  bool `mpeg:"section_syntax_indicator"`
	Sl   int  `mpeg:"section_length"`
	Tsid int  `mpeg:"transport_stream_id"`
	Vn   int  `mpeg:"version_number"`
	Cni  bool `mpeg:"current_next_indicator"`
	Sn   byte `mpeg:"section_number"`
	Lsn  byte `mpeg:"last_section_number"`
	Pgms []pgm
	Crc  uint32 `mpeg:"CRC_32"`
}

type pgm struct {
	Pn    int   `mpeg:"program_number"`
	Pmpid int16 `mpeg:"program_map_PID"`
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
