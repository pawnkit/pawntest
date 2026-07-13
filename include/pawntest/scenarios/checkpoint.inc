#if defined _pt_sc_checkpoint_inc
    #endinput
#endif
#define _pt_sc_checkpoint_inc

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
