/*
** Zabbix
** Copyright (C) 2001-2019 Zabbix SIA
**
** This program is free software; you can redistribute it and/or modify
** it under the terms of the GNU General Public License as published by
** the Free Software Foundation; either version 2 of the License, or
** (at your option) any later version.
**
** This program is distributed in the hope that it will be useful,
** but WITHOUT ANY WARRANTY; without even the implied warranty of
** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
** GNU General Public License for more details.
**
** You should have received a copy of the GNU General Public License
** along with this program; if not, write to the Free Software
** Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
**/

package serial

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//Datatype to return
type Datatype string

//Datatypes currently supported
const (
	Float  Datatype = "float"
	Double Datatype = "double"
	Uint16 Datatype = "uint16"
	Uint32 Datatype = "uint32"
	Uint64 Datatype = "uint64"
	Int16  Datatype = "int16"
	Int32  Datatype = "int32"
	Int64  Datatype = "int64"
	Raw    Datatype = "raw"
	Text   Datatype = "text"
)

//Endianess supported
type Endianess string

//Endianess currently supported
const (
	LittleEndian Endianess = "LE"
	BigEndian    Endianess = "BE"
)

func getUint16(buf []byte, from int, last int, end Endianess) (result uint16, err error) {

	to := from + 2
	err = checkLength(to, last)
	if err != nil {
		return result, err
	}
	if end == LittleEndian {
		result = binary.LittleEndian.Uint16(buf[from:to])
	} else if end == BigEndian {
		result = binary.BigEndian.Uint16(buf[from:to])
	}

	return result, nil
}

func getUint32(buf []byte, from int, last int, end Endianess) (result uint32, err error) {

	to := from + 4
	err = checkLength(to, last)
	if err != nil {
		return result, err
	}
	if end == LittleEndian {
		result = binary.LittleEndian.Uint32(buf[from:to])
	} else if end == BigEndian {
		result = binary.BigEndian.Uint32(buf[from:to])
	}

	return result, nil
}

func getUint64(buf []byte, from int, last int, end Endianess) (result uint64, err error) {

	to := from + 8
	err = checkLength(to, last)
	if err != nil {
		return result, err
	}
	if end == LittleEndian {
		result = binary.LittleEndian.Uint64(buf[from:to])
	} else if end == BigEndian {
		result = binary.BigEndian.Uint64(buf[from:to])
	}

	return result, nil
}

func getInt16(buf []byte, from int, last int, end Endianess) (result int16, err error) {

	x, err := getUint16(buf, from, last, end)
	if err != nil {
		return result, err
	}

	return int16(x), nil

}

func getInt32(buf []byte, from int, last int, end Endianess) (result int32, err error) {

	x, err := getUint32(buf, from, last, end)
	if err != nil {
		return result, err
	}

	return int32(x), nil

}

func getInt64(buf []byte, from int, last int, end Endianess) (result int64, err error) {

	x, err := getUint64(buf, from, last, end)
	if err != nil {
		return result, err
	}

	return int64(x), nil

}

func getFloat32(buf []byte, from int, last int, end Endianess) (result float32, err error) {

	to := from + 4
	err = checkLength(to, last)
	if err != nil {
		return result, err
	}

	bufr := bytes.NewReader(buf[from:to])

	if end == LittleEndian {
		err := binary.Read(bufr, binary.LittleEndian, &result)
		if err != nil {
			return result, err
		}
	} else if end == BigEndian {
		err := binary.Read(bufr, binary.BigEndian, &result)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func getFloat64(buf []byte, from int, last int, end Endianess) (result float64, err error) {

	to := from + 8
	err = checkLength(to, last)
	if err != nil {
		return result, err
	}

	bufr := bytes.NewReader(buf[from:to])

	if end == LittleEndian {
		err := binary.Read(bufr, binary.LittleEndian, &result)
		if err != nil {
			return result, err
		}
	} else if end == BigEndian {
		err := binary.Read(bufr, binary.BigEndian, &result)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func getRaw(buf []byte, from int, last int) (result string, err error) {
	return fmt.Sprintf("%X", buf[from:last]), nil
}

func getText(buf []byte, from int, last int) (result string, err error) {
	return fmt.Sprintf("%q", buf[from:last]), nil
}
