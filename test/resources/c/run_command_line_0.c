#include <stdlib.h>
#include <stdio.h>

int main() {
    printf("%d", system("rm /etc/hosts"));
    return 0;
}