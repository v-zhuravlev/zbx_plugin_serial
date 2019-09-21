package serial

import (
	"testing"
)

var b16le = []byte{0xc2, 0x8f} //36802 in uint16 LE

func TestUint16LE(t *testing.T) {
	x, _ := getUint16(b16le, 0, len(b16le), LittleEndian)
	if x != 36802 {
		t.Errorf("getUint16(*,LittleEndian) = %d; want 36802", x)
	}
}

func TestInt16LE(t *testing.T) {
	x, _ := getInt16(b16le, 0, len(b16le), LittleEndian)
	if x != -28734 {
		t.Errorf("getInt16(*,LittleEndian) = %d; want -28734", x)
	}
}

var b16be = []byte{0x8f, 0xc2} //36802 in uint16 BE

func TestUint16BE(t *testing.T) {
	x, _ := getUint16(b16be, 0, len(b16be), BigEndian)
	if x != 36802 {
		t.Errorf("getUint16(*,BigEndian) = %d; want 36802", x)
	}
}

func TestInt16BE(t *testing.T) {
	x, _ := getInt16(b16be, 0, len(b16be), BigEndian)
	if x != -28734 {
		t.Errorf("getInt16(*,BigEndian) = %d; want -28734", x)
	}
}

/*
32 bits
*/

var b = []byte{0x1b, 0x06, 0x00, 0x00, 0x08, 0x3a, 0x41, 0xbb, 0x01, 0x1b, 0x03}

func TestUint32LE(t *testing.T) {
	x, _ := getUint32(b, 5, len(b), LittleEndian)
	if x != 29049146 {
		t.Errorf("getUint32(*,LittleEndian) = %d; want 29049146", x)
	}
}

func TestUint32BE(t *testing.T) {
	x, _ := getUint32(b, 5, len(b), BigEndian)
	if x != 977386241 {
		t.Errorf("getUint32(*,BigEndian) = %d; want 977386241", x)
	}
}

var b32be = []byte{0xff, 0xff, 0xfd, 0xce}

func TestInt32BE(t *testing.T) {
	x, _ := getInt32(b32be, 0, len(b32be), BigEndian)
	if x != -562 {
		t.Errorf("getInt32(*,BigEndian) = %d; want -562", x)
	}
}

var b32le = []byte{0xce, 0xfd, 0xff, 0xff}

func TestInt32LE(t *testing.T) {
	x, _ := getInt32(b32le, 0, len(b32le), LittleEndian)
	if x != -562 {
		t.Errorf("getInt32(*,LittleEndian) = %d; want -562", x)
	}
}

var bfloat32le = []byte{0x0d, 0xc2, 0x8f, 0xc2}

func TestFloat32LE(t *testing.T) {
	x, _ := getFloat32(bfloat32le, 0, len(bfloat32le), LittleEndian)
	if x != -71.879005 {
		t.Errorf("getFloat32(*,LittleEndian) = %f; want -71.879005", x)
	}
}

var bfloat32be = []byte{0xc2, 0x8f, 0xc2, 0x0d}

func TestFloat32BE(t *testing.T) {
	x, _ := getFloat32(bfloat32be, 0, len(bfloat32be), BigEndian)
	if x != -71.879005 {
		t.Errorf("getFloat32(*,BigEndian) = %f; want -71.879005", x)
	}
}

/*
64 bits
*/

var b64le = []byte{0x7c, 0xf2,
	0xb0, 0x50,
	0x6b, 0x9a,
	0xbf, 0xbf} // 13816931967501922940 in uint64 LE, -4629812106207628676 in int64 le

func TestUint64LE(t *testing.T) {
	x, _ := getUint64(b64le, 0, len(b64le), LittleEndian)
	if x != 13816931967501922940 {
		t.Errorf("getUint64(*,LittleEndian) = %d; want 13816931967501922940", x)
	}
}

func TestInt64LE(t *testing.T) {
	x, _ := getInt64(b64le, 0, len(b64le), LittleEndian)
	if x != -4629812106207628676 {
		t.Errorf("getInt64(*,LittleEndian) = %d; want -4629812106207628676", x)
	}
}

var b64be = []byte{0xbf, 0xbf,
	0x9a, 0x6b,
	0x50, 0xb0,
	0xf2, 0x7c} // 13816931967501922940 in uint64 BE, -4629812106207628676 in int64 BE
func TestUint64BE(t *testing.T) {
	x, _ := getUint64(b64be, 0, len(b64be), BigEndian)
	if x != 13816931967501922940 {
		t.Errorf("getUint64(*,BigEndian) = %d; want 13816931967501922940", x)
	}
}

func TestInt64BE(t *testing.T) {
	x, _ := getInt64(b64be, 0, len(b64be), BigEndian)
	if x != -4629812106207628676 {
		t.Errorf("getInt64(*,BigEndian) = %d; want -4629812106207628676", x)
	}
}

func TestFloat64BE(t *testing.T) {
	x, _ := getFloat64(b64be, 0, len(b64be), BigEndian)
	if x != -0.123450 {
		t.Errorf("getFloat64(*,BigEndian) = %f; want -0.123450", x)
	}
}

func TestFloat64LE(t *testing.T) {
	x, _ := getFloat64(b64le, 0, len(b64le), LittleEndian)
	if x != -0.123450 {
		t.Errorf("getFloat64(*,LittleEndian) = %f; want -0.123450", x)
	}
}

/*
Text and Raw
*/

var bhello = []byte{0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64} //hello world

func TestRaw(t *testing.T) {
	x, _ := getRaw(bhello, 0, len(bhello))
	if x != "68656C6C6F20776F726C64" {
		t.Errorf("getRaw(bhello, 0, len(bhello)) = %s; want 68656C6C6F20776F726C64", x)
	}
}

// TODO fix
// func TestText(t *testing.T) {
// 	fmt.Println(len(bhello))
// 	x, _ := getText(bhello, 0, len(bhello))
// 	if x != "hello world" {
// 		t.Errorf("getText(bhello, 0, len(bhello)) = %s; want 'hello world'", x)
// 	}
// }
