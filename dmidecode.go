package dmidecode

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	//"regexp"
	//"strings"
)

type DMI struct {
	Data map[string]map[string]string
}

func New() *DMI {
	dmi := &DMI{}
	return dmi
}

// Wrapper for FindBin, ExecCmd, ParseDmidecode
func (d *DMI) Run() error {
	bin, findErr := d.FindBin("dmidecode")
	if findErr != nil {
		return findErr
	}

	cmdOutput, cmdErr := d.ExecDmidecode(bin)
	if cmdErr != nil {
		return cmdErr
	}

	if err := d.ParseDmidecode(cmdOutput); err != nil {
		return err
	}

	return nil
}

func (d *DMI) FindBin(binary string) (string, error) {
	locations := []string{"/sbin", "/usr/sbin", "/usr/local/sbin"}

	for _, path := range locations {
		lookup := path + "/" + binary
		fileInfo, err := os.Stat(path + "/" + binary)

		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			return lookup, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Unable to find the '%v' binary", binary))
}

func (d *DMI) ExecDmidecode(binary string) (string, error) {
	cmd := exec.Command(binary)

	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (d *DMI) ParseDmidecode(output string) error {
	// Each record is separated by double newlines

	return nil
}
