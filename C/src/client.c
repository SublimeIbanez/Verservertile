#include "client.h"

#include <stdlib.h>

void error(const char *msg)
{
    perror(msg);
    exit(EXIT_FAILURE);
}