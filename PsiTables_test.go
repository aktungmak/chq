package main

import (
	"testing"
)

// Sample Data

var sampPat = []byte{
	0x00, 0xB0, 0x11, 0x00, 0x00, 0xC5, 0x00, 0x00, 0x00, 0x01,
	0xE0, 0x38, 0x00, 0x00, 0xE0, 0x10, 0xE0, 0xAA, 0x66, 0x62}

var sampPmt = []byte{
	0x02, 0xB0, 0x58, 0x1C, 0x22, 0xC3, 0x00, 0x00, 0xE0, 0x70,
	0xF0, 0x0B, 0x0E, 0x03, 0xC0, 0x22, 0x33, 0x0C, 0x04, 0x80,
	0xB4, 0x81, 0x68, 0x02, 0xE0, 0x70, 0xF0, 0x16, 0x52, 0x01,
	0x01, 0x02, 0x03, 0x1A, 0x48, 0x5F, 0x06, 0x01, 0x02, 0x0E,
	0x03, 0xC0, 0x1E, 0xF6, 0x28, 0x04, 0x4D, 0x40, 0x1E, 0x3F,
	0x04, 0xE0, 0x71, 0xF0, 0x11, 0x52, 0x01, 0x02, 0x03, 0x01,
	0x67, 0x0E, 0x03, 0xC0, 0x00, 0xB4, 0x0A, 0x04, 0x61, 0x66,
	0x72, 0x00, 0x06, 0xE0, 0x72, 0xF0, 0x0A, 0x52, 0x01, 0x03,
	0x0E, 0x03, 0xC0, 0x02, 0x88, 0x56, 0x00, 0xF7, 0xB8, 0xEF,
	0x6C}

// Tests

func TestNewPat(t *testing.T) {
	pat, err := NewPat(sampPat)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if pat.Tid != 0x0 {
		t.Errorf("Table ID should be 0x0, got %x", pat.Tid)
	}
	if pat.Ssi != true {
		t.Errorf("Section Syntax Indicator should be true, got %t", pat.Ssi)
	}
	if pat.Sl != 0x11 {
		t.Errorf("Section Length should be 0x11, got %x", pat.Sl)
	}
	if pat.Tsid != 0x0 {
		t.Errorf("TSID should be 0x0, got %x", pat.Tsid)
	}
	if pat.Vn != 0x2 {
		t.Errorf("Version should be 0x2, got %x", pat.Vn)
	}
	if pat.Cni != true {
		t.Errorf("Current/next indicator should be true, got %t", pat.Cni)
	}
	if pat.Sn != 0x0 {
		t.Errorf("Section number should be 0x0, got %x", pat.Sn)
	}
	if pat.Lsn != 0x0 {
		t.Errorf("Last section number should be 0x0, got %x", pat.Lsn)
	}

	// check the number of programs
	if len(pat.Pgms) != 2 {
		t.Fatalf("Should be 2 programs in test PAT, got %d", len(pat.Pgms))
	}
	// programme 1
	if pat.Pgms[0].Pn != 1 {
		t.Errorf("PGM1 Programme number should be 1, got %d", pat.Pgms[0].Pn)
	}
	if pat.Pgms[0].Pmpid != 56 {
		t.Errorf("PGM1 PMT PID should be 56, got %d", pat.Pgms[0].Pmpid)
	}
	// programme 0
	if pat.Pgms[1].Pn != 0 {
		t.Errorf("PGM0 Programme number should be 0, got %d", pat.Pgms[1].Pn)
	}
	if pat.Pgms[1].Pmpid != 16 {
		t.Errorf("PGM0 PMT PID should be 16, got %d", pat.Pgms[1].Pmpid)
	}

	// ensure CRC was copied correctly
	if pat.Crc != 0xe0aa6662 {
		t.Errorf("PMT CRC was not copied correctly, expected 0xe0aa6662 got %x", pat.Crc)
	}
}

func TestNewPmt(t *testing.T) {
	pmt, err := NewPmt(sampPmt)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if pmt.Tid != 0x2 {
		t.Errorf("Table ID should be 0x2, got %x", pmt.Tid)
	}
	if pmt.Ssi != true {
		t.Errorf("Section Syntax Indicator should be true, got %t", pmt.Ssi)
	}
	if pmt.Sl != 0x58 {
		t.Errorf("Section Length should be 0x58, got %x", pmt.Sl)
	}
	if pmt.Pn != 0x1c22 {
		t.Errorf("Programme Number should be 0x1c22, got %x", pmt.Pn)
	}
	if pmt.Vn != 0x1 {
		t.Errorf("Version should be 0x1, got %x", pmt.Vn)
	}
	if pmt.Cni != true {
		t.Errorf("Current/next indicator should be true, got %t", pmt.Cni)
	}
	if pmt.Sn != 0x0 {
		t.Errorf("Section number should be 0x0, got %x", pmt.Sn)
	}
	if pmt.Lsn != 0x0 {
		t.Errorf("Last section number should be 0x0, got %x", pmt.Lsn)
	}
	if pmt.Pcrpid != 0x70 {
		t.Errorf("PCR PID should be 0x70, got %x", pmt.Pcrpid)
	}
	if pmt.Pil != 0xb {
		t.Errorf("Programme info length should be 0xb, got %x", pmt.Pil)
	}

	// check the number of programme elements
	if len(pmt.Pels) != 3 {
		t.Fatalf("Should be 3 pels in test PMT, got %d", len(pmt.Pels))
	}
	// Programme element 112
	if pmt.Pels[0].St != 0x2 {
		t.Errorf("PID 112 Stream type should be 0x2, got %x", pmt.Pels[0].St)
	}
	if pmt.Pels[0].Pid != 0x70 {
		t.Errorf("PID 112 PID should be 0x70, got %x", pmt.Pels[0].Pid)
	}
	if pmt.Pels[0].Eil != 0x16 {
		t.Errorf("PID 112 ES info len should be 0x16, got %x", pmt.Pels[0].Eil)
	}
	// Programme element 113
	if pmt.Pels[1].St != 0x4 {
		t.Errorf("PID 113 Stream type should be 0x4, got %x", pmt.Pels[1].St)
	}
	if pmt.Pels[1].Pid != 0x71 {
		t.Errorf("PID 113 PID should be 0x71, got %x", pmt.Pels[1].Pid)
	}
	if pmt.Pels[1].Eil != 0x11 {
		t.Errorf("PID 113 ES info len should be 0x11, got %x", pmt.Pels[1].Eil)
	}
	// Programme element 114
	if pmt.Pels[2].St != 0x6 {
		t.Errorf("PID 114 Stream type should be 0x6, got %x", pmt.Pels[2].St)
	}
	if pmt.Pels[2].Pid != 0x72 {
		t.Errorf("PID 114 PID should be 0x72, got %x", pmt.Pels[2].Pid)
	}
	if pmt.Pels[2].Eil != 0xa {
		t.Errorf("PID 114 ES info len should be 0xa, got %x", pmt.Pels[2].Eil)
	}

	// ensure CRC was copied correctly
	if pmt.Crc != 0xf7b8ef6c {
		t.Errorf("PMT CRC was not copied correctly, expected 0xf7b8ef6c got %x", pmt.Crc)
	}

}

// Benchmarks

func BenchmarkNewPat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPat(sampPat)
	}
}

func BenchmarkNewPmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPmt(sampPmt)
	}
}
