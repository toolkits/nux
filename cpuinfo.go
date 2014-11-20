package nux

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strings"
)

func NumCpu() int {
	return runtime.NumCPU()
}

func CpuMHz() (mhz string, err error) {
	f := "/proc/cpuinfo"
	var bs []byte
	bs, err = ioutil.ReadFile(f)
	if err != nil {
		return
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))
	var line []byte

	for {
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			return
		}

		li := string(line)
		if !strings.Contains(li, "MHz") {
			continue
		}

		arr := strings.Split(li, ":")
		if len(arr) != 2 {
			return "", fmt.Errorf("%s content format error", f)
		}

		return strings.TrimSpace(arr[1]), nil
	}

	return "", fmt.Errorf("no MHz in %s", f)
}
