package main

import (
	"database/sql"
	"fmt"
	iou "io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	//VSC really does like a space here

	_ "github.com/go-sql-driver/mysql"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

//<-------------------Data Gathering Functions------------------->

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

//<-------------------Data Gathering Functions End------------------->

//<-------------------Database Functions------------------->

func UpdateDataDB(mid, hostname, cpu, ram_free, ram_total, disk_free string) bool {
	db, err := sql.Open("mysql", "go:go@tcp(127.0.0.1:3306)/go")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Connected to the database!")
	}
	defer db.Close()
	res, err := db.Exec("INSERT INTO stats(MID, MNAME, CPU, RAM_TOTAL, RAM_USED, DISK) VALUES (?, ?)", mid, hostname, cpu, ram_free, ram_total, disk_free)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
	return true
}

//<-------------------Database Functions End------------------->
//<-------------------Single run fucntions------------------->

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}

func setMID() string {
	body, err := iou.ReadFile("/etc/gpm/mid")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	fmt.Println(string(body))
	if string(body) == "" {
		fmt.Println("No MID found, generating one...")
		iou.WriteFile("/etc/gpm/mid", []byte(time.Now().Format("20060102150405")), 0644)
		return time.Now().Format("20060102150405")
	} else {
		return string(body)
	}
}

func firstRunCheck() bool {
	_, err := os.Stat("/etc/gpm/mid")
	if os.IsNotExist(err) {
		fmt.Println("First run detected, generating MID...")
		setMID()
		return true
	} else {
		fmt.Println("MID found, continuing...")
		return false

	}
}

func load() {
	fmt.Println("Loading! Please wait...")
	firstRunCheck()
}

func main() {
	load()
	for {
		UpdateDataDB(setMID(), getHostname(), cpuOut(), memoryUsed(), memoryTotal(), diskOut())
		time.Sleep(time.Second * 5)
		fmt.Println("Updated! Sleeping for 5 seconds...")
	}
}
