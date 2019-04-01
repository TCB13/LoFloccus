#include <windows.h>
#include <process.h>

// windres LoFloccus.rc -O coff -o LoFloccus.res
// gcc -Wall -std=c99 -o LoFloccus.exe LoFloccus.c LoFloccus.res
int main()
{
    ShowWindow(GetConsoleWindow(), SW_HIDE);
    spawnl(P_OVERLAY , "venv\\Scripts\\python.exe", "venv\\Scripts\\python.exe", "app\\LoFloccus.py", NULL );
    return 0;
}
