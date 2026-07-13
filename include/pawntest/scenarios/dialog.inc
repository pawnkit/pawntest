#if defined _pt_sc_dialog_inc
    #endinput
#endif
#define _pt_sc_dialog_inc

#define TEST_RESPOND_DIALOG(%0,%1,%2,%3) \
    __pt_dialog_respond((%0), (%1), (%2), (%3))

#define ASSERT_DIALOG_VISIBLE(%0,%1) \
    if (!__pt_dialog_visible((%0), (%1), true, __file, __line)) return 0

#define ASSERT_DIALOG_HIDDEN(%0) \
    if (!__pt_dialog_visible((%0), -1, false, __file, __line)) return 0

#define ASSERT_DIALOG_TITLE(%0,%1) \
    if (!__pt_dialog_title((%0), (%1), __file, __line)) return 0

#define ASSERT_DIALOG_BODY(%0,%1) \
    if (!__pt_dialog_body((%0), (%1), __file, __line)) return 0
