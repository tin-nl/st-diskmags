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
	flag.StringVar(&strFlag, "input", "DATA.LNK", "Disk Maggie DATA.LNK File")
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
	fileEntries := ReadHeader(f)
	ReadFiles(f, strFolder, fileEntries)
}

type FileEntry struct {
	offset uint32
	length uint32
	name   string
}

func ReadFiles(input *os.File, destFolder string, fileEntries []FileEntry) {
	for _, value := range fileEntries {
		fw, err := os.OpenFile(destFolder+"/"+value.name, os.O_CREATE|os.O_WRONLY, 0764)
		if err != nil {
			panic(err)
		}
		io.CopyN(fw, input, int64(value.length))
		fw.Close()
	}
}

func ReadHeader(input *os.File) []FileEntry {
	var entryCount uint16
	binary.Read(input, binary.BigEndian, &entryCount)
	fmt.Println(entryCount)

	var filename [12]byte
	fileEntries := make([]FileEntry, entryCount, entryCount)
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
	return fileEntries
}
