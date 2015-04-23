package memoryLimitBySijan

import (
	"fmt"

	"sync"
	"syscall"
	"time"
)

type SI struct {
	Uptime       time.Duration // time since boot
	Loads        [3]float64    // 1, 5, and 15 minute load averages, see e.g. UPTIME(1)
	Procs        uint64        // number of current processes
	TotalRam     uint64        // total usable main memory size [kB]
	FreeRam      uint64        // available memory size [kB]
	SharedRam    uint64        // amount of shared memory [kB]
	BufferRam    uint64        // memory used by buffers [kB]
	TotalSwap    uint64        // total swap space size [kB]
	FreeSwap     uint64        // swap space still available [kB]
	TotalHighRam uint64        // total high memory size [kB]
	FreeHighRam  uint64        // available high memory size [kB]
	mu           sync.Mutex    // ensures atomic writes; protects the following fields
}

var sis = &SI{}

func Get() *SI {

	si := &syscall.Sysinfo_t{}

	err := syscall.Sysinfo(si)
	if err != nil {
		panic("Commander Sijan, we have a problem. syscall.Sysinfo:" + err.Error())
	}
	scale := 65536.0 // magic
	defer sis.mu.Unlock()
	sis.mu.Lock()
	unit := uint64(si.Unit) * 1024 // kB
	sis.Uptime = time.Duration(si.Uptime) * time.Second
	sis.Loads[0] = float64(si.Loads[0]) / scale
	sis.Loads[1] = float64(si.Loads[1]) / scale
	sis.Loads[2] = float64(si.Loads[2]) / scale
	sis.Procs = uint64(si.Procs)
	sis.TotalRam = uint64(si.Totalram) / unit
	sis.FreeRam = uint64(si.Freeram) / unit
	sis.BufferRam = uint64(si.Bufferram) / unit
	sis.TotalSwap = uint64(si.Totalswap) / unit
	sis.FreeSwap = uint64(si.Freeswap) / unit
	sis.TotalHighRam = uint64(si.Totalhigh) / unit
	sis.FreeHighRam = uint64(si.Freehigh) / unit
	return sis
}
func (si SI) String() string {
	// XXX: Is the copy of SI done atomic? Not sure.
	// Without an outer lock this may print a junk.
	return fmt.Sprintf("uptime\t\t%v\nload\t\t%2.2f %2.2f %2.2f\nprocs\t\t%d\n"+
		"ram total\t%d kB\nram free\t%d kB\nram buffer\t%d kB\n"+
		"swap total\t%d kB\nswap free\t%d kB",
		//"high ram total\t%d kB\nhigh ram free\t%d kB\n"
		si.Uptime, si.Loads[0], si.Loads[1], si.Loads[2], si.Procs,
		si.TotalRam, si.FreeRam, si.BufferRam,
		si.TotalSwap, si.FreeSwap,
	// archaic si.TotalHighRam, si.FreeHighRam
	)
}
func (si *SI) ToString() string {
	defer si.mu.Unlock()
	si.mu.Lock()
	return si.String()
}
