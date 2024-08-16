#include "reg.h"

int registerNetlink(int socketfd) {
	register_msg_t nlcn_msg;
    memset(&nlcn_msg, 0, sizeof(nlcn_msg));
	nlcn_msg.nl_hdr.nlmsg_len = sizeof(nlcn_msg);
	nlcn_msg.nl_hdr.nlmsg_pid = getpid();
	nlcn_msg.nl_hdr.nlmsg_type = NLMSG_DONE;

	nlcn_msg.cn_msg.id.idx = CN_IDX_PROC;
	nlcn_msg.cn_msg.id.val = CN_VAL_PROC;
	nlcn_msg.cn_msg.len = sizeof(enum proc_cn_mcast_op);

	nlcn_msg.cn_mcast = PROC_CN_MCAST_LISTEN;
    if (send(socketfd, &nlcn_msg, sizeof(nlcn_msg), 0) == -1)
    {
        close(socketfd);
        return -1;
    }
    return 0;
}