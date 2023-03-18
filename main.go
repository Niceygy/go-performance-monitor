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

func cpuOut() string {
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
	floatOut := float64(after.User-before.User) / total * 100
	stringOut := strconv.FormatFloat(floatOut, 'f', -1, 64)
	return stringOut
}

func memoryUsed() string {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		//return
	}

	outputMEM := strconv.FormatUint(memory.Used, 10)
	return outputMEM
}

func memoryTotal() string {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		//return
	}

	outputMEM := strconv.FormatUint(memory.Total, 10)
	return outputMEM
}

func diskOut() string {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		panic(err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)

	fmt.Printf("Total: %d bytes\n", total)
	fmt.Printf("Free: %d bytes\n", free)
	percentFree := (float64(total-free) / float64(total)) * 100
	outputDISK := strconv.FormatUint(uint64(percentFree), 10)
	return outputDISK
}

func dbConnect() {
	db, err := sql.Open("mysql", "go:go@tcp(127.0.0.1:3306)/go")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Connected to the database!")
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

func addData() bool {
	res, err := db.Exec("INSERT INTO go(MID, MNAME, CPU, RAM_TOTAL, RAM_USED, DISK) VALUES (?, ?)", setMID(), getHostname(), cpuOut(), memoryTotal(), memoryUsed(), diskOut())
	if err != nil {
		panic(err.Error())
		return false
	}
	return true
}

func load() {
	fmt.Println("Loading! Please wait...")
	dbConnect()
	addData()
}

func main() {
	load()
}
