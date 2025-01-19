package main

import (
	"cpi-count/pkg/checker"
	"fmt"
	"log"
	"runtime"
)

func init() {
	// Pinning to a specific OS thread ensures that the measurements
	// using both the `go-perf` library and Linux `perf` produce
	// similar measurements.
	// Interestingly, although undesired, if this step is omitted:
	// - Linux perf reports ~4x the number of cycles AND instructions
	// - The CPI remains relatively similar because both cycles and instructions increase proportionally
	runtime.LockOSThread()
}

func main() {

	// Check if resctrl is mounted
	if !checker.CheckResctrlMount() {
		fmt.Println("Resctrl is not mounted, attempting to mount it...")

		// Specify the options you want to use for mounting (e.g., "cdp,mba_MBps")
		options := "cdp,mba_MBps" // Adjust options as needed
		err := checker.MountResctrl(options)
		if err != nil {
			log.Fatalf("Error mounting resctrl: %v", err)
		}
		fmt.Println("Resctrl filesystem mounted successfully")
	} else {
		fmt.Println("Resctrl filesystem is already mounted")
	}

	perfCmd := NewPerfCmd()
	if err := perfCmd.Start(); err != nil {
		log.Fatalf("Failed to execute perf cmd: %v\n", err)
	}

	goperf := NewGoPerf()
	if err := goperf.StartWorkloadMeasurement(); err != nil {
		log.Fatalf("Failed to start goperf measurement: %v\n", err)
	}

	goperfOutput, err := goperf.End()
	if err != nil {
		log.Fatalf("Failed to end goperf measurement: %v\n", err)
	}

	perfOutput, err := perfCmd.End()
	if err != nil {
		log.Fatalf("Failed to end perf cmd: %v\n", err)
	}

	log.Printf("GoPerf Cycles: %d, GoPerf Instrs: %d, GoPerf CPI: %f\n", int64(goperfOutput.Cycles), int64(goperfOutput.Instrs), goperfOutput.CPI())
	log.Printf("PerfCmd Cycles: %d, PerfCmd Instrs: %d, PerfCmd CPI: %f\n", int64(perfOutput.Cycles), int64(perfOutput.Instrs), perfOutput.CPI())
}
