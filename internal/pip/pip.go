package pip

import (
	"os/exec"
	"strings"
)

func GetPipInstallList() ([]PipInstall, error) {
	cmd := exec.Command("pip", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	got := make([]PipInstall, 0)
	lists := strings.Split(string(output), "\n")
	if len(lists) < 3 {
		return got, nil
	}
	for _, value := range lists[2:] {
		if len(value) < 1 {
			continue
		}
		tv := strings.Trim(value, " ")
		dd := strings.Split(tv, " ")
		if len(dd) < 2 {
			continue
		}
		got = append(got, PipInstall{
			Name:    dd[0],
			Version: dd[len(dd)-1],
		})
	}
	return got, nil
}

type PipInstall struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
