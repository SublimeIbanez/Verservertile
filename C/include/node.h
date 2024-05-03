#ifndef NODE_H
#define NODE_H

#include "standardized.h"

#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <string.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

#define MAX_IN_QUEUE 5

typedef enum
{
    SERVER,
    CLIENT,
} Mode;

typedef enum
{
    TCP = SOCK_STREAM,
    UDP = SOCK_DGRAM,
    RUDP = SOCK_SEQPACKET,
} Protocol;

typedef struct
{
    int socket_fd;
    struct sockaddr_storage address;
    socklen_t socklen;
} Node;

Node *create_node(Mode mode, Protocol protocol, const char *port, const char *ip);
void listen_node(Node *node);
void *listener(void *arg);

void free_node(Node *node);

#endif
