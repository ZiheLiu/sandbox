#include <stdlib.h>
#include <iostream>

using namespace std;

int main() {
    cout << system("touch /test.txt") << " ";
    cout << system("reboot");
    return 0;
}