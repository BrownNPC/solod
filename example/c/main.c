#include <netdb.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>

int main(void) {
    struct addrinfo hints = {
        .ai_family = AF_UNSPEC,
        .ai_socktype = SOCK_STREAM,
        .ai_flags = AI_PASSIVE,
    };
    struct addrinfo* res = NULL;
    int err = getaddrinfo(NULL, "8080", &hints, &res);
    if (err) {
        fprintf(stderr, "err: %s\n", gai_strerror(err));
        return 1;
    }
    printf("ok\n");
    freeaddrinfo(res);
    return 0;
}
