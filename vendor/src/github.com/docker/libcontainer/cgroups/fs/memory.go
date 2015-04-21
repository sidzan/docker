package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"github.com/docker/libcontainer/cgroups"
	"github.com/Sirupsen/logrus"
)

type MemoryGroup struct {
}

func (s *MemoryGroup) Set(d *data) error {
logrus.Debugf("!!!!!!!!!!!!!     inside memory set")
	dir, err := d.join("memory")
	// only return an error for memory if it was specified
	logrus.Debugf("This issijan!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!above the check funciont")
	if err != nil && (d.c.Memory != 0 || d.c.MemoryReservation != 0 || d.c.MemorySwap != 0) {
		 logrus.Debugf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! inside the check function")
		return err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()
	logrus.Debugf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!1hecking if memory was specified")
	// Only set values if some config was specified.
	if d.c.Memory != 0 || d.c.MemoryReservation != 0 || d.c.MemorySwap != 0 {
	logrus.Debugf("!!!!!!!!!!!!!    above d.c memory")
		if d.c.Memory != 0 {
		logrus.Debugf("!!!!!!!!!!!!!    inside d.c memory")
			if err := writeFile(dir, "memory.limit_in_bytes", strconv.FormatInt(d.c.Memory, 10)); err != nil {
				logrus.Debugf("!!!!!!!!!!!!!    inside d.c memory and writing")
				return err
			}
		}
		if d.c.MemoryReservation != 0 {
			if err := writeFile(dir, "memory.soft_limit_in_bytes", strconv.FormatInt(d.c.MemoryReservation, 10)); err != nil {
				return err
			}
		}
		// By default, MemorySwap is set to twice the size of RAM.
		// If you want to omit MemorySwap, set it to '-1'.
		if d.c.MemorySwap == 0 {
			if err := writeFile(dir, "memory.memsw.limit_in_bytes", strconv.FormatInt(d.c.Memory*2, 10)); err != nil {
				return err
			}
		}
		if d.c.MemorySwap > 0 {
			if err := writeFile(dir, "memory.memsw.limit_in_bytes", strconv.FormatInt(d.c.MemorySwap, 10)); err != nil {
				return err
			}
		}
	} else {
	 logrus.Debugf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! memory was not specified so trying to set default values ")
		writeFile(dir, "memory.limit_in_bytes", "700000000");
		SijanAnanya()
	}
	return nil
}

func (s *MemoryGroup) Remove(d *data) error {
	return removePath(d.path("memory"))
}

func (s *MemoryGroup) GetStats(path string, stats *cgroups.Stats) error {
	logrus.Debugf("!!!!!!!!!!!!!     inside memory getstat")
	// Set stats from memory.stat.
	statsFile, err := os.Open(filepath.Join(path, "memory.stat"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer statsFile.Close()

	sc := bufio.NewScanner(statsFile)
	for sc.Scan() {
		t, v, err := getCgroupParamKeyValue(sc.Text())
		if err != nil {
			return fmt.Errorf("failed to parse memory.stat (%q) - %v", sc.Text(), err)
		}
		stats.MemoryStats.Stats[t] = v
	}

	// Set memory usage and max historical usage.
	value, err := getCgroupParamUint(path, "memory.usage_in_bytes")
	if err != nil {
		return fmt.Errorf("failed to parse memory.usage_in_bytes - %v", err)
	}
	stats.MemoryStats.Usage = value
	value, err = getCgroupParamUint(path, "memory.max_usage_in_bytes")
	if err != nil {
		return fmt.Errorf("failed to parse memory.max_usage_in_bytes - %v", err)
	}
	stats.MemoryStats.MaxUsage = value
	value, err = getCgroupParamUint(path, "memory.failcnt")
	if err != nil {
		return fmt.Errorf("failed to parse memory.failcnt - %v", err)
	}
	stats.MemoryStats.Failcnt = value

	return nil
}

func SijanAnanya() {
	logrus.Debugf("!!!!!calledSijanAnanya")
	fmt.Println("This is going t change thewhole code")

 }
