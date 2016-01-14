package main

import (
	"errors"
	"fmt"
)

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
	Pn    int `json:"program_number"`
	Pmpid int `json:"program_map_PID"`
}

// parses raw section data and returns a ptr to a Pat
// if the section is malformed, error will be set
func NewPat(data []byte) (*Pat, error) {
	pat := &Pat{}
	if len(data) < 8 {
		return pat, errors.New("PAT data too short to parse!")
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
		pid := ((int(data[i+2]) & 31) << 8) + int(data[i+3])
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
	Tid         byte `json:"table_id"`
	Ssi         bool `json:"section_syntax_indicator"`
	Sl          int  `json:"section_length"`
	Pn          int  `json:"program_number"`
	Vn          int  `json:"version_number"`
	Cni         bool `json:"current_next_indicator"`
	Sn          byte `json:"section_number"`
	Lsn         byte `json:"last_section_number"`
	Pcrpid      int  `json:"PCR_PID"`
	Pil         int  `json:"program_info_length"`
	Descriptors []descriptor
	Pels        []pel
	Crc         uint32 `json:"CRC_32"`
}
type pel struct {
	St          byte `json:"stream_type"`
	Pid         int  `json:"elementary_PID"`
	Eil         int  `json:"ES_info_length"`
	Descriptors []descriptor
}

// parses raw section data and returns a ptr to a Pmt
// if the section is malformed, error will be set
func NewPmt(data []byte) (*Pmt, error) {
	pmt := &Pmt{}

	if len(data) < 12 {
		return pmt, errors.New("PMT data too short to parse!")
	}
	if data[0] != 0x2 {
		return pmt, errors.New(fmt.Sprintf("Invalid Table ID %x for PMT!", data[0]))
	}

	pmt.Tid = data[0]
	pmt.Ssi = data[1]&128 != 0
	pmt.Sl = (int(data[1]&15) << 8) + int(data[2])
	pmt.Pn = (int(data[3]) << 8) + int(data[4])
	pmt.Vn = int(data[5]&62) >> 1
	pmt.Cni = data[5]&1 != 0
	pmt.Sn = data[6]
	pmt.Lsn = data[7]
	pmt.Pcrpid = ((int(data[8]) & 31) << 8) + int(data[9])
	pmt.Pil = ((int(data[10]) & 15) << 8) + int(data[11])

	if len(data) < pmt.Sl+3 {
		return pmt, errors.New("PMT data length and section_length field mismatch!")
	}

	// TODO parse the descriptors, rather than skipping
	i := 12 + pmt.Pil
	pmt.Pels = make([]pel, 0)

	for i < pmt.Sl-4 {
		st := data[i]
		pid := ((int(data[i+1]) & 31) << 8) + int(data[i+2])
		eil := ((int(data[i+3]) & 15) << 8) + int(data[i+4])
		pmt.Pels = append(pmt.Pels, pel{st, pid, eil, nil})

		// TODO don't skip the descriptors
		i += 5 + eil
	}

	pmt.Crc = uint32(data[i+3]) + uint32(data[i+2])<<8 + uint32(data[i+1])<<16 + uint32(data[i])<<24

	return pmt, nil
}

// Represents a Service Description Table
// See ETSI 300468 section 5.2.3
type Sdt struct {
	Tid  byte `json:"table_id"`
	Ssi  bool `json:"section_syntax_indicator"`
	Sl   int  `json:"section_length"`
	Tsid int  `json:"transport_stream_id"`
	Vn   int  `json:"version_number"`
	Cni  bool `json:"current_next_indicator"`
	Sn   byte `json:"section_number"`
	Lsn  byte `json:"last_section_number"`
	Onid int  `json:"original_network_id"`
	Svcs []svc
	Crc  uint32 `json:"CRC_32"`
}

type svc struct {
	Sid         int  `json:"service_id"`
	Esf         bool `json:"EIT_schedule_flag"`
	Epf         bool `json:"EIT_present_following_flag"`
	Rs          int  `json:"running_status"`
	Fcm         bool `json:"free_CA_mode"`
	Dll         int  `json:"descriptors_loop_length"`
	Descriptors []descriptor
}

// parses raw section data and returns a ptr to an Sdt
// if the section is malformed, error will be set
func NewSdt(data []byte) (*Sdt, error) {
	sdt := &Sdt{}
	if len(data) < 8 {
		return sdt, errors.New("SDT data too short to parse!")
	}
	sdt.Tid = data[0]
	sdt.Ssi = data[1]&128 != 0
	sdt.Sl = (int(data[1]&15) << 8) + int(data[2])
	sdt.Tsid = (int(data[3]) << 8) + int(data[4])
	sdt.Vn = int(data[5]&62) >> 1
	sdt.Cni = data[5]&1 != 0
	sdt.Sn = data[6]
	sdt.Lsn = data[7]
	sdt.Onid = (int(data[8]) << 8) + int(data[9])

	if len(data) < sdt.Sl+3 {
		return sdt, errors.New("SDT data length and section_length field mismatch!")
	}

	sdt.Svcs = make([]svc, 0)
	i := 11
	for i < sdt.Sl-4 {
		svc := svc{}
		svc.Sid = (int(data[i]) << 8) + int(data[i+1])
		svc.Esf = data[i+2]&2 != 0
		svc.Epf = data[i+2]&1 != 0
		svc.Rs = int(data[i+3] >> 5)
		svc.Fcm = data[i+3]&16 != 0
		svc.Dll = (int(data[i+3]&15) << 8) + int(data[i+4])

		if len(data) < i+5+svc.Dll {
			return sdt, errors.New("SDT data length and descriptor_loop_length mismatch!")
		}

		svc.Descriptors = ParseDescriptors(data[i+5 : i+5+svc.Dll])

		sdt.Svcs = append(sdt.Svcs, svc)
		i += 5 + svc.Dll
	}

	sdt.Crc = uint32(data[i+3]) + uint32(data[i+2])<<8 + uint32(data[i+1])<<16 + uint32(data[i])<<24

	if crc := CalculateCrc32(data[:i]); crc != sdt.Crc {
		return sdt, errors.New(fmt.Sprintf("CRC error! Calculated %x but SDT says %x", crc, sdt.Crc))
	}
	return sdt, nil
}
