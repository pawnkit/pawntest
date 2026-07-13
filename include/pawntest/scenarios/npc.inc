#if defined _pt_sc_npc_inc
    #endinput
#endif
#define _pt_sc_npc_inc

#define TEST_CREATE_NPC(%0) \
    __pt_npc_create((%0))

#define ASSERT_NPC_VALID(%0) \
    if (!__pt_npc_valid((%0), __file, __line)) return 0

#define ASSERT_NPC_SPAWNED(%0) \
    if (!__pt_npc_spawned((%0), true, __file, __line)) return 0

#define ASSERT_NPC_NOT_SPAWNED(%0) \
    if (!__pt_npc_spawned((%0), false, __file, __line)) return 0

#define ASSERT_NPC_HEALTH(%0,%1,%2) \
    if (!__pt_npc_health((%0), (%1), (%2), __file, __line)) return 0

#define ASSERT_NPC_POS_NEAR(%0,%1,%2,%3,%4) \
    if (!__pt_npc_pos_near((%0), (%1), (%2), (%3), (%4), __file, __line)) return 0

#define ASSERT_NPC_WEAPON(%0,%1) \
    if (!__pt_npc_weapon((%0), (%1), __file, __line)) return 0

#define ASSERT_NPC_AIMING(%0) \
    if (!__pt_npc_aiming((%0), true, __file, __line)) return 0

#define ASSERT_NPC_NOT_AIMING(%0) \
    if (!__pt_npc_aiming((%0), false, __file, __line)) return 0

#define ASSERT_NPC_ANIMATION(%0,%1,%2) \
    if (!__pt_npc_animation((%0), (%1), (%2), __file, __line)) return 0

#define ASSERT_NPC_PLAYBACK(%0) \
    if (!__pt_npc_playback((%0), true, __file, __line)) return 0

#define ASSERT_NPC_PLAYBACK_STOPPED(%0) \
    if (!__pt_npc_playback((%0), false, __file, __line)) return 0

#define ASSERT_NPC_PATH_COUNT(%0) \
    if (!__pt_npc_path_count((%0), __file, __line)) return 0

#define ASSERT_NPC_RECORD_COUNT(%0) \
    if (!__pt_npc_record_count((%0), __file, __line)) return 0
