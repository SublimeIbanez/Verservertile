#ifndef CLIENT_H
#define CLIENT_H

#include "standardized.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdbool.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

typedef struct {
    
} Client;

void error(const char *msg);


#endif