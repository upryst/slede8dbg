# slede8dbg
A debugger for [SLEDE8 VM](https://github.com/PSTNorge/slede8/).

![screenshot](https://raw.githubusercontent.com/julebokk/slede8dbg/main/example/slede8dbg.png)

## Building

1. Download & install [Go](https://golang.org/doc/install).
2. git clone https://github.com/julebokk/slede8dbg.git
3. cd slede8dbg
4. go build .

## Debugging

```
$ ./slede8dbg debug ./example/hello.s8
$ ./slede8dbg debug --input 9090cd219090 ./example/hello.s8
```
