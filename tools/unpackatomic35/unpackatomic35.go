/*
needs GCC
Windows: http://win-builds.org/stable/#_getting_and_running_the_system
*/
package main

/*
//** Atomik 3.5 depacker, universal C version
//** function is reentrand and thread safe
//** data read and writes are checked against their bounds:
//** src to src+atm_35_packedsize() for the source data (read only)
//** dst to dst+atm_35_origsize() for the destination data (read and write)
//**
//** placed in public domain 2007 by Hans Wessels
//**
//** The function:
//** unsigned long atm_35_depack(const unsigned char *src, unsigned char *dst);
//** depacks Atomik 3.5 packed data located at src to dst
//** it returns the size of the depacked data excluding the
//** picture offsets at the end of the file.
//** Or -1 if no valid Atomik 3.5 header was found or in
//** case of an depack error. The memory at dst should be large
//** enough to hold the depacked data. This size can be determined
//** by using the function atm_35_origsize().
//**
//** The function:
//** int atm_35_header(const unsigned char *src);
//** returns 0 if no valid Atomik 3.5 header was found at src
//**
//** The function:
//** unsigned long atm_35_packedsize(const unsigned char *src);
//** returns the size of the Atomik 3.5 packed data located at src.
//** It does not check for a valid Atomik 3.5 header
//**
//** The function:
//** unsigned long atm_35_origsize(const unsigned char *src)
//** returns the unpacked (original) size of the Atomik 3.5 packed
//** data located at src, including the picture offsets at the end
//** of the file.
//** It does not check for a valid Atomik 3.5 header
//**
//** Atomik 3.5 was a popoular data packer for the Atari ST series
//** Atomik 3.5 packed data can be recognized by the characters "ATM5"
//** at the first 4 positions in a file

unsigned long atm_35_get_long(const unsigned char *p)
{
  unsigned long res;
  res=*p++;
  res<<=8;
  res+=*p++;
  res<<=8;
  res+=*p++;
  res<<=8;
  res+=*p;
  return res;
}

unsigned int atm_35_getbits(int len, int *cmdp, int *maskp, const unsigned char **p)
{
  int tmp=0;
  int cmd=*cmdp;
  int mask=*maskp;
  while(len>0)
  {
    tmp+=tmp;
    mask>>=1;
    if(mask==0)
    {
      *p-=1;
      cmd=**p;
      mask=0x80;
    }
    if(cmd&mask)
    {
      tmp+=1;
    }
    len--;
  }
  *cmdp=cmd;
  *maskp=mask;
  return tmp;
}

int atm_35_header(const unsigned char *src)
{ // returns 0 if no Atomik 3.5 header was found
  return (atm_35_get_long(src)==0x41544d35UL);
}

long atm_35_packedsize(const unsigned char *src)
{ // returns packed size of Atomik packed data (+12 -> including headers!)
  return (long) atm_35_get_long(src+8)+12;
}

long atm_35_origsize(const unsigned char *src)
{ // returns origiginal size of Atomik packed data
  return (long) atm_35_get_long(src+4);
}

long atm_35_depack(const unsigned char *src, unsigned char *dst)
{ // Atomik 3.5 depacker
  unsigned char *p;
  const unsigned char *src_start=src;
  int cmd;
  int mask=0;
  long orig_size;
  int picture_cnt;
  if(!atm_35_header(src))
  { // No 'ATM5' header
    return -1;
  }
  orig_size=atm_35_origsize(src); // orig size
  p=dst+orig_size;
  src+=atm_35_packedsize(src); // packed size
  picture_cnt=*--src; // get numer of compressed pictures in data
  // Atomik has an init problem, the LSB in cmd that is set is _NOT_ valid
  //** reaching this bit is a sign to reload the cmd data (1 byte)
  //** as the msb bit is always set there are always 2 bits set in the cmd
  //** block
  { // fix reload
    cmd=*--src; // init cmd
    mask=0x80;
    while(!(cmd&1))
    { // dump all 0 bits
      cmd>>=1;
      mask>>=1;
    }
    // dump one 1 bit
    cmd>>=1;
  }
  while((p>dst) && (src>src_start))
  {
    unsigned char* bdst=dst;
    int length_bits=0;
    int special_offset_tab_full=0;
    int special_offset_tab[7];
    int special_offset=0;
    int short_char_tab_full=0;
    unsigned char short_char_tab[16];
    // get block data
    if(atm_35_getbits(1, &cmd, &mask, &src))
    { // special offsets
      int i;
      int offset=1;
      special_offset=*--src;
      for(i=0;i<7;i++)
      {
        if(offset==special_offset)
        {
          offset+=2;
        }
        special_offset_tab[i]=offset;
        offset+=2;
      }
      special_offset_tab_full=1;
    }
    if(atm_35_getbits(1, &cmd, &mask, &src))
    { // short chars
      int i;
      for(i=0;i<16;i++)
      {
        short_char_tab[i]=*--src;
      }
      short_char_tab_full=1;
    }
    length_bits=*--src;
    { // get block size
      long int block_size;
      block_size=*--src;
      block_size<<=8;
      block_size+=*--src;
      if((p-dst)<block_size)
      {
        block_size=(long int)(p-dst);
      }
      bdst=p-block_size;
    }
    while((p>bdst) && (src>src_start))
    {
      if(atm_35_getbits(1, &cmd, &mask, &src))
      { // sld
        int len=0;
        int offset=0;
        while(!atm_35_getbits(1, &cmd, &mask, &src))
        {
          if(len==length_bits)
          {
            len=atm_35_getbits(4, &cmd, &mask, &src);
            if(len==0)
            {
              len=atm_35_getbits(8, &cmd, &mask, &src);
              if(len==0)
              {
                len=atm_35_getbits(14, &cmd, &mask, &src);
                len+=255;
              }
              len+=15;
            }
            len+=length_bits;
            break;
          }
          len++;
        }
        if(len==0)
        {
          if(atm_35_getbits(1, &cmd, &mask, &src))
          { // short char
            offset=atm_35_getbits(4, &cmd, &mask, &src);
            if(short_char_tab_full)
            {
              *--p=short_char_tab[offset];
            }
            else
            {
              unsigned char tmp=*p;
              if(offset>7)
              {
                offset+=0xf0;
              }
              *--p=tmp-(unsigned char)offset;
            }
          }
          else
          {
            len=2;
            offset=2;
          }
        }
        else
        {
          len+=2;
          offset=3;
        }
        if(len>0)
        {
          int bits;
          unsigned char *q;
          const int bit_table[]={1,2,4,5,6,7,8,9};
          const int offset_table[]={0,32,96,352,864,1888,3936,8032};
          offset=atm_35_getbits(offset, &cmd, &mask, &src);
          bits=bit_table[offset];
          if(special_offset_tab_full)
          {
            int index;
            bits=atm_35_getbits(bits, &cmd, &mask, &src);
            bits<<=4;
            index=atm_35_getbits(3, &cmd, &mask, &src);
            if(index==7)
            {
              if(atm_35_getbits(1, &cmd, &mask, &src))
              {
                index=atm_35_getbits(3, &cmd, &mask, &src);
                index+=index;
                bits|=index;
              }
              else
              {
                bits|=special_offset;
              }
            }
            else
            {
              bits|=special_offset_tab[index];
            }
            offset=offset_table[offset]+bits;
          }
          else
          {
            bits+=4;
            offset=offset_table[offset]+atm_35_getbits(bits, &cmd, &mask, &src);
          }
          q=p+offset+1;
          if((p-bdst)<len)
          {
            len=(int)(p-bdst);
          }
          while(len>0)
          {
            *--p=*--q;
            len--;
          }
        }
      }
      else
      { // literal
        *--p=*--src;
      }
    }
  }
  if(p>dst)
  { // depack error
    return -1;
  }
  { // pick algo: chunky to 4 bitplanes
    while((picture_cnt>0) && (orig_size>32000))
    {
      int i;
      long offset;
      orig_size-=4;
      offset=atm_35_get_long(dst+orig_size);
      if(offset>(orig_size-32000))
      { // check upper bound
        offset=orig_size-32000;
      }
      if(offset<0)
      { // check lower bound
        offset=0;
      }
      p=dst+offset;
      for(i=0;i<4000;i++)
      {
        unsigned int plane0=0;
        unsigned int plane1=0;
        unsigned int plane2=0;
        unsigned int plane3=0;
        int j;
        for(j=0;j<8;j++)
        {
          unsigned char data;
          data=p[j];
          plane0<<=2;
          plane1<<=2;
          plane2<<=2;
          plane3<<=2;
          plane3+=data&1;
          data>>=1;
          plane2+=data&1;
          data>>=1;
          plane1+=data&1;
          data>>=1;
          plane0+=data&1;
          plane3+=data&2;
          data>>=1;
          plane2+=data&2;
          data>>=1;
          plane1+=data&2;
          data>>=1;
          plane0+=data&2;
        }
        *p++=(unsigned char)(plane0>>8);
        *p++=(unsigned char)plane0;
        *p++=(unsigned char)(plane1>>8);
        *p++=(unsigned char)plane1;
        *p++=(unsigned char)(plane2>>8);
        *p++=(unsigned char)plane2;
        *p++=(unsigned char)(plane3>>8);
        *p++=(unsigned char)plane3;
      }
      picture_cnt--;
    }
  }
  return orig_size;
}
*/
import "C"
import "unsafe"

//import "reflect"
import (
	//	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	//	"io"
	"os"
)

var strFlag, strFolder string //= flag.String("long-string", "DATA.LNK", "Disk Maggie DATA.LNK File")

func init() {
	// example with short version for long flag
	flag.StringVar(&strFlag, "input", "data", "ATOMIK-Packed Data File")
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
	fi, _ := f.Stat()
	data := make([]byte, fi.Size(), fi.Size())

	var source *C.uchar = (*C.uchar)(unsafe.Pointer(&data[0]))
	binary.Read(f, binary.BigEndian, &data)
	unpsize := C.atm_35_origsize(source)
	dataunpacked := make([]byte, unpsize, unpsize)
	var dest *C.uchar = (*C.uchar)(unsafe.Pointer(&dataunpacked[0]))
	fmt.Printf("Unpacked: %d\n", unpsize)
	res := C.atm_35_depack(source, dest)
	if res == -1 {
		panic("Data file is not packed with Atomic 3.5.")
	}
	fmt.Println(res)
	fw, err := os.OpenFile(strFolder+"/"+strFlag, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0764)
	if err != nil {
		panic(err)
	}
	defer fw.Close()
	var theGoSlice []byte = C.GoBytes((unsafe.Pointer)(dest), (C.int)(unpsize))
	/*
		sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&theGoSlice)))
		sliceHeader.Cap = int(unpsize)
		sliceHeader.Len = int(unpsize)
		sliceHeader.Data = uintptr(C.GoByte)
	*/
	fmt.Printf("%c", theGoSlice[0])
	fmt.Printf("%c", theGoSlice[1])
	fw.Write(theGoSlice)
}
