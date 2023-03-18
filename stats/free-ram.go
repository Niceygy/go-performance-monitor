package stats

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mackerelio/go-osstat/memory"
)

func memoryUsed() string {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		//return
	}

	outputMEM := strconv.FormatUint(memory.Used, 10)
	return outputMEM
}
