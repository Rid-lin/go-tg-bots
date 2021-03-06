// Copyright (c) 2020 Xelaj Software
//
// This file is a part of go-dry package.
// See https://github.com/xelaj/go-dry/blob/master/LICENSE for details

package dry

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"strings"
)

func BytesReader(data any) io.Reader {
	switch s := data.(type) {
	case io.Reader:
		return s
	case []byte:
		return bytes.NewReader(s)
	case string:
		return strings.NewReader(s)
	case fmt.Stringer:
		return strings.NewReader(s.String())
	case error:
		return strings.NewReader(s.Error())
	}
	return nil
}

func BytesMD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func BytesEncodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func BytesDecodeBase64(base64Str string) string {
	result, _ := base64.StdEncoding.DecodeString(base64Str)
	return string(result)
}

func BytesEncodeHex(str string) string {
	return hex.EncodeToString([]byte(str))
}

func BytesDecodeHex(hexStr string) string {
	result, _ := hex.DecodeString(hexStr)
	return string(result)
}

func BytesDeflate(uncompressed []byte) (compressed []byte) {
	var buf bytes.Buffer
	writer := Deflate.GetWriter(&buf)
	writer.Write(uncompressed)
	Deflate.ReturnWriter(writer)
	return buf.Bytes()
}

func BytesInflate(compressed []byte) (uncompressed []byte) {
	reader := flate.NewReader(bytes.NewBuffer(compressed))
	result, _ := ioutil.ReadAll(reader)
	return result
}

func BytesGzip(uncompressed []byte) (compressed []byte) {
	var buf bytes.Buffer
	writer := Gzip.GetWriter(&buf)
	writer.Write(uncompressed)
	Gzip.ReturnWriter(writer)
	return buf.Bytes()
}

func BytesUnGzip(compressed []byte) (uncompressed []byte) {
	reader, err := gzip.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return nil
	}
	result, _ := ioutil.ReadAll(reader)
	return result
}

// BytesHead returns at most numLines from data starting at the beginning.
// A slice of the remaining data is returned as rest.
// \n is used to detect line ends, a preceding \r will be stripped away.
// BytesHead resembles the Unix head command.
func BytesHead(data []byte, numLines int) (lines []string, rest []byte) {
	if numLines <= 0 {
		panic("numLines must be greater than zero")
	}
	lines = make([]string, 0, numLines)
	begin := 0
	for i := range data {
		if data[i] == '\n' {
			end := i
			if i > 0 && data[i-1] == '\r' {
				end--
			}
			lines = append(lines, string(data[begin:end]))
			begin = i + 1
			if len(lines) == numLines {
				break
			}
		}
	}
	if len(lines) != numLines {
		lines = append(lines, string(data[begin:]))
		begin = len(data)
	}
	return lines, data[begin:]
}

// BytesTail returns at most numLines from the end of data.
// A slice of the remaining data before lines is returned as rest.
// \n is used to detect line ends, a preceding \r will be stripped away.
// BytesTail resembles the Unix tail command.
func BytesTail(data []byte, numLines int) (lines []string, rest []byte) {
	if numLines <= 0 {
		panic("numLines must be greater than zero")
	}
	lines = make([]string, 0, numLines)
	end := len(data)
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '\n' {
			begin := i
			if end < len(data) && data[end-1] == '\r' {
				end--
			}
			lines = append(lines, string(data[begin+1:end]))
			end = begin
			if len(lines) == numLines {
				break
			}
		}
	}
	if len(lines) != numLines {
		lines = append(lines, string(data[:end]))
		end = 0
	}
	return lines, data[:end]
}

// BytesMap maps a function on each element of a slice of bytes.
func BytesMap(f func(byte) byte, data []byte) []byte {
	size := len(data)
	result := make([]byte, size, size)
	for i := 0; i < size; i++ {
		result[i] = f(data[i])
	}
	return result
}

// BytesFilter filters out all bytes where the function does not return true.
func BytesFilter(f func(byte) bool, data []byte) []byte {
	result := make([]byte, 0, 0)
	for _, element := range data {
		if f(element) {
			result = append(result, element)
		}
	}
	return result
}

var bitlen = []int{
	1 << 3,  // 8
	1 << 4,  // 16
	1 << 5,  // 32
	1 << 6,  // 64
	1 << 7,  // 128
	1 << 8,  // 256
	1 << 9,  // 512
	1 << 10, // 1024
	1 << 11, // 2048
}

func BigIntBytes(v *big.Int, bitsize int) []byte {
	vbytes := v.Bytes()
	vbytesLen := len(vbytes)
	for i, b := range bitlen {
		if b == bitsize {
			break
		}

		if i == len(bitlen)-1 {
			panic(fmt.Errorf("bitsize not squaring by 2: bitsize %v", bitsize))
		}
	}

	offset := bitsize/8 - vbytesLen
	if offset < 0 {
		panic(fmt.Errorf("bitsize too small: have %v, want at least %v", bitsize, vbytes))
	}

	return append(make([]byte, offset), vbytes...)
}

func BigIntBytesAutoChoose(v *big.Int) []byte {
	vbytes := v.Bytes()
	vbytesLen := len(vbytes)
	choosedBitsize := bitlen[0]
	for i, bitsize := range bitlen {
		// 8 ?????? ???????????? ?????????? ?? ??????????
		if len(vbytes)*8 <= bitsize {
			choosedBitsize = bitsize
			break
		}

		// ???????? ???????? ?????????? ?????????????? ???????????????? ?????????????????? ???? ????????????????
		if i == len(bitlen)-1 {
			panic(fmt.Errorf("value is too large: bitsize is %v", vbytesLen*8))
		}
	}

	offset := choosedBitsize/8 - vbytesLen
	return append(make([]byte, offset), vbytes...)
}

func BytesXor(a, b []byte) []byte {
	res := make([]byte, len(a))
	copy(res, a)
	for i := range res {
		res[i] ^= b[i]
	}
	return res
}
