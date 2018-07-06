package main

import (
	"flag"
	"fmt"
	"github.com/marcinbor85/crcer/crc"
	"github.com/marcinbor85/gohex"
	"os"
)

func exit(msg string) {
	fmt.Fprintf(os.Stderr, "error: %v\n", msg)
	os.Exit(1)
}

func main() {
	method := flag.Uint("method", 0x00, "crc method (default 0 - double, 1 - single)")
	padding := flag.Uint("padding", 0xFF, "padding byte")
	poly := flag.Uint("poly", 0x04C11DB7, "polynomial")
	start := flag.Uint("start", 0x08000000, "start address")
	end := flag.Uint("end", 0x08040000, "end address")

	startSet := false
	endSet := false

	flag.Parse()
	
	crcTable := make([]uint32, 256)
	crc.MakeCrc32Table(crcTable, (uint32)(*poly))

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "end" {
			endSet = true
		}
		if f.Name == "start" {
			startSet = true
		}
	})

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

	startAdr, endAdr := crc.GetAdressRange(mem)

	if startSet != false {
		startAdr = (uint32)(*start)
	}

	if endSet != false {
		endAdr = (uint32)(*end)
	}
	
	if endAdr <= startAdr {
		exit("end address must be greater than start address")
	}
	
	if endAdr%4 != 0 || startAdr%4 != 0 {
		exit("addresses should be aligned to 4 bytes")
	}

	if *method == 0 {
		err := crc.AddDoubleCrc32(mem, startAdr, endAdr, (byte)(*padding), crcTable)
		if err != nil {
			exit("cannot add crc: " + err.Error())
		}

		mem.DumpIntelHex(os.Stdout, 16)
	} else if *method == 1 {
		err := crc.AddSingleCrc32(mem, startAdr, endAdr, (byte)(*padding), crcTable)
		if err != nil {
			exit("cannot add crc: " + err.Error())
		}

		mem.DumpIntelHex(os.Stdout, 16)
	} else {
		exit("unsupported crc method")
	}
}
