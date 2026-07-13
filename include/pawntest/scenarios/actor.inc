#if defined _pt_sc_actor_inc
    #endinput
#endif
#define _pt_sc_actor_inc

#define TEST_CREATE_ACTOR(%0,%1,%2,%3,%4) \
    __pt_actor_create((%0), (%1), (%2), (%3), (%4))

#define ASSERT_ACTOR_VALID(%0) \
    if (!__pt_actor_valid((%0), __file, __line)) return 0

#define ASSERT_ACTOR_SKIN(%0,%1) \
    if (!__pt_actor_skin((%0), (%1), __file, __line)) return 0

#define ASSERT_ACTOR_HEALTH(%0,%1) \
    if (!__pt_actor_health((%0), (%1), __file, __line)) return 0

#define ASSERT_ACTOR_POS_NEAR(%0,%1,%2,%3,%4) \
    if (!__pt_actor_pos_near((%0), (%1), (%2), (%3), (%4), __file, __line)) return 0
