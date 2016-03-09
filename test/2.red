;redcode-NW
;name Ritual Sword
;author inversed
;strategy Oneshot / Imp
;assert 1

d	equ	(CORESIZE / 60)
step	equ	-5 * d
zofs	equ	-10 * d
aval	equ	-8 * d

x0	equ	scan
org	scan+1

scan	add	#step-1,	ptr
	jmz.f	scan,		>ptr

ptr	mov	>aval,		x0+zofs
	spl	-1,		<ptr
	mov.i	#1,		1
