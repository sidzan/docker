package fs

import (
	//"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	//"github.com/docker/libcontainer/cgroups"
	"github.com/memoryLimitBySijan"
	"os"
	//"path/filepath"

	"strconv"
)

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
