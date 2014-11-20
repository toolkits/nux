package nux

import (
	"fmt"
	"testing"
)

func TestMetrics(t *testing.T) {

	// kernel
	fmt.Print("KernelMaxFiles:")
	fmt.Println(KernelMaxFiles())

	fmt.Print("KernelAllocateFiles:")
	fmt.Println(KernelAllocateFiles())

	fmt.Print("KernelMaxProc:")
	fmt.Println(KernelMaxProc())

	fmt.Print("KernelHostname:")
	fmt.Println(KernelHostname())

	// loadavg
	fmt.Print("LoadAvg:")
	fmt.Println(LoadAvg())

	// cpuinfo
	fmt.Print("NumCpu:")
	fmt.Println(NumCpu())

	fmt.Print("CpuMHz:")
	fmt.Println(CpuMHz())

	// cpustat
	fmt.Print("CurrentProcStat:")
	fmt.Println(CurrentProcStat())

}
