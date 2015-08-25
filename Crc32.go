package main

var crcTable [256]uint32

func init() {
	// only generate the table once, at startup
	crcTable = generateTable()
}

// generate a CRC lookup table appropriate for MPEG2
// see ISO 13818-1 Annex A
func generateTable() [256]uint32 {
	var table [256]uint32
	var i, j, k uint32

	for i = 0; i < 256; i++ {
		k = 0
		for j = ((i << 24) | 0x800000); j != 0x80000000; j <<= 1 {
			if ((k ^ j) & 0x80000000) != 0 {
				k = (k << 1) ^ 0x04c11db7
			} else {
				k = (k << 1) ^ 0
			}
		}
		table[i] = k
	}
	return table
}

// calculate a CRC for the supplied byte slice
func CalculateCrc32(data []byte) uint32 {
	var crc uint32 = 0xFFFFFFFF
	for _, b := range data {
		crc = (crc << 8) ^ crcTable[((crc>>24)^uint32(b))&0xff]
	}
	return crc
}
