package module

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const osReleaseFile = "/etc/os-release"

type SystemVersion struct {
	SystemName string
	Version    string
}

// GetLinuxVersion 获取 Linux 系统名称和版本信息
func GetLinuxVersion() (SystemVersion, error) {
	var sys SystemVersion
	file, err := os.Open(osReleaseFile)
	if err != nil {
		return sys, fmt.Errorf("%s file opend failed, error: %v", osReleaseFile, err)
	}
	defer file.Close()

	var name, version string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NAME=") {
			name = strings.Trim(line[len("NAME="):], `"`)
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(line[len("VERSION_ID="):], `"`)
		}
	}

	if err := scanner.Err(); err != nil {
		return sys, fmt.Errorf("file scanner error: %v", err)
	}

	if name == "" || version == "" {
		return sys, fmt.Errorf("could not find OS name or version")
	}

	sys.SystemName = name
	sys.Version = version

	return sys, nil
}
