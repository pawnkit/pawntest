#if defined _pt_sc_database_inc
    #endinput
#endif
#define _pt_sc_database_inc

#define ASSERT_DATABASE_CONNECTIONS(%0) \
    if (!__pt_database_connections((%0), __file, __line)) return 0

#define ASSERT_DATABASE_RESULTS(%0) \
    if (!__pt_database_results((%0), __file, __line)) return 0
