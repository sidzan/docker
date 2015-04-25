package fs

import (
	"bufio"
	"github.com/Sirupsen/logrus"
	"github.com/docker/libcontainer/cgroups"
	"os"
	"path/filepath"
	"strconv"
)

type CpuGroup struct {
}

func (s *CpuGroup) Set(d *data) error {
	// We always want to join the cpu group, to allow fair cpu scheduling
	// on a container basis
	dir, err := d.join("cpu")
	if err != nil {
		return err
	}
	if d.c.CpuShares != 0 {
		if err := writeFile(dir, "cpu.shares", "100"); err != nil {
			return err
		}
	} else {
		SijanCpuShareLimit(d)
	}
	if d.c.CpuPeriod != 0 {
		if err := writeFile(dir, "cpu.cfs_period_us", strconv.FormatInt(d.c.CpuPeriod, 10)); err != nil {
			return err
		}
	}
	if d.c.CpuQuota != 0 {
		if err := writeFile(dir, "cpu.cfs_quota_us", strconv.FormatInt(d.c.CpuQuota, 10)); err != nil {
			return err
		}
	}
	return nil
}

func (s *CpuGroup) Remove(d *data) error {
	return removePath(d.path("cpu"))
}

func (s *CpuGroup) GetStats(path string, stats *cgroups.Stats) error {
	f, err := os.Open(filepath.Join(path, "cpu.stat"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		t, v, err := getCgroupParamKeyValue(sc.Text())
		if err != nil {
			return err
		}
		switch t {
		case "nr_periods":
			stats.CpuStats.ThrottlingData.Periods = v

		case "nr_throttled":
			stats.CpuStats.ThrottlingData.ThrottledPeriods = v

		case "throttled_time":
			stats.CpuStats.ThrottlingData.ThrottledTime = v
		}
	}
	return nil
}
func SijanCpuShareLimit(d *data) {
	logrus.Debugf("cccccccccccccccccccccccccccccccccccccc this was called")
	dir, err := d.join("cpu")
	writeFile(dir, "cpu.shares", "100")
}
