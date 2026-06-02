package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/saferwall/pe"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing binary path")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	f, err := pe.NewBytes(data, &pe.Options{})
	if err != nil {
		log.Fatalf("Failed to parse PE: %v", err)
	}
	if err := f.Parse(); err != nil {
		log.Fatalf("Failed to parse PE: %v", err)
	}

	buildID, err := scanBuildID(f, data)
	if err != nil {
		log.Fatalf("Failed to scan build ID: %v", err)
	}
	log.Println("PsyBuildID:", buildID)

	featureSet, err := scanFeatureSet(f, data)
	if err != nil {
		log.Fatalf("Failed to scan feature set: %v", err)
	}
	log.Println("FeatureSet:", featureSet)
}

func scanBuildID(f *pe.File, data []byte) (string, error) {
	var exportRVA uint32
	for _, exp := range f.Export.Functions {
		if exp.Name == "GPsyonixBuildID" {
			exportRVA = exp.FunctionRVA
			break
		}
	}
	if exportRVA == 0 {
		return "", fmt.Errorf("GPsyonixBuildID export not found")
	}

	var base uint64
	if f.Is64 {
		base = f.NtHeader.OptionalHeader.(pe.ImageOptionalHeader64).ImageBase
	} else {
		base = uint64(f.NtHeader.OptionalHeader.(pe.ImageOptionalHeader32).ImageBase)
	}

	// export points to a pointer-sized variable; dereference it to get the string VA
	ptrOff := int(f.GetOffsetFromRva(exportRVA))
	if ptrOff+8 > len(data) {
		return "", fmt.Errorf("build ID pointer out of bounds")
	}
	stringVA := binary.LittleEndian.Uint64(data[ptrOff:])
	stringRVA := uint32(stringVA - base)

	strOff := int(f.GetOffsetFromRva(stringRVA))
	if strOff >= len(data) {
		return "", fmt.Errorf("build ID string out of bounds")
	}
	return decodeUTF16(data[strOff:]), nil
}

// "PrimeUpdate" in UTF-16LE
var primeUpdateNeedle = []byte("P\x00r\x00i\x00m\x00e\x00U\x00p\x00d\x00a\x00t\x00e\x00")

func scanFeatureSet(f *pe.File, data []byte) (string, error) {
	var rdataSect *pe.Section
	for i := range f.Sections {
		if strings.HasPrefix(f.Sections[i].String(), ".rdata") {
			rdataSect = &f.Sections[i]
			break
		}
	}
	if rdataSect == nil {
		return "", fmt.Errorf("no .rdata section")
	}
	start := int(rdataSect.Header.PointerToRawData)
	rdata := data[start : start+int(rdataSect.Header.SizeOfRawData)]

	type update struct {
		major  uint64
		suffix string
		full   string
	}
	seen := make(map[[2]string]bool)
	var updates []update

	for offset := 0; offset+len(primeUpdateNeedle) <= len(rdata); {
		pos := bytes.Index(rdata[offset:], primeUpdateNeedle)
		if pos < 0 {
			break
		}
		absPos := offset + pos
		offset = absPos + 1
		if absPos%2 != 0 {
			continue
		}

		rest := rdata[absPos+len(primeUpdateNeedle):]
		var chars []uint16
		for i := 0; i+1 < len(rest); i += 2 {
			c := binary.LittleEndian.Uint16(rest[i:])
			if c <= 0x20 {
				break
			}
			chars = append(chars, c)
		}
		if len(chars) == 0 {
			continue
		}
		tail := string(utf16.Decode(chars))

		digitEnd := 0
		for digitEnd < len(tail) && tail[digitEnd] >= '0' && tail[digitEnd] <= '9' {
			digitEnd++
		}
		if digitEnd == 0 {
			continue
		}
		var major uint64
		for _, ch := range tail[:digitEnd] {
			major = major*10 + uint64(ch-'0')
		}
		suffix := tail[digitEnd:]

		key := [2]string{tail[:digitEnd], suffix}
		if !seen[key] {
			seen[key] = true
			updates = append(updates, update{major, suffix, "PrimeUpdate" + tail})
		}
	}
	if len(updates) == 0 {
		return "", fmt.Errorf("PrimeUpdate string not found")
	}

	sort.Slice(updates, func(i, j int) bool {
		if updates[i].major != updates[j].major {
			return updates[i].major > updates[j].major
		}
		return updates[i].suffix > updates[j].suffix
	})
	return updates[0].full, nil
}

func decodeUTF16(data []byte) string {
	var chars []uint16
	for i := 0; i+1 < len(data); i += 2 {
		c := binary.LittleEndian.Uint16(data[i:])
		if c == 0 {
			break
		}
		chars = append(chars, c)
	}
	return string(utf16.Decode(chars))
}
