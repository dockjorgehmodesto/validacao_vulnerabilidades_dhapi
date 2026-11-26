package manager

import "os"

func HostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostName
}
