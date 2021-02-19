/*
 * Copyright (C) linvon
 * Date  2021/2/18 10:29
 */

package cuckoo

import (
	"encoding/binary"
	"fmt"
)

const DEBUG = false

type PermEncoding struct {
	nEnts    uint
	DecTable []uint16
	EncTable []uint16
}

func (p *PermEncoding) Init() {
	p.nEnts = 3876
	p.DecTable = make([]uint16, p.nEnts)
	p.EncTable = make([]uint16, 1<<16)

	dst := [4]uint8{}
	var idx uint16
	p.genTables(0, 0, dst, &idx)
}

/* unpack one 2-byte number to four 4-bit numbers */
func (p *PermEncoding) unpack(in uint16, out *[4]uint8) {
	out[0] = uint8(in & 0x000f)
	out[2] = uint8((in >> 4) & 0x000f)
	out[1] = uint8((in >> 8) & 0x000f)
	out[3] = uint8((in >> 12) & 0x000f)
}

/* pack four 4-bit numbers to one 2-byte number */
func (p *PermEncoding) pack(in [4]uint8) uint16 {
	var in1, in2 uint16
	in1 = binary.LittleEndian.Uint16([]byte{in[0], in[1]}) &0x0f0f
	in2 = binary.LittleEndian.Uint16([]byte{in[2], in[3]}) << 4

	return in1 | in2 
}

func (p *PermEncoding) Decode(codeword uint16, lowBits *[4]uint8) {
	p.unpack(p.DecTable[codeword], lowBits)
}
func (p *PermEncoding) Encode(lowBits [4]uint8) uint16 {
	if DEBUG {
		fmt.Printf("Perm.encode\n")
		for i := 0; i < 4; i++ {
			fmt.Printf("encode lowBits[%d]=%x\n", i, lowBits[i])
		}
		fmt.Printf("pack(lowBits) = %x\n", p.pack(lowBits))
		fmt.Printf("enc_table[%x]=%x\n", p.pack(lowBits), p.EncTable[p.pack(lowBits)])
	}
	return p.EncTable[p.pack(lowBits)]
}

func (p *PermEncoding) genTables(base, k int, dst [4]uint8, idx *uint16) {
	for i := base; i < 16; i++ {
		/* for fast comparison in binary_search in little-endian machine */
		dst[k] = uint8(i)
		if k+1 < 4 {
			p.genTables(i, k+1, dst, idx)
		} else {
			p.DecTable[*idx] = p.pack(dst)
			p.EncTable[p.pack(dst)] = *idx
			if DEBUG {
				fmt.Printf("enc_table[%04x]=%04x\t%x %x %x %x\n", p.pack(dst), *idx, dst[0],
					dst[1], dst[2], dst[3])
			}
			*idx++
		}
	}
}
