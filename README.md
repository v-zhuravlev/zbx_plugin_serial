# Serial zabbix agent plugin

This plugin for zabbix-agent2 is for getting data from the serial port.

## Install

Run 
`usermod -a -G dialout zabbix`
So zabbix-agent can access serial ports.

## Supported keys

### serial.get

`serial.get[<connection string>,<first_byte_to_read>,<request>,<datatype>,<endianess>]`

where

`connection string`  
Serial connection parameters in a form of:  
`portname [baudrate] [parity:N|E|O] [databits] [stopbits:1|2|15]`  
for example  
/dev/ttyS0 9600 N 8 2  
/dev/ttyUSB0 115200 E 8 1  
or enter only the portname, defaults for the rest will be used:  
/dev/ttyS1  
defaults are: 9600 N 8 1

`first_byte_to_read`  
First byte to read from the response. Useful together when retriving numeric datatypes such as `float` or `uint64`.

`request`  
Optional request in hex form (i.e. `10FBAC`) that must written to port before reading the port. If empty, no command is send.

`datatype`

Currently supported datatypes

- `raw` -  hex string
- `text`-  ASCII decoded string
- `uint16`, `uint32`, `uint64`
- `int16`, `int32`, `int64`
- `float`, `double`

If not set, `raw` is used.

`endianess`

Byteorder, used for numeric values extraction. Can be `LE` for LittleEndian and `BE` for BigEndian. If not set, LittleEndian is used.

Example keys:

```text
    serial.get[/dev/ttyS0,32,4,3]
    serial.get["/dev/ttyS0 9600 N 8 2",10,"6f",float]
    serial.get["/dev/pts/2 9600 N 8 2",5,1b02081b03,uint32,LE]
```

## Changelog

v0.1
Initial version
