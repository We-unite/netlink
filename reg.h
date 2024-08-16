#ifndef _SEND_H
#define _SEND_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <sys/select.h>
#include <signal.h>
#include <linux/netlink.h>
#include <linux/connector.h>
#include <linux/cn_proc.h>
#include <unistd.h>
#include <errno.h>

typedef struct __attribute__((aligned(NLMSG_ALIGNTO)))
{
    struct nlmsghdr nl_hdr;
    struct __attribute__((__packed__))
    {
        struct cn_msg cn_msg;
        enum proc_cn_mcast_op cn_mcast;
    };
} register_msg_t;

int registerNetlink(int socketfd);

#endif