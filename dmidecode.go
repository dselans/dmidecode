package dmidecode

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	DMIDecodeBinary = "dmidecode"
)

type Record map[string]string

type DMI struct {
	Data   map[string][]Record
	Binary string
}

func New() *DMI {
	return &DMI{
		Data:   make(map[string][]Record, 0),
		Binary: DMIDecodeBinary,
	}
}

// Run will attempt to find a a valid `dmidecode` bin, attempt to execute it and
// parse whatever data it gets.
func (d *DMI) Run() error {
	bin, err := d.FindBin(d.Binary)
	if err != nil {
		return err
	}

	output, err := d.ExecDmidecode(bin)
	if err != nil {
		return err
	}

	return d.ParseDmidecode(output)
}

// FindBin will attempt to find a given binary in common bin paths.
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

	return "", fmt.Errorf("Unable to find the '%v' binary", binary)
}

// ExecDmiDecode will attempt to execute a given binary, capture its output and
// return it (or an any errors it encounters)
func (d *DMI) ExecDmidecode(binary string) (string, error) {
	cmd := exec.Command(binary)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// ParseDmiDecode will attempt to parse dmidecode output and place all matching
// content in d.Data.
func (d *DMI) ParseDmidecode(output string) error {
	// Each record is separated by double newlines
	splitOutput := strings.Split(output, "\n\n")

	for _, record := range splitOutput {
		recordElements := strings.Split(record, "\n")

		// Entries with less than 3 lines are incomplete/inactive; skip them
		if len(recordElements) < 3 {
			continue
		}

		handleRegex, _ := regexp.Compile("^Handle\\s+(.+),\\s+DMI\\s+type\\s+(\\d+),\\s+(\\d+)\\s+bytes$")
		handleData := handleRegex.FindStringSubmatch(recordElements[0])

		if len(handleData) == 0 {
			continue
		}

		dmiHandle := handleData[1]

		r := Record{}
		r["DMIType"] = handleData[2]
		r["DMISize"] = handleData[3]

		// Okay, we know 2nd line == name
		r["DMIName"] = recordElements[1]

		inBlockElement := ""
		inBlockList := ""

		// Loop over the rest of the record, gathering values
		for i := 2; i < len(recordElements); i++ {
			// Check whether we are inside a \t\t block
			if inBlockElement != "" {
				inBlockRegex, _ := regexp.Compile("^\\t\\t(.+)$")
				inBlockData := inBlockRegex.FindStringSubmatch(recordElements[i])

				if len(inBlockData) > 0 {
					if len(inBlockList) == 0 {
						inBlockList = inBlockData[1]
					} else {
						inBlockList = inBlockList + "\t\t" + inBlockData[1]
					}
					r[inBlockElement] = inBlockList
					continue
				} else {
					// We are out of the \t\t block; reset it again, and let
					// the parsing continue
					inBlockElement = ""
				}
			}

			recordRegex, _ := regexp.Compile("\\t(.+):\\s+(.+)$")
			recordData := recordRegex.FindStringSubmatch(recordElements[i])

			// Is this the line containing handle identifier, type, size?
			if len(recordData) > 0 {
				r[recordData[1]] = recordData[2]
				continue
			}

			// Didn't match regular entry, maybe an array of data?
			recordRegex2, _ := regexp.Compile("\\t(.+):$")
			recordData2 := recordRegex2.FindStringSubmatch(recordElements[i])

			if len(recordData2) > 0 {
				// This is an array of data - let the loop know we are inside
				// an array block
				inBlockElement = recordData2[1]
				continue
			}
		}

		d.Data[dmiHandle] = append(d.Data[dmiHandle], r)
	}

	if len(d.Data) == 0 {
		return fmt.Errorf("unable to parse 'dmidecode' output")
	}

	return nil
}

// GenericSearchBy will search for any param w/ value in the d.Data map.
func (d *DMI) GenericSearchBy(param, value string) ([]Record, error) {
	if len(d.Data) == 0 {
		return nil, fmt.Errorf("DMI data is empty; make sure to .Run() first")
	}

	var records []Record

	for _, v := range d.Data {
		for _, d := range v {
			if d[param] != value {
				continue
			}

			records = append(records, d)
		}
	}

	return records, nil
}

// SearchByName will search for a specific DMI record by name in d.Data
func (d *DMI) SearchByName(name string) ([]Record, error) {
	return d.GenericSearchBy("DMIName", name)
}

// SearchByType will search for a specific DMI record by its type in d.Data
func (d *DMI) SearchByType(id int) ([]Record, error) {
	return d.GenericSearchBy("DMIType", strconv.Itoa(id))
}
