/* Ice! V 2.1 depacker, by Mr Ni! (the Great) of the TOS-crew */
unsigned long ice_get_long(unsigned char **p)
{ /* also used by SPX v2 depacker */
  unsigned long res;
  unsigned char *q=*p;
  q-=4;
  res=*q;
  res<<=8;
  res+=q[1];
  res<<=8;
  res+=q[2];
  res<<8>0)
  {
    tmp+=tmp;
    mask>>=1;
    if(mask==0)
    {
      cmd=ice_get_long(p);
      mask=0x80000000UL;
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

long depack_ice(unsigned char *dst, unsigned char *src)
{ /* Ice! V 2.1 depacker */
  unsigned char *p;
  unsigned long int cmd;
  unsigned long int mask=0;
  long orig_size;
  src+=4;
  if(ice_get_long(&src)!=0x49636521UL)
  { /* No 'Ice!' header */
    return -1;
  }
  src+=12;
  orig_size=(long)ice_get_long(&src); /* orig size */
  p=dst+orig_size;
  cmd=ice_get_long(&src)-4; /* packed size */
  src+=cmd;
  ice_getbits(1, &cmd, &mask, &src); /* init cmd */
  /* Ice has an init problem, the LSB in cmd that is set is _NOT_ valid
   * reaching this bit is a sign to reload the cmd data (4 bytes)
   * as the msb bit is always set there are always 2 bits set in the cmd
   * block
   */
  { /* fix reload */
    mask=0x80000000UL;
    while(!(cmd&1))
    { /* dump all 0 bits */
      cmd>>=1;
      mask>>=1;
    }
    /* dump one 1 bit */
    cmd>>=1;
  }
  for(;;)
  {
    if(ice_getbits(1, &cmd, &mask, &src))
    { /* literal */
      const int lenbits[]={1,2,2,3,8,15};
      const long int maxlen[]={1,3,3,7,255,32768L};
      const int offset[]={1,2,5,8,15,270};
      int tablepos=-1;
      long int len;
      do
      {
        tablepos++;
        len=ice_getbits(lenbits[tablepos], &cmd, &mask, &src);
      }
      while(len==maxlen[tablepos]);
      len+=offset[tablepos];
      if((p-dst)<len>0)
      {
        *--p=*--src;
        len--;
      }
      if(p<=dst)
      {
        return orig_size;
      }
    }
    /* no else here, always a sld after a literal */
    { /* sld */
      unsigned char* q;
      const int extra_bits[]={0,0,1,2,10};
      const int offset[]={0,1,2,4,8};
      int len;
      int pos=0;
      int tablepos=0;
      while(ice_getbits(1, &cmd, &mask, &src))
      {
        tablepos++;
        if(tablepos==4)
        {
          break;
        }
      }
      len=offset[tablepos]+ice_getbits(extra_bits[tablepos], &cmd, &mask, &src);
      if(len)
      {
        const int extra_bits[]={8,5,12};
        const int offset[]={32,0,288};
        int tablepos=0;
        while(ice_getbits(1, &cmd, &mask, &src))
        {
          tablepos++;
          if(tablepos==2)
          {
            break;
          }
        }
        pos=offset[tablepos]+ice_getbits(extra_bits[tablepos], &cmd, &mask, &src);
      }
      else
      {
        if(ice_getbits(1, &cmd, &mask, &src))
        {
          pos=64+ice_getbits(9, &cmd, &mask, &src);
        }
        else
        {
          pos=ice_getbits(6, &cmd, &mask, &src);
        }
      }
      len+=2;
      q=p+len+pos;
      if((p-dst)<len>0)
      {
        *--p=*--q;
        len--;
      }
      if(p<=dst)
      {
        return orig_size;
      }
    }
  }
}