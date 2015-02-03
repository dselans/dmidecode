package dmidecode

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

const (
	fakeBinary string = "time4soup"
)

func TestFindBin(t *testing.T) {
	dmi := New()

	if _, err := dmi.FindBin("time4soup"); err == nil {
		t.Error("Should not be able to find obscure binary")
	}

	bin, findErr := dmi.FindBin("dmidecode")
	if findErr != nil {
		t.Errorf("Should be able to find dmidecode. Error: %v", findErr)
	}

	_, statErr := os.Stat(bin)

	if statErr != nil {
		t.Errorf("Should be able to lookup found file. Error: %v", statErr)
	}
}

func TestExecDmidecode(t *testing.T) {
	dmi := New()

	if _, err := dmi.ExecDmidecode("/bin/" + fakeBinary); err == nil {
		t.Errorf("Should get an error trying to execute a fake binary. Error: %v", err)
	}

	bin, findErr := dmi.FindBin("dmidecode")
	if findErr != nil {
		t.Errorf("Should be able to find binary. Error: %v", findErr)
	}

	output, execErr := dmi.ExecDmidecode(bin)

	if execErr != nil {
		t.Errorf("Should not get errors executing '%v'. Error: %v", bin, execErr)
	}

	if len(output) == 0 {
		t.Errorf("Output should not be empty")
	}
}

func TestParseDmidecode(t *testing.T) {
	dmi := New()

	bin, findErr := dmi.FindBin("dmidecode")
	if findErr != nil {
		t.Errorf("Should be able to find binary. Error: %v", findErr)
	}

	output, execErr := dmi.ExecDmidecode(bin)

	if execErr != nil {
		t.Errorf("Should not get errors executing '%v'. Error: %v", bin, execErr)
	}

	splitOutput := strings.Split(output, "\n\n")

	for _, record := range splitOutput {
		recordElements := strings.Split(record, "\n")

		// Only care about entries which have 3+ lines
		if len(recordElements) < 3 {
			continue
		}

		handleRegex, _ := regexp.Compile("^Handle\\s+(.+),\\s+DMI\\s+type\\s+(\\d+),\\s+(\\d+)\\s+bytes$")
		handleData := handleRegex.FindStringSubmatch(recordElements[0])

		if len(handleData) == 0 {
			continue
		}

		// dmiHandle := handleData[1]
		// dmiType := handleData[2]
		// dmiSize := handleData[3]

		// okay, we know 2nd line == name
		dmiName := recordElements[1]

		// Loop over the rest of the record, gathering values
		for i := 2; i < len(recordElements); i++ {
			t.Errorf("We have %v elements", len(recordElements))
			recordRegex, _ := regexp.Compile("\\t(.+):\\s+(.+)$")
			recordData := recordRegex.FindStringSubmatch(recordElements[i])

			if len(recordData) > 0 {
				t.Errorf("Found data for dmiName %v: %v => %v", dmiName, recordData[1], recordData[2])
			}
		}
	}
}
