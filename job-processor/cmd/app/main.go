package main

import (
	"fmt"

	"github.com/thutasann/job-processor/internal/job"
)

func main() {
	j := job.New("1", job.High, func() error {
		fmt.Println("Executing job 1")
		return nil
	})

	fmt.Println(j)
}
