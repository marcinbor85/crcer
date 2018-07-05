package crc

import (
	"testing"
	"github.com/marcinbor85/gohex"
	"reflect"
	"encoding/binary"
)

func prepareMem() *gohex.Memory {
	m := gohex.NewMemory()
	
	m.AddBinary(0x08005000, []byte{0x00, 0xA0, 0x00, 0x20})
	m.AddBinary(0x08005004, []byte{0x10, 0xA0, 0x00, 0x08, 0x12, 0xA0, 0x00, 0x08, 0x14, 0xA0, 0x00, 0x08, 0x16, 0xA0, 0x00, 0x08})
	m.AddBinary(0x08005014, []byte{0x00, 0x00, 0x00, 0x00})
	m.AddBinary(0x08005024, []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88})
	
	return m
}

func TestCrc32(t *testing.T) {
	var v uint32
	
	v = Crc32Update(0, 0x12345678)
	if v != 0x188E5750 {
		t.Errorf("wrong crc checksum: %08X", v)
	}
	
	v = Crc32Update(0xFFFFFFFF, 0x12345678)
	if v != 0xDF8A8A2B  {
		t.Errorf("wrong crc checksum: %08X", v)
	}
	
	v = Crc32UpdateBlock(0, []byte{0x78, 0x56, 0x34, 0x12, 0xF0, 0xDE, 0xBC, 0x9A})
	if v != 0x14201842  {
		t.Errorf("wrong crc checksum: %08X", v)
	}
	
	v = Crc32UpdateBlock(0xFFFFFFFF, []byte{0x78, 0x56, 0x34, 0x12, 0xF0, 0xDE, 0xBC, 0x9A})
	if v != 0x7D24A31B  {
		t.Errorf("wrong crc checksum: %08X", v)
	}
}

func TestGetAddressRange(t *testing.T) {
	m := prepareMem()
	
	start, end := GetAdressRange(m)
	if start != 0x08005000 {
		t.Errorf("wrong start address: %08X", start)
	}
	if end != 0x0800502C {
		t.Errorf("wrong end address: %08X", end)
	}
}

func TestAddDoubleCrc32(t *testing.T) {
	
	m := prepareMem()
	
	err := AddDoubleCrc32(m, 0x08005000, 0x0800502C, 0xFF)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	
	data := m.ToBinary(0x0800502C, 12, 0xFF)
	crc1 := data[0:4]
	rand := data[4:8]
	crc2 := data[8:12]
	
	if reflect.DeepEqual(crc1, []byte{132, 37, 27, 32}) == false {
		t.Errorf("wrong data crc sum: %v", crc1)
	}
	
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, Crc32UpdateBlock(0, rand))
	if reflect.DeepEqual(a, crc2) == false {
		t.Errorf("wrong rand crc sum: %v", crc2)
	}
	
	n := len(m.GetDataSegments()[1].Data)
	if n != 20 {
		t.Errorf("wrong datasegment length: %v", n)
	}
	
	m = prepareMem()
	
	err = AddDoubleCrc32(m, 0x08005000, 0x08005018, 0xFF)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	
	data = m.ToBinary(0x08005018, 12, 0xFF)
	crc1 = data[0:4]
	rand = data[4:8]
	crc2 = data[8:12]
	
	if reflect.DeepEqual(crc1, []byte{0xC8, 0xA1, 0x26, 0x87}) == false {
		t.Errorf("wrong data crc sum: %v", crc1)
	}
	
	binary.LittleEndian.PutUint32(a, Crc32UpdateBlock(0, rand))
	if reflect.DeepEqual(a, crc2) == false {
		t.Errorf("wrong rand crc sum: %v", crc2)
	}
	
	n = len(m.GetDataSegments())
	if n != 1 {
		t.Errorf("wrong datasegments num: %v", n)
	}
	
	n = len(m.GetDataSegments()[0].Data)
	if n != 44 {
		t.Errorf("wrong datasegment length: %v", n)
	}
	
	m = prepareMem()
	
	err = AddDoubleCrc32(m, 0x08005000, 0x08005010, 0xFF)
	if err == nil {
		t.Errorf("segment should be overlapped")
	}
	err = AddDoubleCrc32(m, 0x08005000, 0x08005004, 0xFF)
	if err == nil {
		t.Errorf("segment should be overlapped")
	}
}


func TestAddSingleCrc32(t *testing.T) {
	
	m := prepareMem()
	
	err := AddSingleCrc32(m, 0x08005000, 0x0800502C, 0xFF)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	
	data := m.ToBinary(0x0800502C, 4, 0xFF)
	crc1 := data[0:4]
	
	if reflect.DeepEqual(crc1, []byte{132, 37, 27, 32}) == false {
		t.Errorf("wrong data crc sum: %v", crc1)
	}
	
	n := len(m.GetDataSegments()[1].Data)
	if n != 12 {
		t.Errorf("wrong datasegment length: %v", n)
	}
}

