package nux

import (
	"fmt"
	"testing"
)

func TestMetrics(t *testing.T) {

	fmt.Println("======kernel======")
	fmt.Print("KernelMaxFiles:")
	fmt.Println(KernelMaxFiles())

	fmt.Print("KernelAllocateFiles:")
	fmt.Println(KernelAllocateFiles())

	fmt.Print("KernelMaxProc:")
	fmt.Println(KernelMaxProc())

	fmt.Print("KernelHostname:")
	fmt.Println(KernelHostname())

	fmt.Println("======loadavg======")
	fmt.Print("LoadAvg:")
	fmt.Println(LoadAvg())

	fmt.Println("======cpuinfo======")
	fmt.Print("NumCpu:")
	fmt.Println(NumCpu())

	fmt.Print("CpuMHz:")
	fmt.Println(CpuMHz())

	fmt.Println("======cpustat======")
	fmt.Print("CurrentProcStat:")
	fmt.Println(CurrentProcStat())

	fmt.Println("======dfstat======")
	fmt.Print("ListMountPoint:")
	mountPoints, err := ListMountPoint()
	fmt.Println(mountPoints, err)
	fmt.Print("DeviceUsage:")
	for _, arr := range mountPoints {
		fmt.Println(BuildDeviceUsage(arr[0], arr[1], arr[2]))
	}

}
