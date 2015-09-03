package main

import (
	"testing"
)

// Sample Data

var sampSis = []byte{
	0xFC, 0x30, 0x20, 0x00, 0x00, 0xC3, 0x7D, 0xFF, 0x13, 0x00,
	0xFF, 0xFF, 0xFF, 0x05, 0x44, 0x00, 0x09, 0x2A, 0x7F, 0xCF,
	0xFE, 0xFB, 0x81, 0xDA, 0x16, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x00, 0xD0, 0xF1, 0xED, 0xB5}

// Tests

func TestNewScte35SpliceInfo(t *testing.T) {
	sis, err := NewScte35SpliceInfo(sampSis)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if sis.Tid != 0xfc {
		t.Errorf("Table ID should be 0xfc, got %x", sis.Tid)
	}
	if sis.Ssi != false {
		t.Errorf("Section Syntax Indicator should be false, got %t", sis.Ssi)
	}
	if sis.Pi != false {
		t.Errorf("Private Indicator should be false, got %t", sis.Sl)
	}
	if sis.Sl != 0x25 {
		t.Errorf("Section Length should be 0x25, got %x", sis.Sl)
	}
	if sis.Pv != 0x0 {
		t.Errorf("Protocol Version should be 0x0, got %x", sis.Pv)
	}
	if sis.Ep != false {
		t.Errorf("Encrypted flag should be false, got %t", sis.Ep)
	}
	if sis.Ea != 0x0 {
		t.Errorf("Encryption algorithm should be 0x0, got %x", sis.Ea)
	}
	if sis.Ptsa != 0xc37dff13 {
		t.Errorf("PTS adjustment should be 0xc37dff13, got %x", sis.Ptsa)
	}
	if sis.Cwi != 0x0 {
		t.Errorf("CW index should be 0x0, got %x", sis.Cwi)
	}
	if sis.Tier != 0x0 {
		t.Errorf("Tier should be 0x0, got %x", sis.Tier)
	}
	if sis.Scl != 0xfff {
		t.Errorf("Splice Command len should be 0xfff, got %x", sis.Scl)
	}
	if sis.Sct != 0x25 {
		t.Errorf("Splice Command type should be 0x25, got %x", sis.Sct)
	}

	// ensure CRC was copied correctly
	if sis.Crc != 0xd0f1edb5 {
		t.Errorf("CRC was not copied correctly, expected 0xd0f1edb5 got %x", sis.Crc)
	}

}
