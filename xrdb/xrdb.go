package xrdb

import (
	"fmt"
	"os/exec"
	"strings"
)

func Get(name string) (string, error) {
	cmd := exec.Command("xrdb", "-get", name)

	stderr := &strings.Builder{}
	cmd.Stderr = stderr

	b, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to get value from xrdb:%s: %w", stderr.String(), err)
	}
	return strings.TrimSpace(string(b)), nil
}
