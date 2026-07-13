#if defined _pt_sc_textlabel_inc
    #endinput
#endif
#define _pt_sc_textlabel_inc

#define TEST_CREATE_TEXT_LABEL(%0,%1,%2,%3,%4,%5,%6) \
    __pt_text_label_create((%0), (%1), (%2), (%3), (%4), (%5), (%6))

#define ASSERT_TEXT_LABEL_VALID(%0) \
    if (!__pt_text_label_valid((%0), __file, __line)) return 0

#define ASSERT_TEXT_LABEL_TEXT(%0,%1) \
    if (!__pt_text_label_text((%0), (%1), __file, __line)) return 0

#define TEST_CREATE_PLAYER_TEXT_LABEL(%0,%1,%2,%3,%4,%5,%6) \
    __pt_player_text_label_create((%0), (%1), (%2), (%3), (%4), (%5), (%6))

#define ASSERT_PLAYER_TEXT_LABEL_VALID(%0,%1) \
    if (!__pt_player_text_label_valid((%0), (%1), __file, __line)) return 0

#define ASSERT_PLAYER_TEXT_LABEL_TEXT(%0,%1,%2) \
    if (!__pt_player_text_label_text((%0), (%1), (%2), __file, __line)) return 0
