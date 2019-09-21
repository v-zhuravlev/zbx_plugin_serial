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
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"zabbix/pkg/plugin"
	"zabbix/pkg/std"

	"github.com/tarm/serial"
)

// Plugin -
type Plugin struct {
	plugin.Base
}

var impl Plugin
var stdOs std.Os

const readTimeout = 5
const baudDefault = 9600
const sizeDefault = 8

type config struct {
	ConnString string
	Request    string
	ByteToRead int // first byte to read
	Datatype   Datatype
	Endianess  Endianess
}

// Export -
func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	if len(params) < 1 {
		return nil, errors.New("Please provide at least <connection string>, for example in form of /dev/ttyS0 9600 N 8 2")
	}
	if len(params) > 5 {
		return nil, errors.New("Too many parameters")
	}

	c := config{
		ConnString: params[0],
		ByteToRead: 0,
		Datatype:   Raw,
		Endianess:  LittleEndian}

	if len(params) > 1 && len(params[1]) > 0 {
		c.ByteToRead, err = strconv.Atoi(params[1])
		if err != nil {
			return "", fmt.Errorf("Bad byte to start from provided '%s'", params[1])
		}
	}

	if len(params) > 2 {
		c.Request = params[2]
	}

	if len(params) > 3 && len(params[3]) > 0 {
		switch c.Datatype = Datatype(params[3]); c.Datatype {
		case Float:
			//supported
		case Double:
			//supported
		case Raw:
			//supported
		case Text:
			//supported
		case Uint16:
			//supported
		case Uint32:
			//supported
		case Uint64:
			//supported
		case Int16:
			//supported
		case Int32:
			//supported
		case Int64:
			//supported
		default:
			return nil, errors.New("Bad datatype provided")
		}
	}

	if len(params) > 4 && len(params[4]) > 0 {
		switch c.Endianess = Endianess(params[4]); c.Endianess {
		case LittleEndian:
			//supported
		case BigEndian:
			//supported
		default:
			return nil, errors.New("Bad endianess provided, must be LE or BE")
		}
	}

	return getSerial(c)
}

func getSerial(config config) (response string, err error) {

	c, err := getPort(config.ConnString)
	if err != nil {
		return "", err
	}

	s, err := serial.OpenPort(&c)
	if err != nil {
		return "", fmt.Errorf("Failed to open the port %s %d %v %v %v", c.Name, c.Baud, string(c.Parity), c.Size, c.StopBits)
	}

	// fromTo, err := getBytesToRead(config.ByteToRead)
	// if err != nil {
	// 	return "", err
	// }

	var n int
	if len(config.Request) > 0 {

		src := []byte(config.Request)
		dst := make([]byte, hex.DecodedLen(len(src)))
		n, err = hex.Decode(dst, src)
		if err != nil {
			return "", fmt.Errorf("Failed to parse command string '%s'", config.Request)
		}

		n, err = s.Write([]byte(dst[:n]))
		if err != nil {
			return "", errors.New("Failed to write command to the port")
		}
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		return "", errors.New("Timeout: failed to read from the port in time")
	}

	if n < config.ByteToRead {
		return "", errors.New("First byte to read is out of bounds of the reply")
	}

	from := config.ByteToRead

	switch config.Datatype {
	case Text:
		return getText(buf, from, n)
	case Raw:
		return getRaw(buf, from, n)
	case Uint16:

		x, err := getUint16(buf, from, n, config.Endianess)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", x), nil

	case Uint32:

		x, err := getUint32(buf, from, n, config.Endianess)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", x), nil
	case Uint64:

		x, err := getUint64(buf, from, n, config.Endianess)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", x), nil
	case Int16:

		x, err := getInt16(buf, from, n, config.Endianess)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", x), nil
	case Int32:

		x, err := getInt32(buf, from, n, config.Endianess)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", x), nil
	case Int64:

		x, err := getInt64(buf, from, n, config.Endianess)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", x), nil

	case Float:
		x, err := getFloat32(buf, from, n, config.Endianess)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%f", x), nil
	case Double:
		x, err := getFloat64(buf, from, n, config.Endianess)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%f", x), nil
	}

	return "", errors.New("Uknown error")

}

func checkLength(to int, lastIndex int) (err error) {
	if to > lastIndex {
		return errors.New("First byte to read is out of bounds of the reply")
	}
	return nil
}

func getPort(connStr string) (c serial.Config, err error) {

	conn := strings.Split(connStr, " ")

	baud := baudDefault
	if len(conn) > 1 {
		baud, err = strconv.Atoi(conn[1])
		if err != nil {
			return c, fmt.Errorf("Failed to parse baudrate of '%s'", conn[1])
		}
	}

	var parity serial.Parity = serial.ParityNone
	if len(conn) > 2 {
		switch parityChar := conn[2]; parityChar {
		case "N":
			parity = serial.ParityNone
		case "E":
			parity = serial.ParityEven
		case "O":
			parity = serial.ParityOdd
		case "M":
			parity = serial.ParityMark
		case "S":
			parity = serial.ParitySpace
		default:
			return c, fmt.Errorf("Failed to parse parity from '%s', expected 'N','E','O','M' or 'S'", conn[2])
		}
	}

	size := sizeDefault
	if len(conn) > 3 {
		size, err = strconv.Atoi(conn[3])
		if err != nil {
			return c, fmt.Errorf("Failed to parse databits size from '%s'", conn[3])
		}
	}

	var stopBits serial.StopBits = serial.Stop2
	if len(conn) > 4 {

		switch stopBitsStr := conn[4]; stopBitsStr {
		case "1":
			stopBits = serial.Stop1
		case "2":
			stopBits = serial.Stop2
		case "15":
			stopBits = serial.Stop1Half
		default:
			return c, fmt.Errorf("Failed to parse stopbits from '%s', expected '1', '2' or '15'", conn[4])
		}
	}

	c = serial.Config{
		Name:        conn[0],
		Baud:        baud,
		Parity:      parity,
		Size:        byte(size),
		StopBits:    stopBits,
		ReadTimeout: time.Second * readTimeout,
	}
	return c, nil

}

func init() {
	stdOs = std.NewOs()
	plugin.RegisterMetrics(&impl, "Serial simple get", "serial.get", "Simple request/response to serial port.")
}
