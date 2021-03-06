# Serial zabbix agent plugin

This plugin for zabbix-agent2(exporter) is for getting data from the serial port.

## Build

Make sure golang is [installed](https://golang.org/doc/install) and properly configured.

Checkout zabbix branch with zabbix-agent2:  
`git clone https://git.zabbix.com/scm/zbx/zabbix.git --depth 1 zabbix-agent2`  
`cd zabbix-agent2`  
Checkout this plugin repo:  
`git submodule add https://github.com/v-zhuravlev/zbx_plugin_serial.git src/go/plugins/serial`  

Edit file `src/go/plugins/plugins.go` by appending `_ "zabbix.com/plugins/serial"`:

```go
package plugins

import (
	_ "zabbix.com/plugins/log"
	_ "zabbix.com/plugins/systemrun"
	_ "zabbix.com/plugins/zabbix/async"
	_ "zabbix.com/plugins/zabbix/stats"
	_ "zabbix.com/plugins/zabbix/sync"
	_ "zabbix.com/plugins/serial"
)
```

`./bootstrap.sh`   
`pushd .`  
`cd src/go`  
`go mod vendor`  
`popd`  
`./configure --enable-agent2`   
`make`  

You will then find new agent with plugin included in `src/go/bin` dir

Test it by running
`./zabbix-agent2 -t agent.ping -c ../conf/zabbix_agent2.conf`

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

## Next steps

- Add resource locking, make sure single serial port is not accessed simuletaneusly
- Make read timeout configurable
- Add new item, that will implement Watcher pattern
- More tests coverage

## Examples

```text
zabbix_get -s localhost -k serial.get["/dev/ttyr00",5,1b02081b03,uint32,LE]
31832326
```

```text
zabbix_get -s localhost -k serial.get["/dev/ttyr00",0,1b02081b03,raw]
1B0600000841B9E5011B03
```

## Changelog

### v0.1

Initial version
