package main

// maybe the TsPacket should be changed to a generic packet
// which has a TYPE field and several fields, one for TS,
// one for PES etc.
type TsPacket struct {
	Header          *TsPacketHeader
	AdaptationField *AdaptationField `adaptation_field`
	Payload         []byte           `data_byte`
	Comment         string           // can be used for logging, don't use newlines
}

func NewTsPacket(data []byte) TsPacket {
	hdr := NewTsPacketHeader(data)
	af := NewAdaptationField(data)
	// payld := data[4+af.Length : 188]

	return TsPacket{
		Header:          hdr,
		AdaptationField: af,
		// Payload:         payld,
	}
}

// Represents the 4-byte transport stream packet header
// See ISO 13818-1 Table 2-2
type TsPacketHeader struct {
	SyncByte byte  `sync_byte`
	Tei      bool  `transport_error_indicator`
	Pusi     bool  `payload_unit_start_indicator`
	Tp       bool  `transport_priority`
	Pid      int16 `PID`
	Tsc      byte  `transport_scrambling_control`
	Afc      byte  `adaptation_field_control`
	Cc       byte  `continuity_counter`
}

// Constructor to create a new TS header struct
// expects data to start with a sync byte
func NewTsPacketHeader(data []byte) *TsPacketHeader {
	return &TsPacketHeader{
		SyncByte: data[0],
		Tei:      data[1]&128 != 0,
		Pusi:     data[1]&64 != 0,
		Tp:       data[1]&32 != 0,
		Pid:      ((int16(data[1]) & 31) << 8) + int16(data[2]),
		Tsc:      (data[3] & 192) >> 6,
		Afc:      (data[3] & 48) >> 4,
		Cc:       data[3] & 15,
	}
}

// represents
type AdaptationField struct {
	Length byte `adaptation_field_length`
	Di     bool `discontinuity_indicator`
	Rai    bool `random_access_indicator`
	Espi   bool `elementary_stream_priority_indicator`
	Pcrf   bool `PCR_flag`
	Opcrf  bool `OPCR_flag`
	Spf    bool `splicing_point_flag`
	Tpdf   bool `transport_private_data_flag`
	Afef   bool `adaptation_field_extension_flag`
	//if PCR flag == 1
	Pcrb int64 `program_clock_reference_base`
	Pcre int64 `program_clock_reference_extension`
	//if OPCR flag == 1
	Opcrb int64 `original_program_clock_reference_base`
	Opcre int64 `original_program_clock_reference_extension`
	//if splicing_point_flag == 1
	Sc byte `splice_countdown`
	//if transport_private_data_flag == 1
	Tpdl byte   `transport_private_data_length`
	Tpd  []byte `private_data_byte`
	//if adaptation_field_extension_flag == 1
	Afel byte `adaptation_field_extension_length`
	Ltwf bool `ltw_flag`
	Pwrf bool `piecewise_rate_flag`
	Ssf  bool `seamless_splice_flag`
	//if `ltw_flag` == 1
	Ltwvf bool  `ltw_valid_flag`
	Ltwo  int16 `ltw_offset`
	//if `piecewise_rate_flag` == 1
	Pwr int `piecewise_rate`
	//if `seamless_splice_flag` == 1
	St  int   `splice_type`
	Dna int64 `DTS_next_AU`
}

// Constructor to create a new adaptation field struct
// expects data to start with a sync byte
// If AF is not present in the provided TS packet,
// returns zero struct if no af (Length = 0).
func NewAdaptationField(data []byte) *AdaptationField {
	af := AdaptationField{}
	if (data[3] & 32) == 0 {
		//no af
		return nil
	}
	// log.Printf("%v", data)
	af.Length = data[4]
	af.Di = data[5]&128 != 0
	af.Rai = data[5]&64 != 0
	af.Espi = data[5]&32 != 0
	af.Pcrf = data[5]&16 != 0
	af.Opcrf = data[5]&8 != 0
	af.Spf = data[5]&4 != 0
	af.Tpdf = data[5]&2 != 0
	af.Afef = data[5]&1 != 0

	//keep track of byte offset depending on flags
	ofs := 6

	if af.Pcrf {
		af.Pcrb = 0
		af.Pcrb += int64(data[ofs]) << 25
		af.Pcrb += int64(data[ofs+1]) << 17
		af.Pcrb += int64(data[ofs+2]) << 9
		af.Pcrb += int64(data[ofs+3]) << 1
		af.Pcrb += int64(data[ofs+4] >> 7)

		af.Pcre = int64(data[ofs+4]&1) << 9
		af.Pcre += int64(data[ofs+5])
		ofs += 6
	}
	if af.Opcrf {
		af.Opcrb = 0
		af.Opcrb += int64(data[ofs]) << 25
		af.Opcrb += int64(data[ofs+1]) << 17
		af.Opcrb += int64(data[ofs+2]) << 9
		af.Opcrb += int64(data[ofs+3]) << 1
		af.Opcrb += int64(data[ofs+4] >> 7)

		af.Opcre = int64(data[ofs+4]&1) << 9
		af.Opcre += int64(data[ofs+5])
		ofs += 6
	}
	if af.Spf {
		af.Sc = data[ofs]
		ofs += 1
	}
	if af.Tpdf {
		af.Tpdl = data[ofs]
		af.Tpd = data[ofs+1 : ofs+int(af.Tpdl)+1]
		ofs += int(af.Tpdl + 1)
	}
	if af.Afef {
		af.Afel = data[ofs]
		ofs += 1

		af.Ltwf = data[ofs]&128 != 0
		af.Pwrf = data[ofs]&64 != 0
		af.Ssf = data[ofs]&32 != 0
		ofs += 1

	}
	return &af
}
