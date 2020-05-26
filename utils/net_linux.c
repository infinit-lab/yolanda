#ifdef LINUX

#include "net.h"
#include <stdio.h>
#include <string.h>
#include <net/if.h>
#include <sys/ioctl.h>
#include <arpa/inet.h>
#include <errno.h>
#include <unistd.h>

int GetAdaptersCount() {
	int fd = 0;
	int interfaceNum = 0;
	struct ifreq buf[16] = {0};
	struct ifconf ifc = {0};

	if ((fd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
		close(fd);
		return 0;
	}

	ifc.ifc_len = sizeof(buf);
	ifc.ifc_buf = (caddr_t)buf;
	if (!ioctl(fd, SIOCGIFCONF, (char *)&ifc)) {
		interfaceNum = ifc.ifc_len / sizeof(struct ifreq);
	}
	close(fd);
	return interfaceNum;
}

int GetAdapter(int i, PT_ADAPTER pAdapter) {
	if (!pAdapter) {
		return NET_ERROR_SUCCESS;
	}
	int fd = 0;
	int interfaceNum = 0;
	struct ifreq buf[16] = {0};
	struct ifconf ifc = {0};
	char ip[32] = {0};
	char broadAddr[32] = {0};
	char subnetMask[32] = {0};

	if ((fd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
		goto failed;
	}

	ifc.ifc_len = sizeof(buf);
	ifc.ifc_buf = (caddr_t)buf;
	if (ioctl(fd, SIOCGIFCONF, (char *)&ifc) != 0) {
		goto failed;
	}
	interfaceNum = ifc.ifc_len / sizeof(struct ifreq);
	if (i >= interfaceNum) {
		goto failed;
	}
	strncpy(pAdapter->name, buf[i].ifr_name, NAME_LENGTH - 1);

	if (ioctl(fd, SIOCGIFHWADDR, (char*)(&buf[i])) != 0) {
		goto failed;
	}
	memset(pAdapter->mac, 0, sizeof(pAdapter->mac));
	sprintf(pAdapter->mac, "%02x-%02x-%02x-%02x-%02x-%02x",
		(unsigned char)buf[i].ifr_hwaddr.sa_data[0],
		(unsigned char)buf[i].ifr_hwaddr.sa_data[1],
		(unsigned char)buf[i].ifr_hwaddr.sa_data[2],
		(unsigned char)buf[i].ifr_hwaddr.sa_data[3],
		(unsigned char)buf[i].ifr_hwaddr.sa_data[4],
		(unsigned char)buf[i].ifr_hwaddr.sa_data[5]);

	if (ioctl(fd, SIOCGIFADDR, (char*)(&buf[i])) != 0) {
		goto failed;
	}
	sprintf(pAdapter->ip, "%s", (char*)inet_ntoa(((struct sockaddr_in*)&(buf[i].ifr_addr))->sin_addr));

	if (ioctl(fd, SIOCGIFNETMASK, (char*)(&buf[i])) != 0) {
		goto failed;
	}
	sprintf(pAdapter->mask, "%s", (char*)inet_ntoa(((struct sockaddr_in*)&(buf[i].ifr_netmask))->sin_addr));

	close(fd);
	return NET_ERROR_SUCCESS;
failed:
	close(fd);
	return NET_ERROR_FAILED;
}

#endif

