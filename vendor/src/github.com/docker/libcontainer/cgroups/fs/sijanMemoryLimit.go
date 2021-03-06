// This package is maintained by Sijan Shrestha <sijanshrestha2@gmail.com>. This package helps u decide if u want to set default value for the cgroups. I dont take resposnibilty for any damage due to modification of the below code
package fs

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/memoryLimitBySijan"
	"os"
	"strings"
	"strconv"
)

func SetToDefaultValue(d *data) {
	dir, err := d.join("memory")

	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	file, err := os.Open("/etc/default/docker")
	if err != nil {
		defaultfunction(d)
	} else {

		data := make([]byte, 10000)
		file.Read(data)
		s := string(data)
		start := strings.Index(s, "MEMDEFAULT") + 12
		end := strings.Index(s, "DEFAULTMEM") - 1

		option := s[start:end]
		if option == "default" {
			defaultfunction(d)

		} else {
			writeFile(dir, "memory.limit_in_bytes", option)
		}
	}
	//this is going to set default values

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

func defaultfunction(d *data) {
	dir, err := d.join("memory")

	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	logrus.Debugf("go do default")
	si := memoryLimitBySijan.Get()
	TotalMemory := si.TotalRam - 300
	LimitForEachContainer := TotalMemory * 20 / 100
	ByteConverter := 1000 * LimitForEachContainer
	var a int64
	a = Num64(ByteConverter)
	str := strconv.FormatInt(a, 10)
	writeFile(dir, "memory.limit_in_bytes", str)
	logrus.Debugf("!!!!!setting value to %v", str)

}

/*func main() {
	file, err := os.Open("/etc/default/docker")
	check(err)
	data := make([]byte, 10000)
	file.Read(data)
	check(err)
	s := string(data)
	start := strings.Index(s, "MEMDEFAULT") + 12
	end := strings.Index(s, "DEFAULTMEM") - 1
	option := s[start:end]

	if option == "default" {
		fmt.Println("go do default")
	} else {
		fmt.Println(option)
		fmt.Println(reflect.TypeOf(option))

	}
}
*/
