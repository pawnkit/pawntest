#if defined _pt_sc_menu_inc
    #endinput
#endif
#define _pt_sc_menu_inc

#define TEST_SELECT_MENU_ROW(%0,%1) \
    __pt_menu_select((%0), (%1))

#define TEST_EXIT_MENU(%0) \
    __pt_menu_exit((%0))

#define ASSERT_MENU_VALID(%0) \
    if (!__pt_menu_valid((%0), __file, __line)) return 0

#define ASSERT_MENU_VISIBLE(%0,%1) \
    if (!__pt_menu_visible((%0), (%1), true, __file, __line)) return 0

#define ASSERT_MENU_HIDDEN(%0,%1) \
    if (!__pt_menu_visible((%0), (%1), false, __file, __line)) return 0

#define ASSERT_MENU_ITEMS(%0,%1,%2) \
    if (!__pt_menu_items((%0), (%1), (%2), __file, __line)) return 0
