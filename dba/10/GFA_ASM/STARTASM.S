; How to get started in assembler.
; Program coded by O.T.M. of the D.B.A.
; This is a very simple example.
;
	section	text

	bsr.s	init
	bsr	main
	bsr	exit

init:
	dc.w	$a00a		; kill the mouse

	clr.l	-(sp)		; enter supervisor
	move.w	#$20,-(sp)	; (<--- 1 word and ^--- 1 longword)
	trap	#1		; Gemdos
	addq.l	#6,sp		; (1 word = 2 bytes and 1 longw = 4 bytes)
	move.l	d0,oldstack	; save this important value

	move.w	#4,-(sp)	; this routine checks the resolution
	trap	#14		; xbios --> in G.f.a.: oldrez=xbios(4)
	addq.l	#2,sp		; correct stack: 1 word : 2 bytes
	move.w	d0,oldrez	; save the value

* The value that is returned by the previous routine contains the 
* following information
* 0 - Low resolution
* 1 - Medium Resolution
* 2 - High resolution	(Monochrome monitor connected)
* Higher values are reserved for TT-resolutions
	
	cmp.w	#2,d0		; test the resolution
	bne.s	rez_good	; if it's not equal then continue

	pea	err_text(pc)	; else give an error message
	move.w	#9,-(sp)	; (on screen)
	trap	#1		; using Gemdos
	addq.l	#6,sp		; correct stack

	move.w	#7,-(sp)	; wait for a key
	trap	#1		; Gemdos
	addq.l	#2,sp		; correct stack

	bra	clear_off	; and quit nice and neatly
rez_good:			; and continue
	move.w	#2,-(sp)	; get the screen-address
	trap	#14		; it's an Xbios (compare it with G.f.a.)
	addq.l	#2,sp		; (It would be like this oldscr=XBIOS(2))
	move.l	d0,oldscr	; save the returned value

* We are now going to clear the screen without using Gemdos, Bios or
* Xbios
* Before we can do that you'll have to know that the screen is build up
* out of 32000 bytes. 
* The screenaddress can be found in d0 coz of the previous routine
* We put it in an addressregister, and we will work our way through
* the memorie until the whole screen is cleared.

	movea.l	d0,a0		; move it to an address register
	move.w	#8000-1,d0	; use d0 as a counter (8000 longwords
				; is equal to 32000 bytes)
cls:
	clr.l	(a0)+		; clear a longword and go to the next
	dbf	d0,cls		; do until d0 is "FALSE" (-1)


* done it.....

	move.w	#0,-(sp)	; enter low resolution-------------v
	move.l	#-1,-(sp)	; keep screenaddress------------v
	move.l	#-1,-(sp)	; keep second screenaddress-v
	move.w	#5,-(sp)	; the routine	--------v        
	trap	#14		; In G.f.a.: VOID XBIOS(5,l:-1,l:-1,0)
	add.l	#12,sp		; correct the stack	

	move.w	#1,-(sp)	; init some musix by our own musix-man
	pea	song		; Nice yamaha tunes by our own SONIC
	jsr	routine+6	;
	addq.l	#6,sp		; correct the stack again
	
	rts			; all inits done 

main:
	bsr.s	vsync		; wait till the screen is finished
	bsr.s	vu_meter	; our own routine

	move.w	#11,-(sp)	; test if a key is pressed
	trap	#1		; using Gemdos
	addq.l	#2,sp

	tst.w	d0		; if d0 is zero then no key is pressed
	beq.s	main		; no ---> continue
	rts			; else stop and go back

vsync:
	movem.l	d0-a7,-(sp)	; save all registers
				; (Xbios messes em up)

* First we change the colour, we wait for the screen to finish
* and then we change the background colour again to see how much
* time we have left 
* Also called rastertimer.

	move.w	#0,-(sp)	; colour to change it into (black)
	move.w	#0,-(sp)	; this is the colour we wanna change
	move.w	#7,-(sp)	; This xbios routine is equal to setcolor
	trap	#14		; xbios
	addq.l	#6,sp		; correct stack: 3 word == 6 bytes

	move.w	#$25,-(sp)	; official vsync routine
	trap	#14		; xbios	
	addq.l	#2,sp		; correct the stack
	
	move.w	#$777,-(sp)	; change the colour again (white)
	move.w	#0,-(sp)	; the backgroundcolour
	move.w	#7,-(sp)	; setcolor
	trap	#14		; xbios routine
	addq.l	#6,sp		; correct the stack

	movem.l	(sp)+,d0-a7	; give the registers back
	rts			; and back to main

vu_meter:			; our own routine

	movem.l	d0-a7,-(sp)	; save registers

	move.w	#8,-(sp)	; G.f.a.:volume1=xbios(28,0,8)
	move.w	#0,-(sp)
	move.w	#28,-(sp)	; xbios routine to check upon the musix
	trap	#14		; gives the volume of musix-channel 1
	addq.l	#6,sp		; correct the stack
	move.w	d0,volume1	; save the value

	move.w	#9,-(sp)	; get the volume of channel 2
	move.w	#0,-(sp)	; read it 
	move.w	#28,-(sp)	
	trap	#14
	addq.l	#6,sp		; correct stack
	move.w	d0,volume2	; save the value

	move.w	#10,-(sp)	; get the volume of channel 2
	move.w	#0,-(sp)	; read it 
	move.w	#28,-(sp)	
	trap	#14
	addq.l	#6,sp		; correct stack
	move.w	d0,volume3	; save the value

* I think it's best to take the routines below for granted. 
* To understand them you have to know a lot more about the screen
* and assembler.
* And I think it is not very interesting for you to know yet.
* It would only bore you, and you wouldn't code anymore because you
* would be so depressed, and conviced that coding is difficult that 
* you wouldn't continue coding, and that would be a terrible shame, 
* coz coding isn't that difficult. It only needs some pursistance 
* and normal common sense.

	move.l	oldscr,a0	; the screenaddress
	moveq	#15,d0	; the counter
clearline1:
	clr.w	(a0)		; clear the contense if this addres
	addq.l	#8,a0		; next address 
	dbf	d0,clearline1

	move.l	oldscr,a0	; the screenaddress
	lea	$5*$a0(a0),a0	; the offset
	moveq	#15,d0	; counter
clearline2:
	clr.w	(a0)		; clear it
	addq.l	#8,a0
	dbf	d0,clearline2	

	move.l	oldscr,a0	; same screen
	lea	$a*$a0(a0),a0	; another ofset
	moveq	#15,d0
clearline3:
	clr.w	(a0)
	addq.l	#8,a0
	dbf	d0,clearline3

	move.l	oldscr,a0	; re-install the screenaddress
	move.w	volume1,d0	; use d0 as a counter
newline1:
	move.w	#-1,(a0)	; and form a new line
	addq.l	#8,a0		; next line in same plane
	dbf	d0,newline1
	
	move.l	oldscr,a0
	lea	$5*$a0(a0),a0
	move.w	volume2,d0
newline2:
	move.w	#-1,(a0)
	addq.l	#8,a0
	dbf	d0,newline2
	
	move.l	oldscr,a0
	lea	$a*$a0(a0),a0
	move.w	volume3,d0
newline3:
	move.w	#-1,(a0)	
	addq.l	#8,a0	
	dbf	d0,newline3

* So, and that was it for the "difficult" part.
* let's continue our lesson

	movem.l	(sp)+,d0-a7	; give the stored addresses back
	rts			; return back to main

exit:				; and now to quit everything
	clr.l	-(sp)		; shut up ze musix
	clr.w	-(sp)		; total silence
	bsr	routine+6	; from here
	addq.l	#6,sp		; correct stack

	move.w	#$777,-(sp)	; colour to change it into (white)
	move.w	#0,-(sp)	; this is the colour we wanna change
	move.w	#7,-(sp)	; This xbios routine is equal to setcolor
	trap	#14		; xbios
	addq.l	#6,sp		; correct stack: 3 word == 6 bytes

	move.w	oldrez,-(sp)	; change back to old-resolution
	move.l	#-1,-(sp)	; keep screen-address
	move.l	#-1,-(sp)
	move.w	#5,-(sp)	; change it
	trap	#14		; xbios
	add	#12,sp		; correct stack


clear_off:			; in case of wrong rez we jump to this 
				; label to quit in a good way
	move.l	oldstack,-(sp)	; enter user-mode again
	move.w	#$20,-(sp)	; set it
	trap	#1		; Gemdos
	addq.l	#6,sp		; correct stack	

	dc.w	$a009		; make mouse active

	clr.w	-(sp)		; quit program, and return to desktop
	trap	#1		; (When assembled)

	section	data

routine:incbin	'f:\music\sonic\include.sng\music_qu.inc'

song:	incbin	'f:\music\sonic\itgoeson.mod'

err_text:
	dc.b	27,'E','Sorry monochrome is not supported so reconnect a colour-monitor.',13,10,0

	section bss

oldstack:
	ds.l 1
oldscr:	ds.l 1
oldrez:	ds.w 1
volume1:ds.w 1
volume2:ds.w 1
volume3:ds.w 1


	end