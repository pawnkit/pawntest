#if defined _pt_sc_player_inc
    #endinput
#endif
#define _pt_sc_player_inc

#define TEST_CREATE_PLAYER(%0) \
    __pt_player_create((%0))

#define ASSERT_PLAYER_CONNECTED(%0) \
    if (!__pt_player_connected((%0), true, __file, __line)) return 0

#define ASSERT_PLAYER_DISCONNECTED(%0) \
    if (!__pt_player_connected((%0), false, __file, __line)) return 0

#define ASSERT_PLAYER_MONEY(%0,%1) \
    if (!__pt_player_money((%0), (%1), __file, __line)) return 0

#define ASSERT_PLAYER_MESSAGE(%0,%1) \
    if (!__pt_player_message((%0), (%1), false, __file, __line)) return 0

#define ASSERT_PLAYER_MESSAGE_CONTAINS(%0,%1) \
    if (!__pt_player_message((%0), (%1), true, __file, __line)) return 0

#define ASSERT_PLAYER_POS_NEAR(%0,%1,%2,%3,%4) \
    if (!__pt_player_pos_near((%0), (%1), (%2), (%3), (%4), __file, __line)) return 0
