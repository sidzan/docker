// +build !linux

package systemd

import (
	"fmt"	
	"github.com/Sirupsen/logrus"
	"github.com/docker/libcontainer/cgroups"
)

func UseSystemd() bool {
	log.Debugf("ccccccccccccccccccccccc setupCgroups called  calling use systemd?")
	return false
}

func Apply(c *cgroups.Cgroup, pid int) (map[string]string, error) {
	log.Debugf("ccccccccccccccccccccccc setupCgroups called  aplly called from systemd")
	return nil, fmt.Errorf("Systemd not supported")
}

func GetPids(c *cgroups.Cgroup) ([]int, error) {
	log.Debugf("ccccccccccccccccccccccc setupCgroups called  get pidcalled")
	return nil, fmt.Errorf("Systemd not supported")
}

func ApplyDevices(c *cgroups.Cgroup, pid int) error {
	return fmt.Errorf("Systemd not supported")
}

func Freeze(c *cgroups.Cgroup, state cgroups.FreezerState) error {
	return fmt.Errorf("Systemd not supported")
}
