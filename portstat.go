package nux

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/toolkits/file"
	"github.com/toolkits/slice"
	"github.com/toolkits/sys"
	"io"
	"strconv"
	"strings"
)

func ListeningPorts() ([]int64, error) {
	ports := []int64{}

	bs, err := sys.CmdOutBytes("ss", "-t", "-l", "-n")
	if err != nil {
		return ports, err
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))

	// ignore the first line
	line, err := file.ReadLine(reader)
	if err != nil {
		return ports, err
	}

	for {
		line, err = file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return ports, err
		}

		fields := strings.Fields(string(line))
		fieldsLen := len(fields)

		if fieldsLen != 4 && fieldsLen != 5 {
			return ports, fmt.Errorf("output of [ss -t -l -n] format not supported")
		}

		portColumnIndex := 2
		if fieldsLen == 5 {
			portColumnIndex = 3
		}

		location := strings.LastIndex(fields[portColumnIndex], ":")
		port := fields[portColumnIndex][location+1:]

		if p, e := strconv.ParseInt(port, 10, 64); e != nil {
			return ports, fmt.Errorf("parse port to int64 fail: %s", e.Error())
		} else {
			ports = append(ports, p)
		}

	}

	return slice.UniqueInt64(ports), nil
}
