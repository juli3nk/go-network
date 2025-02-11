package iwgetid

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/martinlindhe/unit"
)

func run(intf, flag string) (string, error) {
	execPath, err := exec.LookPath("iwgetid")
	if err != nil {
		return "", err
	}

	var args []string

	if len(intf) > 0 {
		args = append(args, intf)
	}
	if len(flag) > 0 {
		args = append(args, flag)
	}

	out, err := exec.Command(execPath, args...).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func GetSSID(intf string) (*string, error) {
	ssid, err := run(intf, "-r")
	if err != nil {
		return nil, err
	}

	return &ssid, nil
}

func GetAccessPointMAC(intf string) (*string, error) {
	apm, err := run(intf, "-a")
	if err != nil {
		return nil, err
	}

	return &apm, nil
}

func GetChannel(intf string) (*int, error) {
	ch, err := run(intf, "-c")
	if err != nil {
		return nil, err
	}
	chi, _ := strconv.Atoi(ch)

	return &chi, nil
}

func GetFrequency(intf string) (*float64, error) {
	freq, err := run(intf, "-f")
	if err != nil {
		return nil, err
	}
	freqFloat, _ := strconv.ParseFloat(freq, 64)
	frequency := unit.Frequency(freqFloat).Hertz()

	return &frequency, nil
}
