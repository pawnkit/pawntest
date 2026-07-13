#if defined _pt_sc_server_inc
    #endinput
#endif
#define _pt_sc_server_inc

#define ASSERT_SERVER_WEATHER(%0) \
    if (!__pt_server_weather((%0), __file, __line)) return 0

#define ASSERT_SERVER_TIME(%0) \
    if (!__pt_server_time((%0), __file, __line)) return 0

#define ASSERT_SERVER_GRAVITY(%0,%1) \
    if (!__pt_server_gravity((%0), (%1), __file, __line)) return 0

#define ASSERT_GAME_MODE_TEXT(%0) \
    if (!__pt_server_mode_text((%0), __file, __line)) return 0

#define ASSERT_SERVER_RULE(%0,%1) \
    if (!__pt_server_rule((%0), (%1), __file, __line)) return 0
