package sys

import (
	"bytes"
	"os/exec"
)

func CmdOutBytes(args ...string) ([]byte, error) {
	var params []string
	params = append(params, "-c")
	for _, value := range args {
		params = append(params, value)
	}

	cmd := exec.Command("sh", params...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}
