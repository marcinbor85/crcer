package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/marcinbor85/crcer/crc"
	"github.com/marcinbor85/gohex"
	"math/rand"
	"os"
	"time"
)

func exit(msg string) {
	fmt.Fprintf(os.Stderr, "error: %v\n", msg)
	os.Exit(1)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	method := flag.Uint("method", 0x00, "crc method")
	padding := flag.Uint("padding", 0xFF, "padding byte")

	flag.Parse()

	if len(flag.Args()) != 1 {
		exit("no input filename")
	}
	filename := flag.Args()[0]

	file, err := os.Open(filename)
	if err != nil {
		exit(err.Error())
	}
	defer file.Close()

	mem := gohex.NewMemory()
	err = mem.ParseIntelHex(file)
	if err != nil {
		exit(err.Error())
	}

	firstSegment := mem.GetDataSegments()[0]
	lastSegment := mem.GetDataSegments()[len(mem.GetDataSegments())-1]
	startAdr := firstSegment.Address
	endAdr := lastSegment.Address + (uint32)(len(lastSegment.Data))

	data := mem.ToBinary(startAdr, endAdr-startAdr, (byte)(*padding))

	if *method == 0 {
		a := make([]byte, 4)

		binary.LittleEndian.PutUint32(a, crc.Crc32UpdateBlock(0, data))
		mem.AddBinary(endAdr, a)

		binary.LittleEndian.PutUint32(a, (uint32)(rand.Intn(0xFFFFFFFF)))
		mem.AddBinary(endAdr+4, a)

		binary.LittleEndian.PutUint32(a, crc.Crc32UpdateBlock(0, a))
		mem.AddBinary(endAdr+8, a)

		mem.DumpIntelHex(os.Stdout, 16)
	} else {
		exit("unsupported crc method")
	}
}
