package fs

import (
	"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/docker/libcontainer/cgroups"
	//"github.com/memoryLimitBySijan"
	"os"
	"path/filepath"

	"strconv"
)

type MemoryGroup struct {
}

func (s *MemoryGroup) Set(d *data) error {
	logrus.Debugf("!!!!!!!!!!!!!     inside memory set")
	dir, err := d.join("memory")
	// only return an error for memory if it was specified
	logrus.Debugf("XXXXXXXXXX!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!above the check funciont")
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
		//

		SijanAnanya(d)
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

/*
func SijanAnanya(d *data) {
	dir, err := d.join("memory")

	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	logrus.Debugf("!!!!!calledSijanAnanya")
	fmt.Println("This is going t change thewhole code")
	si := memoryLimitBySijan.Get()
	TotalMemory := si.TotalRam
	logrus.Debugf("!!!!!!!!!!!!!!!!!!calledSijanAnanya%v\n", si.TotalRam)
	//fmt.Printf("%v\n", si.TotalRam)
	//	logrus.Debugf(reflect.TypeOf(si.TotalRam))
	LimitForEachContainer := TotalMemory * 20 / 100
	ByteConverter := 1000000 * LimitForEachContainer
	var a int64
	a = Num64(ByteConverter)
	str := strconv.FormatInt(a, 10)
	writeFile(dir, "memory.limit_in_bytes", str)
	//s	str := strconv.FormatUInt(ByteConverter, 10)

	//fmt.Println(reflect.TypeOf(strval))
}
func Num64(n interface{}) int64 {
	s := fmt.Sprintf("%d", n)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	} else {
		return i
	}
}
*/
