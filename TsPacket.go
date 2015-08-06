package main

// TSPacket represents a single MPEG Transport Stream packet
// It holds pointers to each part of the packet, to reduce copying
// The Comment field is used to mark errors in the packet etc
type TsPacket struct {
	Header          *TsPacketHeader
	AdaptationField *AdaptationField `mpeg:"adaptation_field"`
	Payload         []byte           `mpeg:"data_byte"`
	Comment         string           // can be used for logging, don't use newlines
}

func NewTsPacket(data []byte) TsPacket {
	hdr := NewTsPacketHeader(data)
	af := NewAdaptationField(data)
	payld := data[4+af.Length:]

	return TsPacket{
		Header:          hdr,
		AdaptationField: af,
		Payload:         payld,
	}
}

// Represents the 4-byte transport stream packet header
// See ISO 13818-1 Table 2-2
type TsPacketHeader struct {
	SyncByte byte  `mpeg:"sync_byte"`
	Tei      bool  `mpeg:"transport_error_indicator"`
	Pusi     bool  `mpeg:"payload_unit_start_indicator"`
	Tp       bool  `mpeg:"transport_priority"`
	Pid      int16 `mpeg:"PID"`
	Tsc      byte  `mpeg:"transport_scrambling_control"`
	Afc      byte  `mpeg:"adaptation_field_control"`
	Cc       byte  `mpeg:"continuity_counter"`
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

// represents the optional adaptation field
// see ISO 13818-1 Table 2-6
type AdaptationField struct {
	Length byte `mpeg:"adaptation_field_length"`
	Di     bool `mpeg:"discontinuity_indicator"`
	Rai    bool `mpeg:"random_access_indicator"`
	Espi   bool `mpeg:"elementary_stream_priority_indicator"`
	Pcrf   bool `mpeg:"PCR_flag"`
	Opcrf  bool `mpeg:"OPCR_flag"`
	Spf    bool `mpeg:"splicing_point_flag"`
	Tpdf   bool `mpeg:"transport_private_data_flag"`
	Afef   bool `mpeg:"adaptation_field_extension_flag"`
	//if PCR flag == 1
	Pcrb int64 `mpeg:"program_clock_reference_base"`
	Pcre int64 `mpeg:"program_clock_reference_extension"`
	//if OPCR flag == 1
	Opcrb int64 `mpeg:"original_program_clock_reference_base"`
	Opcre int64 `mpeg:"original_program_clock_reference_extension"`
	//if splicing_point_flag == 1
	Sc byte `mpeg:"splice_countdown"`
	//if transport_private_data_flag == 1
	Tpdl byte   `mpeg:"transport_private_data_length"`
	Tpd  []byte `mpeg:"private_data_byte"`
	//if adaptation_field_extension_flag == 1
	Afel byte `mpeg:"adaptation_field_extension_length"`
	Ltwf bool `mpeg:"ltw_flag"`
	Pwrf bool `mpeg:"piecewise_rate_flag"`
	Ssf  bool `mpeg:"seamless_splice_flag"`
	//if `mpeg:"ltw_flag` == "1
	Ltwvf bool  `mpeg:"ltw_valid_flag"`
	Ltwo  int16 `mpeg:"ltw_offset"`
	//if `mpeg:"piecewise_rate_flag` == "1
	Pwr int `mpeg:"piecewise_rate"`
	//if `mpeg:"seamless_splice_flag` == "1
	St  int   `mpeg:"splice_type"`
	Dna int64 `mpeg:"DTS_next_AU"`
}

// Constructor to create a new adaptation field struct
// expects data to start with a sync byte
// If AF is not present in the provided TS packet,
// returns zero struct if no af (Length = 0).
func NewAdaptationField(data []byte) *AdaptationField {
	af := AdaptationField{}
	if (data[3] & 32) == 0 {
		//no af
		return &af
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
		// TODO add the extension fields

	}
	return &af
}
