#include "info_windows.h"
#include <stdio.h>
#include <Windows.h>

void GetVolumeInfo(char *dir)
{
    BOOL ret;
    __int64 freeBytesAvailableToCaller;
    __int64 totalNumberOfBytes;
    __int64 totalNumberOfFreeBytes;

    ret = GetDiskFreeSpaceExA(dir, (PULARGE_INTEGER)&freeBytesAvailableToCaller, (PULARGE_INTEGER)&totalNumberOfBytes, (PULARGE_INTEGER)&totalNumberOfFreeBytes);
    if (ret == TRUE)
    {
        goCallbackVolumeInformation(totalNumberOfFreeBytes, totalNumberOfBytes);
    }
    else
    {
        printf("[CGO] [ERROR] GetVolumeInformation() failed, error % u\n", GetLastError());
    }
}