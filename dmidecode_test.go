package dmidecode

import (
	"os"
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

	if err := dmi.ParseDmidecode(output); err != nil {
		t.Error("Should not receive an error after parsing dmidecode output")
	}

	if len(dmi.Data) == 0 {
		t.Error("Parsed data structure should have more than 0 entries")
	}
}
