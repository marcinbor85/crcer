package crc

import (
	"encoding/binary"
	"github.com/marcinbor85/gohex"
	"math/rand"
	"time"
)

func GetAdressRange(mem *gohex.Memory) (uint32, uint32) {
	firstSegment := mem.GetDataSegments()[0]
	lastSegment := mem.GetDataSegments()[len(mem.GetDataSegments())-1]
	return firstSegment.Address, lastSegment.Address + (uint32)(len(lastSegment.Data))
}

func AddDoubleCrc32(mem *gohex.Memory, startAdr uint32, endAdr uint32, pad byte) error {
	rand.Seed(time.Now().UTC().UnixNano())

	data := mem.ToBinary(startAdr, endAdr - startAdr, pad)
	
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, Crc32UpdateBlock(0, data))
	err := mem.AddBinary(endAdr, a)
	if err != nil {
		return err
	}

	a = make([]byte, 4)
	binary.LittleEndian.PutUint32(a, (uint32)(rand.Intn(0xFFFFFFFF)))
	err = mem.AddBinary(endAdr+4, a)
	if err != nil {
		return err
	}

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, Crc32UpdateBlock(0, a))
	err = mem.AddBinary(endAdr+8, b)
	if err != nil {
		return err
	}
	
	return nil
}
