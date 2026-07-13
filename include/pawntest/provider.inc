#if defined _pawntest_provider_included
    #endinput
#endif
#define _pawntest_provider_included

#define PAWNTEST_PROVIDER_ABI (1)

forward PawntestProviderVersion();
public PawntestProviderVersion()
{
    return PAWNTEST_PROVIDER_ABI;
}

native bool:__pawntest_provider_register(const native_name[], const handler_name[], const signature[]);
native __pawntest_provider_arg_cell(index);
native bool:__pawntest_provider_arg_string(index, output[], size = sizeof(output));
native __pt_provider_arg_array(index, offset);
native bool:__pawntest_provider_set_cell(index, value);
native bool:__pt_provider_set_array(index, offset, value);
native bool:__pawntest_provider_set_string(index, const value[], size);
native bool:__pawntest_provider_call(const public_name[], {Float,_}:...);
native bool:__pt_provider_test_name(output[], size = sizeof(output));

#define PROVIDE_NATIVE(%0,%1,%2) \
    __pawntest_provider_register(#%0, #%1, (%2))

#define PROVIDER_ARG_CELL(%0) \
    __pawntest_provider_arg_cell((%0))

#define PROVIDER_ARG_STRING(%0,%1) \
    __pawntest_provider_arg_string((%0), (%1), sizeof(%1))

#define PROVIDER_ARG_ARRAY_CELL(%0,%1) \
    __pt_provider_arg_array((%0), (%1))

#define PROVIDER_SET_CELL(%0,%1) \
    __pawntest_provider_set_cell((%0), (%1))

#define PROVIDER_SET_ARRAY_CELL(%0,%1,%2) \
    __pt_provider_set_array((%0), (%1), (%2))

#define PROVIDER_SET_STRING(%0,%1,%2) \
    __pawntest_provider_set_string((%0), (%1), (%2))

#define PROVIDER_CALL0(%0) \
    __pawntest_provider_call((%0))

#define PROVIDER_TEST_NAME(%0) \
    __pt_provider_test_name((%0), sizeof(%0))

#define PROVIDER_CALL1(%0,%1) \
    __pawntest_provider_call((%0), (%1))

#define PROVIDER_CALL2(%0,%1,%2) \
    __pawntest_provider_call((%0), (%1), (%2))

#define PROVIDER_CALL3(%0,%1,%2,%3) \
    __pawntest_provider_call((%0), (%1), (%2), (%3))

#define PROVIDER_CALL4(%0,%1,%2,%3,%4) \
    __pawntest_provider_call((%0), (%1), (%2), (%3), (%4))

#define PROVIDER_CALL5(%0,%1,%2,%3,%4,%5) \
    __pawntest_provider_call((%0), (%1), (%2), (%3), (%4), (%5))

#define PROVIDER_CALL6(%0,%1,%2,%3,%4,%5,%6) \
    __pawntest_provider_call((%0), (%1), (%2), (%3), (%4), (%5), (%6))

#define PROVIDER_CALL7(%0,%1,%2,%3,%4,%5,%6,%7) \
    __pawntest_provider_call((%0), (%1), (%2), (%3), (%4), (%5), (%6), (%7))

#define PROVIDER_CALL8(%0,%1,%2,%3,%4,%5,%6,%7,%8) \
    __pawntest_provider_call((%0), (%1), (%2), (%3), (%4), (%5), (%6), (%7), (%8))
