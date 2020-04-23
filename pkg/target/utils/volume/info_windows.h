#ifndef WIN_H_
#define WIN_H_

#include <stdlib.h>
#include <stdio.h>
#include <Windows.h>

void GetVolumeInfo(char *dir);

// this is a call-back function from the go code
extern void goCallbackVolumeInformation(__int64 free, __int64 total);

#endif // WIN_H_