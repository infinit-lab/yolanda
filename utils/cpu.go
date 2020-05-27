package utils

/*
#include <cpuid.h>
#include <stdio.h>

void getCpuId() {
	char id[256] = {0};
	unsigned int cpuId[4] = {0};

	__cpuid(0, cpuId[0], cpuId[1], cpuId[2], cpuId[3]);
	sprintf(id, "%08x %08x %08x %08x", cpuId[0], cpuId[1], cpuId[2], cpuId[3]);
	printf("%s\n", id);
}
*/
import "C"

func GetCpuID() (string, error) {
	C.getCpuId()
	return "", nil
}
