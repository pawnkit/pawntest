#if defined _pt_sc_vehicle_inc
    #endinput
#endif
#define _pt_sc_vehicle_inc

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
