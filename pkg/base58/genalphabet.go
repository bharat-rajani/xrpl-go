//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	start = []byte(`// In compliance with copyfree initiative we are acknowledging that base58 codec were originally developed 
// by btcsuite developers and this package modifies the dictionary as per ripple's directive for base58 encoding and decoding.
// AUTOGENERATED by genalphabet.go; do not edit.

package base58

const (
	// alphabet is the modified base58 alphabet used by Ripple.
	alphabet = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"

	alphabetIdx0 = 'r'
)

var b58 = [256]byte{`)

	end = []byte(`}`)

	alphabet = []byte("rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz")
	tab      = []byte("\t")
	invalid  = []byte("255")
	comma    = []byte(",")
	space    = []byte(" ")
	nl       = []byte("\n")
)

func write(w io.Writer, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fi, err := os.Create("alphabet.go")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	write(fi, start)
	write(fi, nl)
	for i := byte(0); i < 32; i++ {
		write(fi, tab)
		for j := byte(0); j < 8; j++ {
			idx := bytes.IndexByte(alphabet, i*8+j)
			if idx == -1 {
				write(fi, invalid)
			} else {
				write(fi, strconv.AppendInt(nil, int64(idx), 10))
			}
			write(fi, comma)
			if j != 7 {
				write(fi, space)
			}
		}
		write(fi, nl)
	}
	write(fi, end)
	write(fi, nl)
}
