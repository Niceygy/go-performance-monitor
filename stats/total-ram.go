package stats

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mackerelio/go-osstat/memory"
)

func memoryTotal() string {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		//return
	}

	outputMEM := strconv.FormatUint(memory.Total, 10)
	return outputMEM
}
