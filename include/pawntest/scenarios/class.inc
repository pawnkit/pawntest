#if defined _pt_sc_class_inc
    #endinput
#endif
#define _pt_sc_class_inc


#define TEST_SELECT_CLASS(%0,%1) \
    __pt_class_select((%0), (%1))

#define ASSERT_CLASS_COUNT(%0) \
    if (!__pt_class_count((%0), __file, __line)) return 0

#define ASSERT_PLAYER_CLASS(%0,%1) \
    if (!__pt_player_class((%0), (%1), __file, __line)) return 0

#define ASSERT_PLAYER_SELECTING_CLASS(%0) \
    if (!__pt_player_selecting_class((%0), true, __file, __line)) return 0

#define ASSERT_PLAYER_NOT_SELECTING_CLASS(%0) \
    if (!__pt_player_selecting_class((%0), false, __file, __line)) return 0
