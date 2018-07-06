package crc

import (
	"encoding/binary"
)


func MakeCrc32Table(table []uint32, poly uint32) {
	for i := 0; i < 256; i += 1 {
		table[i] = Crc32UpdateBits(0, (uint32)(i), poly)
	}
}

func Crc32UpdateBits(crc uint32, data uint32, poly uint32) uint32 {
	for i := 0; i < 32; i += 1 {
		if ((data ^ crc) & 0x80000000) != 0 {
			crc <<= 1
			crc ^= poly
		} else {
			crc <<= 1
		}
		data <<= 1
	}
	return crc
}

func Crc32Update(crc uint32, data uint32, crcTable []uint32) uint32 {
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, data)
	crc = crcTable[a[3]^(byte(crc>>24)&0xff)] ^ (crc << 8)
	crc = crcTable[a[2]^(byte(crc>>24)&0xff)] ^ (crc << 8)
	crc = crcTable[a[1]^(byte(crc>>24)&0xff)] ^ (crc << 8)
	crc = crcTable[a[0]^(byte(crc>>24)&0xff)] ^ (crc << 8)
	return crc
}

func Crc32UpdateBlock(crc uint32, data []byte, crcTable []uint32) uint32 {
	for len(data)%4 != 0 {
		data = append(data, 0)
	}
	for i := 0; i < len(data); i += 4 {
		w := binary.LittleEndian.Uint32(data[i : i+4])
		crc = Crc32Update(crc, w, crcTable)
	}
	return crc
}
