package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

var strFlag, strFolder string //= flag.String("long-string", "DATA.LNK", "Disk Maggie DATA.LNK File")

func init() {
	// example with short version for long flag
	flag.StringVar(&strFlag, "input", "DATA", "UCM 1-7 DATA File")
	flag.StringVar(&strFolder, "output-folder", "unpacked", "Output Folder (will be created)")
}

func main() {
	flag.Parse()
	fmt.Println(len(os.Args), os.Args)
	fmt.Println(strFlag)
	if _, err := os.Stat(strFolder); os.IsNotExist(err) {
		fmt.Println("Create folder \"" + strFolder + "\"")
		os.Mkdir(strFolder, 0764)
	}
	f, err := os.Open(strFlag)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var head byte
	binary.Read(f, binary.BigEndian, &head)
	f.Seek(0, 0)
	if head == 0x9d {
		ReadHeaderV1(f, strFolder)
	} else {
		fileEntries := ReadHeaderV2(f, strFolder)
		ReadFilesV2(f, strFolder, fileEntries)
	}
}

type FileEntry struct {
	offset     uint32
	packlength uint32
	length     uint32
	name       string
}

var eor []byte = []byte{'w', 'H', 'a', 'T', ' ', 'y', 'O', 'u', ' ', 'r', 'E', 'a', 'P', ' ', 'i', 'S', ' ', 'w', 'H', 'a', 'T', ' ', 'y', 'O', 'u', ' ', 's', 'O', 'w'}

func unxor(data []byte) {
	var D4 uint32 = 0x31415926
	eorIdx := 0
	for key, value := range data {
		eorbyte := byte(eor[eorIdx])
		//>fmt.Printf("%x ", eorbyte)
		eorIdx++
		if eorIdx >= len(eor) {
			eorIdx = 0
		}
		eorbyte += byte(D4 & 0xff)
		data[key] = value ^ eorbyte
		//ror.l #1
		D4 = (D4 >> 1) | ((D4 & 1) << 31)
		D4++
	}
	//>fmt.Println()
}

/*
PC=426f4
a0=acb1c
lea eor,A1
move.l #$31415926,D4
.wrap:
move.l A1,A2
.loop:
move.b (A2)+,D5
beq.s .wrap
add.b 	d4,d5
eor.b   D5,(A0)+
ror.l  #1,D4
addq   #1,D4
dbra  D0,.loop
sub.l 	#$10000,D0
bpl.s  .loop
*/
func ReadFilesV2(input *os.File, destFolder string, fileEntries []FileEntry) {
	for _, value := range fileEntries {
		fw, err := os.OpenFile(destFolder+"/"+value.name, os.O_CREATE|os.O_WRONLY, 0764)
		if err != nil {
			panic(err)
		}
		input.Seek(int64(value.offset), 0)
		io.CopyN(fw, input, int64(value.packlength))
		fw.Close()
	}
}

func ReadHeaderV1(input *os.File, destFolder string) {
	var maxlen = 0x1802
	var data = make([]byte, maxlen, maxlen)
	binary.Read(input, binary.BigEndian, &data)
	unxor(data)
	var entryCount uint16 = uint16(data[0])<<8 | uint16(data[1])
	fmt.Println(entryCount)
	fileEntries := make([]FileEntry, entryCount, entryCount)
	/*
		for _, value := range data {
			fmt.Printf("%x", byte(value))
			fmt.Printf("[%s]", string(byte(value)))
		}
	*/
	var filename []byte
	idx := 2
	for iCnt := 0; iCnt < int(entryCount); iCnt++ {
		filename = data[idx : idx+12]
		fmt.Println(string(filename))
		idx += 12
		startoffset := uint32(data[idx])<<24 | uint32(data[idx+1])<<16 | uint32(data[idx+2])<<8 | uint32(data[idx+3])
		fileEntries[iCnt].offset = startoffset
		idx += 4
		lenpack := uint32(data[idx])<<24 | uint32(data[idx+1])<<16 | uint32(data[idx+2])<<8 | uint32(data[idx+3])
		fileEntries[iCnt].packlength = lenpack
		idx += 4
		lenunpack := uint32(data[idx])<<24 | uint32(data[idx+1])<<16 | uint32(data[idx+2])<<8 | uint32(data[idx+3])
		fileEntries[iCnt].length = lenunpack
		idx += 4
		fmt.Printf("Offset: %d\n", (startoffset))
		fmt.Printf("Packlen: %d\n", (lenpack))
		fmt.Printf("Uplen: %d\n", (lenunpack))
		input.Seek(int64(startoffset), 0)
		nulpos := bytes.Index(filename[:12], []byte{0})
		if nulpos > 0 {
			fileEntries[iCnt].name = string(filename[:nulpos])
		} else {
			fileEntries[iCnt].name = string(filename[:12])

		}
		fw, err := os.OpenFile(destFolder+"/"+fileEntries[iCnt].name, os.O_CREATE|os.O_WRONLY, 0764)
		if err != nil {
			panic(err)
		}
		filedata := make([]byte, lenpack, lenpack)
		io.ReadAtLeast(input, filedata, int(lenpack))
		unxor(filedata)
		fw.Write(filedata)
		fw.Close()
	}
	/*
		fmt.Println(entryCount)

		var filename [12]byte
		for iCnt := 0; iCnt < int(entryCount); iCnt++ {
			binary.Read(input, binary.BigEndian, &fileEntries[iCnt].offset)
			binary.Read(input, binary.BigEndian, &filename)
			nulpos := bytes.Index(filename[:12], []byte{0})
			if nulpos > 0 {
				fileEntries[iCnt].name = string(filename[:nulpos])
			} else {
				fileEntries[iCnt].name = string(filename[:12])

			}
			if iCnt > 0 {
				fileEntries[iCnt-1].length = fileEntries[iCnt].offset - fileEntries[iCnt-1].offset
			}
			fmt.Println(fileEntries[iCnt].name + "...")
		}
		var lastoffset uint32
		binary.Read(input, binary.BigEndian, &lastoffset)
		fileEntries[entryCount-1].length = lastoffset - fileEntries[entryCount-1].offset
	*/
}
func ReadHeaderV2(input *os.File, destFolder string) []FileEntry {
	entryCount := 256
	fmt.Println(entryCount)
	fileEntries := make([]FileEntry, 0, entryCount)
	var filename [16]byte
	offset := 0x1400
	lastoffset := offset
	iCnt := 0
	for iCnt = 0; iCnt < int(256); iCnt++ {
		var foffs uint32 = 0
		binary.Read(input, binary.BigEndian, &filename)
		binary.Read(input, binary.BigEndian, &foffs)
		var fname string
		nulpos := bytes.Index(filename[:16], []byte{0})
		if nulpos > 0 {
			fname = string(filename[:nulpos])
		} else {
			fname = string(filename[:16])
		}
		if foffs == 0 && iCnt > 0 {
			break
		}
		fileEntries = append(fileEntries, FileEntry{name: fname, offset: uint32(lastoffset), packlength: uint32((offset + int(foffs)) - lastoffset)})
		lastoffset = (offset + int(foffs))
		fmt.Println(fileEntries[len(fileEntries)-1].name)
		fmt.Println(fileEntries[len(fileEntries)-1].packlength)
	}
	return fileEntries
}
