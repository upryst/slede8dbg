; SLEDE8++ is an extended version of SLEDE8 computer. It is equipped
; with an experimental VGA graphics chip. VGA was named after the
; initials of the three NPST employees working on the project:
; Vegard, Gunnar and Adrian (it was later revealed that Vegard's
; real name was Vladimir; he was fired for breaching the non-compete
; agreement).

    SETT r11, 1

loop:
    SETT r5, r0
    PLUSS r5, r4
    SETT r6, r1
    PLUSS r6, r4
    XELLER r5, r6
    VLAGR r5
    PLUSS r0, r11
    ULIK r0, r10
    BHOPP loop
    PLUSS r1, r11
    ULIK r1, r10
    BHOPP loop
    VSYNK
    PLUSS r4, r11
    HOPP loop
