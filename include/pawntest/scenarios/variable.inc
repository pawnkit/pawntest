#if defined _pt_sc_variable_inc
    #endinput
#endif
#define _pt_sc_variable_inc

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
