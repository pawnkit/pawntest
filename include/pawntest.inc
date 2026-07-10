#if defined _pawntest_included
    #endinput
#endif
#define _pawntest_included

#include <pawntest/core>
#include <pawntest/assertions>
#include <pawntest/mocks>
#include <pawntest/scenarios>

forward __pawntest_marker();
public __pawntest_marker()
{
    return 1;
}

#define TEST(%0) \
    forward test_%0(); \
    public test_%0()

#define BEFORE_EACH() \
    forward test_setup(); \
    public test_setup()

#define AFTER_EACH() \
    forward test_teardown(); \
    public test_teardown()

#define BEFORE_ALL() \
    forward test_suite_setup(); \
    public test_suite_setup()

#define AFTER_ALL() \
    forward test_suite_teardown(); \
    public test_suite_teardown()

#define FIXTURE_SETUP(%0) \
    forward fixture_%0_setup(); \
    public fixture_%0_setup()

#define FIXTURE_TEARDOWN(%0) \
    forward fixture_%0_teardown(); \
    public fixture_%0_teardown()

#define USE_FIXTURE(%0) \
    if (!__pawntest_use_fixture("fixture_" #%0 "_setup", "fixture_" #%0 "_teardown", __file, __line)) return 0

#define TEST_CREATE_PLAYER(%0) \
    __pt_player_create((%0))

#define ASSERT_PLAYER_CONNECTED(%0) \
    if (!__pt_player_connected((%0), true, __file, __line)) return 0

#define ASSERT_PLAYER_DISCONNECTED(%0) \
    if (!__pt_player_connected((%0), false, __file, __line)) return 0

#define ASSERT_PLAYER_MONEY(%0,%1) \
    if (!__pt_player_money((%0), (%1), __file, __line)) return 0

#define ASSERT_PLAYER_MESSAGE(%0,%1) \
    if (!__pt_player_message((%0), (%1), false, __file, __line)) return 0

#define ASSERT_PLAYER_MESSAGE_CONTAINS(%0,%1) \
    if (!__pt_player_message((%0), (%1), true, __file, __line)) return 0

#define ASSERT_PLAYER_POS_NEAR(%0,%1,%2,%3,%4) \
    if (!__pt_player_pos_near((%0), (%1), (%2), (%3), (%4), __file, __line)) return 0

#define TEST_CREATE_VEHICLE(%0,%1,%2,%3) \
    __pt_vehicle_create((%0), (%1), (%2), (%3))

#define ASSERT_VEHICLE_VALID(%0) \
    if (!__pt_vehicle_valid((%0), __file, __line)) return 0

#define ASSERT_VEHICLE_MODEL(%0,%1) \
    if (!__pt_vehicle_model((%0), (%1), __file, __line)) return 0

#define ASSERT_VEHICLE_HEALTH(%0,%1) \
    if (!__pt_vehicle_health((%0), (%1), __file, __line)) return 0

#define ASSERT_VEHICLE_POS_NEAR(%0,%1,%2,%3,%4) \
    if (!__pt_vehicle_pos_near((%0), (%1), (%2), (%3), (%4), __file, __line)) return 0

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

#define TEST_CREATE_PICKUP(%0,%1,%2,%3,%4,%5) \
    __pt_pickup_create((%0), (%1), (%2), (%3), (%4), (%5))

#define ASSERT_PICKUP_VALID(%0) \
    if (!__pt_pickup_valid((%0), __file, __line)) return 0

#define ASSERT_PICKUP_MODEL(%0,%1) \
    if (!__pt_pickup_model((%0), (%1), __file, __line)) return 0

#define ASSERT_PICKUP_POS_NEAR(%0,%1,%2,%3,%4) \
    if (!__pt_pickup_pos_near((%0), (%1), (%2), (%3), (%4), __file, __line)) return 0

#define TEST_CREATE_PLAYER_PICKUP(%0,%1,%2,%3,%4,%5,%6) \
    __pt_player_pickup_create((%0), (%1), (%2), (%3), (%4), (%5), (%6))

#define ASSERT_PLAYER_PICKUP_VALID(%0,%1) \
    if (!__pt_player_pickup_valid((%0), (%1), __file, __line)) return 0

#define ASSERT_CHECKPOINT_ACTIVE(%0) \
    if (!__pt_checkpoint_active((%0), true, __file, __line)) return 0

#define ASSERT_CHECKPOINT_INACTIVE(%0) \
    if (!__pt_checkpoint_active((%0), false, __file, __line)) return 0

#define ASSERT_PLAYER_IN_CHECKPOINT(%0) \
    if (!__pt_checkpoint_inside((%0), true, __file, __line)) return 0

#define ASSERT_PLAYER_NOT_IN_CHECKPOINT(%0) \
    if (!__pt_checkpoint_inside((%0), false, __file, __line)) return 0

#define ASSERT_RACE_CHECKPOINT_ACTIVE(%0) \
    if (!__pt_race_checkpoint_active((%0), true, __file, __line)) return 0

#define ASSERT_RACE_CHECKPOINT_INACTIVE(%0) \
    if (!__pt_race_checkpoint_active((%0), false, __file, __line)) return 0

#define ASSERT_PLAYER_IN_RACE_CHECKPOINT(%0) \
    if (!__pt_race_checkpoint_inside((%0), true, __file, __line)) return 0

#define ASSERT_PLAYER_NOT_IN_RACE_CHECKPOINT(%0) \
    if (!__pt_race_checkpoint_inside((%0), false, __file, __line)) return 0

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

#define ASSERT_SVAR_INT(%0,%1) \
    if (!__pt_svar_int((%0), (%1), __file, __line)) return 0

#define ASSERT_SVAR_FLOAT(%0,%1,%2) \
    if (!__pt_svar_float((%0), (%1), (%2), __file, __line)) return 0

#define ASSERT_SVAR_STRING(%0,%1) \
    if (!__pt_svar_string((%0), (%1), __file, __line)) return 0

#define ASSERT_PVAR_INT(%0,%1,%2) \
    if (!__pt_pvar_int((%0), (%1), (%2), __file, __line)) return 0

#define ASSERT_PVAR_FLOAT(%0,%1,%2,%3) \
    if (!__pt_pvar_float((%0), (%1), (%2), (%3), __file, __line)) return 0

#define ASSERT_PVAR_STRING(%0,%1,%2) \
    if (!__pt_pvar_string((%0), (%1), (%2), __file, __line)) return 0

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

#define ASSERT_DATABASE_CONNECTIONS(%0) \
    if (!__pt_database_connections((%0), __file, __line)) return 0

#define ASSERT_DATABASE_RESULTS(%0) \
    if (!__pt_database_results((%0), __file, __line)) return 0

#define ASSERT(%0) \
    if (!__pawntest_assert((%0), #%0, __file, __line)) return 0

#define ASSERT_TRUE(%0) \
    ASSERT(%0)

#define ASSERT_FALSE(%0) \
    if (!__pawntest_assert(!(%0), "!(" #%0 ")", __file, __line)) return 0

#define ASSERT_EQ(%0,%1) \
    if (!__pawntest_assert_eq((%0), (%1), #%0, #%1, __file, __line)) return 0

#define ASSERT_NE(%0,%1) \
    if (!__pawntest_assert_ne((%0), (%1), #%0, #%1, __file, __line)) return 0

#define ASSERT_STR_EQ(%0,%1) \
    if (!__pawntest_assert_str_eq(%0, %1, #%0, #%1, __file, __line)) return 0

#define ASSERT_GT(%0,%1) \
    if (!__pawntest_assert_compare((%0), (%1), 1, #%0, #%1, __file, __line)) return 0

#define ASSERT_GE(%0,%1) \
    if (!__pawntest_assert_compare((%0), (%1), 2, #%0, #%1, __file, __line)) return 0

#define ASSERT_LT(%0,%1) \
    if (!__pawntest_assert_compare((%0), (%1), 3, #%0, #%1, __file, __line)) return 0

#define ASSERT_LE(%0,%1) \
    if (!__pawntest_assert_compare((%0), (%1), 4, #%0, #%1, __file, __line)) return 0

#define ASSERT_STR_CONTAINS(%0,%1) \
    if (!__pawntest_assert_string(%0, %1, 1, #%0, #%1, __file, __line)) return 0

#define ASSERT_STR_PREFIX(%0,%1) \
    if (!__pawntest_assert_string(%0, %1, 2, #%0, #%1, __file, __line)) return 0

#define ASSERT_STR_SUFFIX(%0,%1) \
    if (!__pawntest_assert_string(%0, %1, 3, #%0, #%1, __file, __line)) return 0

#define ASSERT_FLOAT_NEAR(%0,%1,%2) \
    if (!__pawntest_assert_float_near((%0), (%1), (%2), #%0, #%1, __file, __line)) return 0

#define ASSERT_ARRAY_EQ(%0,%1,%2) \
    if (!__pawntest_assert_array_eq(%0, %1, (%2), #%0, #%1, __file, __line)) return 0

#define ASSERT_BETWEEN(%0,%1,%2) \
    if (!__pawntest_assert_between((%0), (%1), (%2), #%0, __file, __line)) return 0

#define ASSERT_HAS_FLAG(%0,%1) \
    if (!__pawntest_assert_has_flag((%0), (%1), #%0, __file, __line)) return 0

#define ASSERT_ARRAY_CONTAINS(%0,%1,%2) \
    if (!__pt_assert_array_contains(%0, (%1), (%2), #%0, __file, __line)) return 0

#define ASSERT_STR_IEQ(%0,%1) \
    if (!__pawntest_assert_str_ieq(%0, %1, #%0, #%1, __file, __line)) return 0

#define ASSERT_EQ_MSG(%0,%1,%2) \
    if (!__pawntest_assert_eq_message((%0), (%1), #%0, #%1, (%2), __file, __line)) return 0

#define EXPECT_ERROR(%0,%1) \
    if (!__pawntest_expect_error(#%0, (%1), __file, __line)) return 0

#define ASSERT_NO_ERROR(%0) \
    if (!__pawntest_expect_no_error(#%0, __file, __line)) return 0

#define XFAIL(%0,%1) \
    return __pawntest_xfail(#%0, (%1), __file, __line)

#define TEST_CASE(%0,%1,%2) \
    forward test_%0(); \
    public test_%0() { return %1(%2); }

#define TEST_CASE2(%0,%1,%2,%3) \
    forward test_%0(); \
    public test_%0() { return %1(%2, %3); }

#define TEST_CASE3(%0,%1,%2,%3,%4) \
    forward test_%0(); \
    public test_%0() { return %1(%2, %3, %4); }

#define TEST_TAG(%0) \
    forward ptt_%0(); \
    public ptt_%0() { return 1; }

#define ASSERT_SNAPSHOT(%0,%1) \
    if (!__pawntest_snapshot((%0), (%1), __file, __line)) return 0

#define FUZZ_INT(%0,%1,%2,%3) \
    return __pawntest_fuzz_int(#%0, (%1), (%2), (%3), __file, __line)

#define CHECK(%0) \
    __pawntest_assert((%0), #%0, __file, __line)

#define CHECK_EQ(%0,%1) \
    __pawntest_assert_eq((%0), (%1), #%0, #%1, __file, __line)

#define CHECK_NE(%0,%1) \
    __pawntest_assert_ne((%0), (%1), #%0, #%1, __file, __line)

#define CHECK_STR_EQ(%0,%1) \
    __pawntest_assert_str_eq(%0, %1, #%0, #%1, __file, __line)

#define FAIL(%0) \
    return __pawntest_fail(%0, __file, __line)

#define SKIP(%0) \
    return __pawntest_skip(%0, __file, __line)

#define MOCK_RETURN(%0,%1) \
    __pawntest_mock_return(#%0, (%1))

#define MOCK_RETURN_ONCE(%0,%1) \
    __pawntest_mock_return_once(#%0, (%1))

#define MOCK_OUTPUT(%0,%1,%2) \
    __pawntest_mock_output(#%0, (%1), (%2))

#define MOCK_OUTPUT_STRING(%0,%1,%2,%3) \
    __pawntest_mock_output_string(#%0, (%1), (%2), (%3))

#define MOCK_CALLBACK(%0,%1) \
    __pawntest_mock_callback(#%0, #%1)

#define MOCK_RESET(%0) \
    __pawntest_mock_reset(#%0)

#define EXPECT_NATIVE_CALLS(%0,%1) \
    __pawntest_expect_calls(#%0, (%1), (%1), __file, __line)

#define EXPECT_NATIVE_CALLS_BETWEEN(%0,%1,%2) \
    __pawntest_expect_calls(#%0, (%1), (%2), __file, __line)

#define EXPECT_NATIVE_CALLED(%0) \
    __pawntest_expect_calls(#%0, 1, 2147483647, __file, __line)

#define EXPECT_NATIVE_NOT_CALLED(%0) \
    __pawntest_expect_calls(#%0, 0, 0, __file, __line)

#define EXPECT_NATIVE_ARG(%0,%1,%2,%3) \
    __pawntest_expect_arg(#%0, (%1), (%2), (%3), __file, __line)

#define EXPECT_NATIVE_STRING_ARG(%0,%1,%2,%3) \
    __pawntest_expect_string_arg(#%0, (%1), (%2), (%3), __file, __line)

#define EXPECT_NATIVE_ORDER(%0) \
    __pawntest_expect_order(#%0, __file, __line)

#define TEST_SCHEDULE(%0,%1) \
    __pawntest_schedule((%0), #%1)

#define TEST_ADVANCE_TIME(%0) \
    __pawntest_advance_time((%0))

#define TEST_RUN_PENDING() \
    __pawntest_run_pending()

#define TEST_NOW() \
    __pawntest_now()
