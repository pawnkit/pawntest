#if defined _pt_sc_object_inc
    #endinput
#endif
#define _pt_sc_object_inc

#define TEST_CREATE_OBJECT(%0,%1,%2,%3) \
    __pt_object_create((%0), (%1), (%2), (%3))

#define ASSERT_OBJECT_VALID(%0) \
    if (!__pt_object_valid((%0), __file, __line)) return 0

#define ASSERT_OBJECT_MODEL(%0,%1) \
    if (!__pt_object_model((%0), (%1), __file, __line)) return 0

#define ASSERT_OBJECT_POS_NEAR(%0,%1,%2,%3,%4) \
    if (!__pt_object_pos_near((%0), (%1), (%2), (%3), (%4), __file, __line)) return 0

#define TEST_CREATE_PLAYER_OBJECT(%0,%1,%2,%3,%4) \
    __pt_player_object_create((%0), (%1), (%2), (%3), (%4))

#define ASSERT_PLAYER_OBJECT_VALID(%0,%1) \
    if (!__pt_player_object_valid((%0), (%1), __file, __line)) return 0
