package stats

import (
	"fmt"
	"os"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
)

func cpuOut() float64 {
	//if error run below
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		//return
	} //get usage over one second (below)
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		//return
	} //print cpu usage
	total := float64(after.Total - before.Total)
	//fmt.Printf("cpu use: %.2f %%\n", float64(after.User-before.User)/total*100) //print usage
	return float64(after.User-before.User) / total * 100

}
