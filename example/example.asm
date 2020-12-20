; Syntax should be fully compatible with https://slede8.npst.no/
;
; In addition, .DATA can take strings enclosed in "" and chars enclosed in ''

    HOPP start

hello:
    .DATA "Hello world!", 0

start:
    FINN hello
    SETT r11, 1

next:
    LAST r5
    LIK r5, r10
    BHOPP kthxbye
    SKRIV r5
    PLUSS r0, r11
    HOPP next

kthxbye:
    STOPP
