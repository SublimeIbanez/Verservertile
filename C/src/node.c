#include "node.h"
#include <stdlib.h>
#include <stdbool.h>
#include <signal.h>

volatile sig_atomic_t run_server = true;

void sig_handler(int sig)
{
    run_server = false;
}

Node *create_node(Mode mode, Protocol protocol, const char *port, const char *ip)
{
    Node *node = malloc(sizeof(Node));
    if (!node)
    {
        return NULL;
    }

    struct sigaction sa;
    sa.sa_handler = sig_handler;
    sigemptyset(&sa.sa_mask);
    sigaction(SIGINT, &sa, NULL);
    sigaction(SIGTERM, &sa, NULL);

    // Hints tell the type of information getaddrinfo() will provide
    struct addrinfo hints, *address, *address_pointer;
    memset(&hints, 0, sizeof(struct addrinfo));
    hints.ai_family = AF_UNSPEC;
    hints.ai_socktype = protocol;
    hints.ai_flags = (mode == SERVER ? AI_PASSIVE : 0);

    // Try to get addressinfo
    int s = getaddrinfo(ip, port, &hints, &address);
    if (s != 0)
    {
        free(node);
        fprintf(stderr, "getaddrinfo %s\n", gai_strerror(s));
        exit(EXIT_FAILURE);
    }

    // Go through each address struct and try each address until one is successfully bound
    bool bound = false;
    for (address_pointer = address; address_pointer != NULL; address_pointer = address_pointer->ai_next)
    {
        node->socket_fd = socket(address_pointer->ai_family, address_pointer->ai_socktype, address_pointer->ai_protocol);
        if (node->socket_fd == -1)
        {
            continue;
        }

        int result = (mode == SERVER ? bind(node->socket_fd, address_pointer->ai_addr, address_pointer->ai_addrlen)
                                     : connect(node->socket_fd, address_pointer->ai_addr, address_pointer->ai_addrlen));
        if (result == 0)
        {
            memcpy(&node->address, address_pointer->ai_addr, address_pointer->ai_addrlen);
            node->socklen = address_pointer->ai_addrlen;
            bound = true;
            freeaddrinfo(address);
            return node;
        }

        close(node->socket_fd);
    }

    freeaddrinfo(address);
    free_node(node);
    return NULL;
}

void listen_node(Node *node)
{
    if (listen(node->socket_fd, MAX_IN_QUEUE) == -1)
    {
        perror("Could not listen");
        close(node->socket_fd);
        free_node(node);
        exit(EXIT_FAILURE);
    }

    printf("Listening on socket: %d\n", node->socket_fd);

    while (run_server)
    {
        struct sockaddr_storage client_addr;
        socklen_t client_addr_size = sizeof(client_addr);
        int *client_fd = malloc(sizeof(int));
        if (client_fd == NULL)
        {
            perror("Failed to allocate memory for client_fd");
            continue;
        }

        *client_fd = accept(node->socket_fd, (struct sockaddr *)&client_addr, &client_addr_size);
        if (*client_fd == -1)
        {
            perror("Accept failed");
            close(node->socket_fd);
            free(client_fd);
            continue;
        }

        printf("Accepted connection on socket %d\n", *client_fd);

        pthread_t thread;
        if (pthread_create(&thread, NULL, listener, client_fd) == -1)
        {
            perror("Failed to create thread");
            close(*client_fd);
            free(client_fd);
            continue;
        }
        pthread_detach(thread);
    }
}

void *listener(void *arg)
{
    int client_fd = *((int *)arg);
    free(arg);

    char buffer[BUFFER_SIZE];
    ssize_t nread = recv(client_fd, buffer, sizeof(buffer), 0);
    if (nread != -1)
    {
        printf("Received message: %s\n", buffer);
        send(client_fd, buffer, nread, 0); // Testing - send back received message
    }

    close(client_fd);
    return NULL;
}

void free_node(Node *node)
{
    if (node != NULL)
    {
        free(node);
        node = NULL;
    }
}