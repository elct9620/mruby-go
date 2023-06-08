package mruby

import (
	"io"
	"os"
	"os/exec"
)

func Compile(script io.Reader) ([]byte, error) {
	temp, err := os.CreateTemp("", "rb")
	if err != nil {
		return nil, err
	}
	defer os.Remove(temp.Name())
	io.Copy(temp, script)

	cmd := exec.Command(
		"mrbc",
		"-o",
		"-",
		temp.Name(),
	)
	bin, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return bin, nil
}
