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
	_, err = io.Copy(temp, script)
	if err != nil {
		return nil, err
	}

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
