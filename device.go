package keypresslog

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unsafe"
)

type DeviceReader interface {
	Read() (chan Event, error)
	GetId() int
	GetName() string
}

type device struct {
	id   int
	name string
}

func Find() []DeviceReader {
	var ret []DeviceReader

	for i := 0; i < MAX_FILES; i++ {
		buff, err := ioutil.ReadFile(fmt.Sprintf(EVENTS_PATH_TPL, i))
		if err != nil {
			break
		}
		ret = append(ret, newDeviceReader(buff, i))
	}

	return ret
}

func (d *device) GetName() string {
	return d.name
}

func (d *device) GetId() int {
	return d.id
}

func (d *device) Read() (chan Event, error) {
	ret := make(chan Event)
	eventSize := int(unsafe.Sizeof(Event{}))
	fd, err := os.Open(fmt.Sprintf(DEVICE_FILE_TPL, d.GetId()))

	switch {
	case os.IsPermission(err):
		close(ret)
		fd.Close()
		return ret, errors.New("You don't have permissions to read this device")
	case err != nil:
		close(ret)
		fd.Close()
		return ret, err
	}

	go func() {
		tmp := make([]byte, eventSize)

		for {

			n, err := fd.Read(tmp)

			if err != nil {
				close(ret)
				fd.Close()
				panic(err)
			}

			if n <= 0 {
				continue
			}

			var event Event
			if err := binary.Read(bytes.NewBuffer(tmp), binary.LittleEndian, &event); err != nil {
				close(ret)
				fd.Close()
				panic(err)
			}

			ret <- event
		}
	}()

	return ret, nil
}

func newDeviceReader(buff []byte, id int) DeviceReader {
	rd := bufio.NewReader(bytes.NewReader(buff))
	rd.ReadLine()
	dev, _, _ := rd.ReadLine()
	split := strings.Split(string(dev), "=")

	return DeviceReader(&device{id, split[1]})
}