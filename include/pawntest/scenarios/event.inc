#if defined _pt_sc_event_inc
    #endinput
#endif
#define _pt_sc_event_inc

#define TEST_CONNECT_PLAYER(%0) \
    __pt_event_connect((%0))

#define TEST_DISCONNECT_PLAYER(%0,%1) \
    __pt_event_disconnect((%0), (%1))

#define TEST_SPAWN_PLAYER(%0) \
    __pt_event_spawn((%0))

#define TEST_KILL_PLAYER(%0,%1,%2) \
    __pt_event_death((%0), (%1), (%2))

#define TEST_PLAYER_TEXT(%0,%1) \
    __pt_event_text((%0), (%1))

#define TEST_PLAYER_COMMAND(%0,%1) \
    __pt_event_command((%0), (%1))

#define TEST_PLAYER_KEYS(%0,%1,%2,%3) \
    __pt_event_keys((%0), (%1), (%2), (%3))

#define TEST_ENTER_VEHICLE(%0,%1,%2) \
    __pt_event_enter_vehicle((%0), (%1), (%2))

#define TEST_EXIT_VEHICLE(%0) \
    __pt_event_exit_vehicle((%0))

#define TEST_PICKUP(%0,%1) \
    __pt_event_pickup((%0), (%1))

#define TEST_PLAYER_PICKUP(%0,%1) \
    __pt_event_player_pickup((%0), (%1))

#define TEST_ENTER_CHECKPOINT(%0) \
    __pt_event_enter_checkpoint((%0))

#define TEST_LEAVE_CHECKPOINT(%0) \
    __pt_event_leave_checkpoint((%0))

#define TEST_ENTER_RACE_CHECKPOINT(%0) \
    __pt_event_enter_race_cp((%0))

#define TEST_LEAVE_RACE_CHECKPOINT(%0) \
    __pt_event_leave_race_cp((%0))

#define TEST_CLICK_TEXTDRAW(%0,%1) \
    __pt_event_click_textdraw((%0), (%1))

#define TEST_CLICK_PLAYER_TEXTDRAW(%0,%1) \
    __pt_event_click_player_td((%0), (%1))

#define TEST_DAMAGE_PLAYER(%0,%1,%2,%3,%4) \
    __pt_event_damage_player((%0), (%1), (%2), (%3), (%4))

#define TEST_DAMAGE_ACTOR(%0,%1,%2,%3,%4) \
    __pt_event_damage_actor((%0), (%1), (%2), (%3), (%4))

#define TEST_WEAPON_SHOT(%0,%1,%2,%3,%4,%5,%6) \
    __pt_event_weapon_shot((%0), (%1), (%2), (%3), (%4), (%5), (%6))

#define TEST_FINISH_OBJECT_MOVE(%0) \
    __pt_event_object_moved((%0))

#define TEST_FINISH_PLAYER_OBJECT_MOVE(%0,%1) \
    __pt_event_player_obj_moved((%0), (%1))

#define TEST_STREAM_PLAYER_IN(%0,%1) \
    __pt_event_player_stream((%0), (%1), true)

#define TEST_STREAM_PLAYER_OUT(%0,%1) \
    __pt_event_player_stream((%0), (%1), false)

#define TEST_STREAM_VEHICLE_IN(%0,%1) \
    __pt_event_vehicle_stream((%0), (%1), true)

#define TEST_STREAM_VEHICLE_OUT(%0,%1) \
    __pt_event_vehicle_stream((%0), (%1), false)

#define TEST_STREAM_ACTOR_IN(%0,%1) \
    __pt_event_actor_stream((%0), (%1), true)

#define TEST_STREAM_ACTOR_OUT(%0,%1) \
    __pt_event_actor_stream((%0), (%1), false)

#define TEST_MOVE_PLAYER(%0,%1,%2,%3) \
    __pt_event_move_player((%0), (%1), (%2), (%3))

#define TEST_DAMAGE_VEHICLE(%0,%1,%2) \
    __pt_event_vehicle_damage((%0), (%1), (%2))

#define TEST_VEHICLE_DAMAGE_STATUS(%0,%1,%2,%3,%4,%5) \
    __pt_event_vehicle_status((%0), (%1), (%2), (%3), (%4), (%5))

#define TEST_RESPAWN_VEHICLE(%0) \
    __pt_event_vehicle_respawn((%0))

#define TEST_CANCEL_TEXTDRAW_SELECTION(%0) \
    __pt_event_cancel_textdraw((%0))

#define TEST_CLICK_PLAYER(%0,%1,%2) \
    __pt_event_click_player((%0), (%1), (%2))

#define TEST_CLICK_MAP(%0,%1,%2,%3) \
    __pt_event_click_map((%0), (%1), (%2), (%3))
