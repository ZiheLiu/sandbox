#include <stdlib.h>
#include <stdio.h>

int main() {
    printf("%d", system("shutdown -h now"));
    return 0;
}