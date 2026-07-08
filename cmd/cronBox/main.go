package main

import (
	"cronBox/executor"
	"cronBox/logger"
	"cronBox/parser"
	"cronBox/scheduler"
	"fmt"
)

func main() {
	jobs, err := parser.ParseConfig("config.json")

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Println("Parsed jobs:")
	for _, job := range jobs {
		fmt.Printf("- %s\n", job)
	}

	executor := executor.New()
	logger, err := logger.New("cronbox.log")

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	scheduler := scheduler.New(executor, logger)

	scheduler.AddJobs(jobs)
	scheduler.Start()

	select {}
}
