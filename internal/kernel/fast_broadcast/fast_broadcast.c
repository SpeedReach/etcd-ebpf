//go:build ignore

#include "../vmlinux.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#define TC_ACT_OK   0
#define TC_ACT_SHOT 2

#define ETH_P_IP    0x0800

#define IP_P_TCP    6
#define IP_P_UDP    17

#define ETH_SIZE    sizeof(struct ethhdr)
#define IP_SIZE	    sizeof(struct iphdr)
#define UDP_SIZE    sizeof(struct udphdr)
#define TCP_SIZE    sizeof(struct tcphdr)

#define MAX_ENTRIES_PER_PACKET 20

char __license[] SEC("license") = "Dual MIT/GPL";



// Map to store the current term
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY); 
    __type(key, __u32);
    __type(value, __u32);
    __uint(max_entries, 1);
} term SEC(".maps");



SEC("xdp")
int fast_return(struct xdp_md *ctx) {
	void* data = (void*)(long)ctx->data;
	void* data_end = (void*)(long)ctx->data_end;

	struct hdr header = try_parse_udp(data, data_end);
	if (header.udp == NULL) {
		return XDP_PASS;
	}
	if (header.udp->dest != bpf_ntohs(8080)){
		return XDP_PASS;
	}

	u32 term_key = 0;
	u32* term_value = bpf_map_lookup_elem(&term, &term_key);
	if (term_value == NULL) {
		return XDP_PASS;
	}
	bpf_printk("term_value: %d\n", *term_value);

	if (data + ETH_SIZE  + IP_SIZE + UDP_SIZE + sizeof(struct append_entries_reply) + sizeof(struct append_entries_args) > data_end) {
		bpf_printk("not enough data w\n");
		return XDP_PASS;
	}
	

	struct append_entries_args* args = data + ETH_SIZE + IP_SIZE + UDP_SIZE + sizeof(struct append_entries_reply);
	bpf_printk("term: %d\n", args->term);
	bpf_printk("leader_id: %d\n", args->leader_id);
	bpf_printk("prev_log_index: %d\n", args->prev_log_index);
	bpf_printk("prev_log_term: %d\n", args->prev_log_term);
	bpf_printk("leader_commit: %d\n", args->leader_commit);
	bpf_printk("entry_count: %d\n", args->entry_count);

	void* cursor = (void*) args + sizeof(struct append_entries_args);

	void* ringbuf = bpf_ringbuf_reserve(&new_entries, sizeof(struct append_entries_args) + MAX_ENTRIES_PER_PACKET * sizeof(struct log_entry), 0);
	if (!ringbuf) {
		bpf_printk("ringbuf_reserve failed\n");
		return XDP_PASS;
	}

	struct append_entries_args* arg_ring_buf = ringbuf;
	*arg_ring_buf = *args;

	for (int i=0;i<MAX_ENTRIES_PER_PACKET;i++){
		if (i >= args->entry_count) {
			break;
		}
		if (cursor + sizeof(struct log_entry) > data_end) {
			bpf_ringbuf_discard(ringbuf, 0);
			bpf_printk("not enough data\n");
			return XDP_PASS;
		}
		struct log_entry* entry = cursor;
		cursor += sizeof(struct log_entry);
		bpf_printk("parsed entry %d: %d\n", i, entry->term);
		struct log_entry* log_ringbuf_ptr = ringbuf + sizeof(struct append_entries_args) + i * sizeof(struct log_entry);
		*log_ringbuf_ptr = *entry;
	}
	
	bpf_ringbuf_submit(ringbuf, 0);

	struct append_entries_reply reply = {
		.PeerId = 1,
		.Term = *term_value,
		.Success = 1
	};

	int ret = bpf_xdp_store_bytes(ctx, ETH_SIZE + IP_SIZE + UDP_SIZE, &reply, sizeof(reply));
	bpf_printk("ret: %d\n", ret);

	header.udp->dest = header.udp->source;
	header.udp->source = bpf_htons(8080);

	return XDP_TX;
}





static __always_inline struct hdr try_parse_udp(void* data, void* data_end){
	if(data + ETH_SIZE > data_end)
		return (struct hdr) {NULL,NULL, NULL};
	
	struct ethhdr* eth = data;
	if(bpf_ntohs(eth->h_proto) != ETH_P_IP)
		return (struct hdr) {NULL,NULL, NULL};

	if(data + ETH_SIZE + IP_SIZE > data_end)
		return (struct hdr) {NULL,NULL, NULL};
	
	struct iphdr* ip = data + ETH_SIZE;
	if(ip->protocol != IP_P_UDP)
		return (struct hdr) {NULL,NULL, NULL};
	
	if(data + ETH_SIZE + IP_SIZE + UDP_SIZE > data_end)
		return (struct hdr) {NULL,NULL, NULL};
	
	struct udphdr* udp = data + ETH_SIZE + IP_SIZE;

	return (struct hdr){eth,ip, udp};
}


