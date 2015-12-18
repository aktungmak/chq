package main

// TSPacket represents a single MPEG Transport Stream packet
// It holds pointers to each part of the packet, to reduce copying
// The Comment field is used to mark errors in the packet etc
type TsPacket struct {
	Header          *TsPacketHeader
	AdaptationField *AdaptationField `json:"adaptation_field"`
	Payload         []byte           `json:"data_byte"`
	Comment         string           // can be used for logging, don't use newlines
}

func NewTsPacket(data []byte) TsPacket {
	hdr := NewTsPacketHeader(data)
	af := NewAdaptationField(data)
	// copy the data, don't reuse the backing array
	// otherwise you will get a race!
	payld := append([]byte(nil), data[4+af.Length:]...)

	return TsPacket{
		Header:          hdr,
		AdaptationField: af,
		Payload:         payld,
	}
}

// Represents the 4-byte transport stream packet header
// See ISO 13818-1 Table 2-2
type TsPacketHeader struct {
	SyncByte byte `json:"sync_byte"`
	Tei      bool `json:"transport_error_indicator"`
	Pusi     bool `json:"payload_unit_start_indicator"`
	Tp       bool `json:"transport_priority"`
	Pid      int  `json:"PID"`
	Tsc      byte `json:"transport_scrambling_control"`
	Afc      byte `json:"adaptation_field_control"`
	Cc       byte `json:"continuity_counter"`
}

// Constructor to create a new TS header struct
// expects data to start with a sync byte
func NewTsPacketHeader(data []byte) *TsPacketHeader {
	return &TsPacketHeader{
		SyncByte: data[0],
		Tei:      data[1]&128 != 0,
		Pusi:     data[1]&64 != 0,
		Tp:       data[1]&32 != 0,
		Pid:      ((int(data[1]) & 31) << 8) + int(data[2]),
		Tsc:      (data[3] & 192) >> 6,
		Afc:      (data[3] & 48) >> 4,
		Cc:       data[3] & 15,
	}
}

// represents the optional adaptation field
// see ISO 13818-1 Table 2-6
type AdaptationField struct {
	Length byte `json:"adaptation_field_length"`
	Di     bool `json:"discontinuity_indicator"`
	Rai    bool `json:"random_access_indicator"`
	Espi   bool `json:"elementary_stream_priority_indicator"`
	Pcrf   bool `json:"PCR_flag"`
	Opcrf  bool `json:"OPCR_flag"`
	Spf    bool `json:"splicing_point_flag"`
	Tpdf   bool `json:"transport_private_data_flag"`
	Afef   bool `json:"adaptation_field_extension_flag"`
	//if PCR flag == 1
	Pcrb int64 `json:"program_clock_reference_base"`
	Pcre int64 `json:"program_clock_reference_extension"`
	//if OPCR flag == 1
	Opcrb int64 `json:"original_program_clock_reference_base"`
	Opcre int64 `json:"original_program_clock_reference_extension"`
	//if splicing_point_flag == 1
	Sc byte `json:"splice_countdown"`
	//if transport_private_data_flag == 1
	Tpdl byte   `json:"transport_private_data_length"`
	Tpd  []byte `json:"private_data_byte"`
	//if adaptation_field_extension_flag == 1
	Afel byte `json:"adaptation_field_extension_length"`
	Ltwf bool `json:"ltw_flag"`
	Pwrf bool `json:"piecewise_rate_flag"`
	Ssf  bool `json:"seamless_splice_flag"`
	//if `json:"ltw_flag` == "1
	Ltwvf bool `json:"ltw_valid_flag"`
	Ltwo  int  `json:"ltw_offset"`
	//if `json:"piecewise_rate_flag` == "1
	Pwr int `json:"piecewise_rate"`
	//if `json:"seamless_splice_flag` == "1
	St  int   `json:"splice_type"`
	Dna int64 `json:"DTS_next_AU"`
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
