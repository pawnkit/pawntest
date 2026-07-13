#if defined _pt_sc_textdraw_inc
    #endinput
#endif
#define _pt_sc_textdraw_inc

#define TEST_CREATE_TEXTDRAW(%0,%1,%2) \
    __pt_textdraw_create((%0), (%1), (%2))

#define ASSERT_TEXTDRAW_VALID(%0) \
    if (!__pt_textdraw_valid((%0), __file, __line)) return 0

#define ASSERT_TEXTDRAW_TEXT(%0,%1) \
    if (!__pt_textdraw_text((%0), (%1), __file, __line)) return 0

#define ASSERT_TEXTDRAW_VISIBLE(%0,%1) \
    if (!__pt_textdraw_visible((%0), (%1), true, __file, __line)) return 0

#define ASSERT_TEXTDRAW_HIDDEN(%0,%1) \
    if (!__pt_textdraw_visible((%0), (%1), false, __file, __line)) return 0

#define TEST_CREATE_PLAYER_TEXTDRAW(%0,%1,%2,%3) \
    __pt_player_textdraw_create((%0), (%1), (%2), (%3))

#define ASSERT_PLAYER_TEXTDRAW_VALID(%0,%1) \
    if (!__pt_player_textdraw_valid((%0), (%1), __file, __line)) return 0

#define ASSERT_PLAYER_TEXTDRAW_TEXT(%0,%1,%2) \
    if (!__pt_player_textdraw_text((%0), (%1), (%2), __file, __line)) return 0

#define ASSERT_PLAYER_TEXTDRAW_VISIBLE(%0,%1) \
    if (!__pt_player_textdraw_visible((%0), (%1), true, __file, __line)) return 0

#define ASSERT_PLAYER_TEXTDRAW_HIDDEN(%0,%1) \
    if (!__pt_player_textdraw_visible((%0), (%1), false, __file, __line)) return 0
