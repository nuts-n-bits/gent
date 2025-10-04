package main

import (
	"fmt"
)

type DecoderWorktable struct {
	secFdByte int;
	secLdByte int;
	secBodyByte int;
	inputQ [][]byte;
	inputSumLen int;	
}

func NewDefault() DecoderWorktable {
	return DecoderWorktable {
		secFdByte: 8,
		secLdByte: 8,
		secBodyByte: 1<<30,
		inputQ: make([][]byte, 0),
		inputSumLen: 0,
	}
}

func NewMaxBig() DecoderWorktable {
		return DecoderWorktable {
		secFdByte: 8,
		secLdByte: 8,
		secBodyByte: 1<<56,
		inputQ: make([][]byte, 0),
		inputSumLen: 0,
	}
}

func (dwt *DecoderWorktable) inputPeek(i int) (byte, bool) {
	if dwt.inputSumLen <= i { 
		return 0, false; 
	}
	currentChunkIndex := 0;
	inputPartialSumLen := len(dwt.inputQ[0]);
	for i > inputPartialSumLen {
		currentChunkIndex += 1;
		if (currentChunkIndex >= len(dwt.inputQ)) {
			return 0, false;
		}
		inputPartialSumLen += len(dwt.inputQ[currentChunkIndex]);
	}
	currentChunk := dwt.inputQ[currentChunkIndex];
	return currentChunk[i - inputPartialSumLen + len(currentChunk)], true;
}

func (dwt *DecoderWorktable) Load(chunk []byte) {
	if len(chunk) == 0 {
		return;
	}
	dwt.inputQ = append(dwt.inputQ, chunk);
	dwt.inputSumLen += len(chunk);
}

func (dwt *DecoderWorktable) Unload() [][]byte {
	temp := dwt.inputQ;
	dwt.inputQ = make([][]byte, 0);
	dwt.inputSumLen = 0;
	return temp;
}

// It returns: 
// 1) A frame id, 
// 2) The contents of the frame, 
// 3) An integer != 0 indicating that one stepping cannot be completed
// 4) An error object
// It will only either return (1 and 2) or (3) or (4)
func (dwt *DecoderWorktable) Step() (frameId uint64, oneDataChunkFragmented [][]byte, signal int, err error) {
	if dwt.inputSumLen == 0 {
		return 0, make([][]byte, 0), -1, nil;
	}
	fdBytes := make([]byte, 0);
	i := 0;
	for {
		if i >= dwt.secFdByte {
			return 0, make([][]byte, 0), 0, fmt.Errorf("frame descriptor over length limit %d", dwt.secFdByte);
		}
		varintByte, exists := dwt.inputPeek(i);
		if !exists {
			return 0, make([][]byte, 0), -1, nil;
		}
		fdBytes = append(fdBytes, varintByte);
		if (varintByte & 0x80) == 0 {
			break;
		}
		i += 1;
	}
	ldBytes := make([]byte, 0);
	i = 0;
	for {
		if i >= dwt.secLdByte {
			return 0, make([][]byte, 0), 0, fmt.Errorf("length descriptor over length limit %d", dwt.secLdByte);
		}
		varintByte, exists := dwt.inputPeek(i + len(fdBytes));
		if !exists {
			return 0, make([][]byte, 0), -1, nil;
		}
		ldBytes = append(ldBytes, varintByte);
		if (varintByte & 0x80 == 0) { 
			break;
		}
		i += 1;
	}
	frameId, bodylen := decodeVarint(fdBytes), int(decodeVarint(ldBytes));
	if bodylen > dwt.secBodyByte {
		return 0, make([][]byte, 0), 0, fmt.Errorf("body over length limit %d", dwt.secBodyByte);
	}
	if dwt.inputSumLen < (len(fdBytes) + len(ldBytes) + bodylen) {
		return 0, make([][]byte, 0), (bodylen + len(fdBytes) + len(ldBytes)), nil;
	}
	dwt.unload(len(fdBytes) + len(ldBytes));
	body := dwt.unload(bodylen)
	return frameId, body, 0, nil;
}

func (dwt *DecoderWorktable) StepCompact() (oneDataChunkFragmented [][]byte, signal int, err error) {
	if dwt.inputSumLen == 0 {
		fmt.Printf("444")
		return make([][]byte, 0), -1, nil;
	}
	ldBytes := make([]byte, 0);
	i := 0;
	for {
		if i >= dwt.secLdByte {
			return make([][]byte, 0), 0, fmt.Errorf("length descriptor over length limit %d", dwt.secLdByte);
		}
		varintByte, exists := dwt.inputPeek(i);
		if !exists {
			return make([][]byte, 0), -1, nil;
		}
		ldBytes = append(ldBytes, varintByte);
		if (varintByte & 0x80 == 0) { 
			break;
		}
		i += 1;
	}
	bodylen := int(decodeVarint(ldBytes));
	if bodylen > dwt.secBodyByte {
		return make([][]byte, 0), 0, fmt.Errorf("body over length limit %d", dwt.secBodyByte);
	}
	if dwt.inputSumLen < (len(ldBytes) + bodylen) {
		return make([][]byte, 0), (bodylen + len(ldBytes)), nil;
	}
	dwt.unload(len(ldBytes));
	body := dwt.unload(bodylen)
	return body, 0, nil;
}

// danger: this will always assume inputQ has enough bytes to unload!
func (dwt *DecoderWorktable) unload(unloadlen int) [][]byte {
	if unloadlen == 0 {
		return make([][]byte, 0);
	}
	coll := make([][]byte, 0);
	deficit := unloadlen;
	for deficit > 0 {
		firstBufLen := len(dwt.inputQ[0]);
		if deficit >= firstBufLen {
			deficit -= firstBufLen;
			coll = append(coll, dwt.inputQ[0]);
			dwt.inputQ = dwt.inputQ[1:];
		} else {
			coll = append(coll, dwt.inputQ[0][0:deficit]);
			dwt.inputQ[0] = dwt.inputQ[0][deficit:];
			deficit = 0;
		}
	}
	dwt.inputSumLen -= unloadlen;
	return coll;
}

func encodeVarint(varint uint64) []byte {
	if varint == 0 { 
		return []byte{0};
	}
	coll := make([]byte, 0);
	for varint != 0 {
		coll = append(coll, uint8((varint & 0x7F) | 0x80));
		varint = varint >> 7;
	}
	coll[len(coll)-1] &= 0x7F;
	return coll;
}

func encodeFrameHeader(varint uint64, bytes []byte) []byte {
	coll := make([]byte, 0);
	coll = append(coll, encodeVarint(varint)...);
	coll = append(coll, encodeVarint(uint64(len(bytes)))...);
	return coll;
}

func DelimitCompactNoRem(buffers [][]byte) ([][][]byte, error) {
	dwt := NewDefault();
	for _, buf := range buffers {
		dwt.Load(buf);
	}
	coll := make([][][]byte, 0);
	for {
		if dwt.Remainder() == 0 {
			break;
		}
		buf, signal, err := dwt.StepCompact();
		if err != nil {
			return make([][][]byte, 0), err;
		}
		if signal != 0 {
			return make([][][]byte, 0), fmt.Errorf("incomplete message detected, signal = %d", signal);
		}
		coll = append(coll, buf);
	}
	return coll, nil;
}

func CompatRepr(buffers [][]byte) []byte {
	coll := make([]byte, 0);
	for _, buffer := range buffers {
		lengthDesc := encodeVarint(uint64(len(buffer)));
		coll = append(coll, lengthDesc...);
		coll = append(coll, buffer...);
	}
	return coll;
}

func (dwt *DecoderWorktable) Remainder() int {
	return dwt.inputSumLen;
}

// assumes incoming buffer's varint bounds is correct. assume at least 1 byte
func decodeVarint(bytes []byte) uint64 {
	ret := uint64(0);
	for i, byte := range bytes {
		ret |= uint64(byte & 0x7F) << (i * 7);
	}
	return ret;
}

func subsets(nums []int) [][]int {
    coll := [][]int{};
    
	var search func(start int, acc []int);
    search = func (start int, acc []int) {
        if start == len(nums) {
            coll = append(coll, acc);
            return;
        }
        search(start + 1, acc);
        search(start + 1, append(acc[:], nums[start]));
    };
	
    return coll;
}