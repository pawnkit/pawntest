#if defined _pt_sc_gangzone_inc
    #endinput
#endif
#define _pt_sc_gangzone_inc

#define TEST_CREATE_GANGZONE(%0,%1,%2,%3) \
    __pt_gangzone_create((%0), (%1), (%2), (%3))

#define ASSERT_GANGZONE_VALID(%0) \
    if (!__pt_gangzone_valid((%0), __file, __line)) return 0

#define ASSERT_GANGZONE_VISIBLE(%0,%1) \
    if (!__pt_gangzone_visible((%0), (%1), true, __file, __line)) return 0

#define ASSERT_GANGZONE_HIDDEN(%0,%1) \
    if (!__pt_gangzone_visible((%0), (%1), false, __file, __line)) return 0

#define ASSERT_PLAYER_IN_GANGZONE(%0,%1) \
    if (!__pt_gangzone_inside((%0), (%1), true, __file, __line)) return 0

#define ASSERT_PLAYER_NOT_IN_GANGZONE(%0,%1) \
    if (!__pt_gangzone_inside((%0), (%1), false, __file, __line)) return 0

#define TEST_CREATE_PLAYER_GANGZONE(%0,%1,%2,%3,%4) \
    __pt_player_gangzone_create((%0), (%1), (%2), (%3), (%4))

#define ASSERT_PLAYER_GANGZONE_VALID(%0,%1) \
    if (!__pt_player_gangzone_valid((%0), (%1), __file, __line)) return 0

#define ASSERT_PLAYER_GANGZONE_VISIBLE(%0,%1) \
    if (!__pt_player_gangzone_visible((%0), (%1), true, __file, __line)) return 0

#define ASSERT_PLAYER_GANGZONE_HIDDEN(%0,%1) \
    if (!__pt_player_gangzone_visible((%0), (%1), false, __file, __line)) return 0
