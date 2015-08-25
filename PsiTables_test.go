package main

import (
	"testing"
)

// sample Program Association Table
var sampPat = []byte{
	0x00, 0xB0, 0x11, 0x00, 0x00, 0xC5, 0x00, 0x00, 0x00, 0x01,
	0xE0, 0x38, 0x00, 0x00, 0xE0, 0x10, 0xE0, 0xAA, 0x66, 0x62}

func TestNewPat(t *testing.T) {
	pat, err := NewPat(sampPat)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if pat.Tid != 0x0 {
		t.Errorf("Table ID should be 0x0, got %x", pat.Tid)
	}
	if pat.Ssi != true {
		t.Errorf("Section Syntax Indicator should be true, got %s", pat.Ssi)
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
		t.Errorf("Current/next indicator should be true, got %s", pat.Cni)
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

func BenchmarkNewPat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPat(sampPat)
	}
}
