#ifndef RUST_C_CHANNEL_MANAGER
#define RUST_C_CHANNEL_MANAGER

#include <stdint.h>

// EXPORTED STRUCT
typedef struct RootChannel root_channel_t;
typedef struct DailyChannelManager daily_channel_t;

typedef struct ChannelInfo{
    char *channel_id;
    char *announce_id;
} channel_info_t;

typedef struct KeyNonce{
    uint8_t key[32];
    uint8_t nonce[24];
} key_nonce_t;

typedef struct RawPacket raw_packet_t;

//EXPORTED FUNCTIONS
extern void hello_from_rust(const char *str);
extern root_channel_t *new_root_channel(int mainnet);
extern root_channel_t *import_root_channel_from_tangle(channel_info_t const *, char const *channel_psw, int mainnet);
extern void drop_root_channel(root_channel_t *);
extern channel_info_t const *open_root_channel(root_channel_t *, char const *channel_psw);
extern daily_channel_t *get_create_daily_actor_channel(root_channel_t *, int category, char const *actor_id,
                                                       char const *state_psw, uint16_t day, uint16_t month, uint16_t year);
extern channel_info_t const *root_channel_info(root_channel_t *);
extern void drop_daily_channel_manager(daily_channel_t *);
extern void print_channel_tree(root_channel_t *);
extern char const *send_raw_packet(daily_channel_t *, raw_packet_t const *, key_nonce_t const *);
extern channel_info_t const *daily_channel_info(daily_channel_t *);
extern channel_info_t const *new_channel_info(char const *channel_id, char const *announce_id);
extern void drop_channel_info(channel_info_t *);
extern key_nonce_t const *new_encryption_key_nonce(char const *key, char const *nonce);
extern void drop_key_nonce(key_nonce_t *);
extern raw_packet_t const *new_raw_packet(uint8_t *pub, uint64_t p_len, uint8_t *mask, uint64_t m_len);
extern void drop_raw_packet(raw_packet_t *);
extern void drop_str(char *);



#endif //RUST_C_CHANNEL_MANAGER
