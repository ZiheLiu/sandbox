#include <stdio.h>
#include <sys/signal.h>

int main() {
    printf("%d", kill(1, SIGSEGV));
    return 0;
}