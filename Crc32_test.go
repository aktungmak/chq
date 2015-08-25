package main

import (
	"testing"
)

func TestgenerateTable(t *testing.T) {
	// data := []byte{0x00, 0xB0, 0x0D, 0x00, 0x01, 0xC5, 0x00, 0x00, 0x00, 0x01, 0xE1, 0xE0}
	tbl := generateTable()

	// check a few of the values are correct
	if tbl[0] != 0x0 ||
		tbl[1] != 0x4c11db7 ||
		tbl[254] != 0xb5365d03 ||
		tbl[255] != 0xb1f740b4 {
		t.Error("Values in the CRC table are not correct!")
	}
}

func TestCalculateCrc32(t *testing.T) {
	data := []byte{0x00, 0xb0, 0x0d, 0x00, 0x01, 0xc5, 0x00, 0x00, 0x00, 0x01, 0xe1, 0xe0}
	expect := uint32(0x14ccc5f7)
	crc := CalculateCrc32(data)
	if crc != expect {
		t.Errorf("CRC calculation incorrect, expected %x got %x", expect, crc)
	}
}

func BenchmarkGenerateTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateTable()
	}
}
