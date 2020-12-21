# slede8dbg
A debugger / assembler for [SLEDE8 VM](https://github.com/PSTNorge/slede8/).

![screenshot](https://raw.githubusercontent.com/julebokk/slede8dbg/main/example/slede8dbg.png)

## Building

1. Download & install [Go](https://golang.org/doc/install).
2. `git clone https://github.com/julebokk/slede8dbg.git`
3. `cd slede8dbg`
4. `go build .`

## Debugger

```
$ ./slede8dbg debug ./example/hello.s8
$ ./slede8dbg debug ./example/example.asm # compiles it for you
$ ./slede8dbg debug --input 9090cd219090 ./example/hello.s8
```

Using alternative syntax for `debug`:
```
$ ./slede8dbg ./example/hello.s8
$ ./slede8dbg ./example/hello.s8 f09f8e85      # with føde
$ ./slede8dbg ./example/hello.s8 f09f8e85 2600 # and cycle limit
```

## Assembler

```
$ ./slede8dbg compile ./example/example.asm # default binary name is a.s8
$ ./slede8dbg compile -o example.s8 ./example/example.asm
```
