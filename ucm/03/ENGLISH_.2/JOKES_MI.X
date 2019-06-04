                                  JOKES MIX
                                  ---------

-------------------------------------------------------------------------------

How many programmers are necessary to change a light bulb? None, it's a problem
of hardware...

-------------------------------------------------------------------------------

Two nuns are out of fuel on a lonely road. Luckily there is a farm nearby where
they get gasoline. But the farmer does not have a can and gives them a chamber-
pot. When the nuns are refueling the car a passing driver stops saying: "I wan-
na have your believe!"

-------------------------------------------------------------------------------

Multimedia on ATARI: Using Tempus Word  to print out a 3D-picture computed with
Seb's method, while  playing  the  Heavy-module "Tribute to Priest" from Tim on
the DSP-Replay for what I can't understand anything of "Silence of Lambs" (Hope
the title is translated correctly... )

-------------------------------------------------------------------------------

The ATARI-TOS: I had to look at the routine  calculating the  checksum of boot-
sectors - no problem you think?

in the TOS ( A5 = address of the bootsektor, D6 = it's length ):

         clr.w    D7
         bra.s    test
loop:    add.w    (A5)+,D7
test:    move.w   D6,D0
         subq.w   #1,D6
         tst.w    D0
         bne.s    loop
         move.w   D7,D0

Here the solution every assembler programmer will offer:

         clr.w    D0
         subq.w   #1,D6
loop:    add.w    (A5)+,D0
         dbra     D6,loop

The TOS seems to be a typical  product of  a C-compiler. And  possibly even not
TURBO-C, since  in the  TOS is used  much LINK/UNLK... I do not want to see any
other ( possibly  time-critical ) ATARI-routine. I'd advice a manual optimizing
of the code!

-------------------------------------------------------------------------------

                                                     Apache of tnb / 27.03.1994
