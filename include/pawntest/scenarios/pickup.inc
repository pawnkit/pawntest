#if defined _pt_sc_pickup_inc
    #endinput
#endif
#define _pt_sc_pickup_inc

#define TEST_CREATE_PICKUP(%0,%1,%2,%3,%4,%5) \
    __pt_pickup_create((%0), (%1), (%2), (%3), (%4), (%5))

#define ASSERT_PICKUP_VALID(%0) \
    if (!__pt_pickup_valid((%0), __file, __line)) return 0

#define ASSERT_PICKUP_MODEL(%0,%1) \
    if (!__pt_pickup_model((%0), (%1), __file, __line)) return 0

#define ASSERT_PICKUP_POS_NEAR(%0,%1,%2,%3,%4) \
    if (!__pt_pickup_pos_near((%0), (%1), (%2), (%3), (%4), __file, __line)) return 0

#define TEST_CREATE_PLAYER_PICKUP(%0,%1,%2,%3,%4,%5,%6) \
    __pt_player_pickup_create((%0), (%1), (%2), (%3), (%4), (%5), (%6))

#define ASSERT_PLAYER_PICKUP_VALID(%0,%1) \
    if (!__pt_player_pickup_valid((%0), (%1), __file, __line)) return 0
