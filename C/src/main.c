#include "node.h"

#include <unistd.h>
#include <getopt.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <strings.h>

// Forward declarations
void process_args(int argc, char **argv);

static struct option long_options[] = {
    {"mode", required_argument, 0, 'm'},
    {"host", required_argument, 0, 'h'},
    {0, 0, 0, 0},
};

// port = "8080";
// host = "172.18.0.3";
char *port;
char *host;
Mode mode = SERVER;

int main(int argc, char **argv)
{
    if (argc < 2)
    {
        perror("Invalid arguments passed\n  Must pass -h ipaddress:port");
        exit(EXIT_FAILURE);
    }

    process_args(argc, argv);

    Node *node = create_node(mode, TCP, port, host);
    if (!node)
    {
        perror("Could not generate node!");
        exit(EXIT_FAILURE);
    }

    if (mode == SERVER)
    {
        listen_node(node);
    }

    free(host);
    free(port);
    free_node(node);
    return EXIT_SUCCESS;
}

void process_args(int argc, char **argv)
{
    int opt;
    bool hflag = false;
    bool mflag = false;

    while ((opt = getopt_long(argc, argv, "h:m:", long_options, NULL)) != -1)
    {
        switch (opt)
        {
        case 'h':
        {
            char *input = optarg;
            char *colon = strchr(input, ':');
            if (colon == NULL)
            {
                fprintf(stderr, "Error: Incorrect format for host and port:\n  Use: hostname:port\n");
                exit(EXIT_FAILURE);
            }

            int host_length = colon - input;
            host = malloc(host_length + 1);
            if (host == NULL)
            {
                perror("Failed to allocate memory for the host");
                exit(EXIT_FAILURE);
            }
            strncpy(host, input, host_length);
            host[host_length] = '\0';

            port = strdup(colon + 1);
            if (port == NULL)
            {
                perror("Failed to allocate memory for port");
                exit(EXIT_FAILURE);
            }
            hflag = true;

            break;
        }

        case 'm':
        {
            if (strncasecmp(optarg, "client", 6) == 0)
            {
                mode = CLIENT;
            }
            if (strncasecmp(optarg, "server", 6) == 0)
            {
                mode = SERVER;
            }
            break;
        }

        default:
        {
            fprintf(stderr, "Usage: %s -h hostname:port\n", argv[0]);
            exit(EXIT_FAILURE);
        }
        }
    }

    if (!hflag)
    {
        fprintf(stderr, "Usage: %s -h hostname:port\n", argv[0]);
        exit(EXIT_FAILURE);
    }

    printf("%s connecting on Host: %s, Port: %s\n", mode == SERVER ? "Server" : "Client", host, port);
}