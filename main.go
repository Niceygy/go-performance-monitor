package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	//VSC really does like a space here
	_ "github.com/go-sql-driver/mysql"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

func memoryUsed() string {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		//return
	}

	//fmt.Printf("memory total: %d bytes\n", memory.Total)
	//fmt.Printf("memory used: %d bytes\n", memory.Used)
	//fmt.Printf("memory cached: %d bytes\n", memory.Cached)
	outputMEM := strconv.FormatUint(memory.Used, 10)
	return outputMEM
}

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

func diskOut() {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		panic(err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)

	fmt.Printf("Total: %d bytes\n", total)
	fmt.Printf("Free: %d bytes\n", free)
}

func dbConnect() {
	db, err := sql.Open("mysql", "go:go@tcp(127.0.0.1:3306)/go")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}

func setMID() string {
	return time.Now().Format("20060102150405")
}

// func addData() bool {
// 	res, err := db.Exec("INSERT INTO go(MID, MNAME, CPU, RAM, DISK) VALUES (?, ?)", setMID(), getHostname(), cpuOut(), memoryOut(), diskOut())
// 	if err != nil {
// 		panic(err.Error())
// 		return false
// 	}
// 	return true
// }

func main() {

	fmt.Println("Loading!")
	dbConnect()
	fmt.Println("Hostname: " + getHostname())
	fmt.Println("MID: " + setMID())

	i := 0
	fmt.Println("Loaded!")
	for {
		//fmt.Println(i) //debug
		i++
		memoryOut()
		//cpuOut()
		fmt.Println(cpuOut())
		time.Sleep(time.Duration(1) * time.Second)

	}

}
