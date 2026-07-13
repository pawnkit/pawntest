#if defined _pawntest_provider_included
    #endinput
#endif
#define _pawntest_provider_included

native bool:__pawntest_provider_register(const native_name[], const handler_name[]);
native __pawntest_provider_arg_cell(index);
native bool:__pawntest_provider_arg_string(index, output[], size = sizeof(output));
native __pt_provider_arg_array(index, offset);
native bool:__pawntest_provider_set_cell(index, value);
native bool:__pt_provider_set_array(index, offset, value);
native bool:__pawntest_provider_set_string(index, const value[], size);
native bool:__pawntest_provider_call(const public_name[], {Float,_}:...);

#define PROVIDE_NATIVE(%0,%1) \
    __pawntest_provider_register(#%0, #%1)

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

#define PROVIDER_CALL1(%0,%1) \
    __pawntest_provider_call((%0), (%1))

#define PROVIDER_CALL2(%0,%1,%2) \
    __pawntest_provider_call((%0), (%1), (%2))

#define PROVIDER_CALL3(%0,%1,%2,%3) \
    __pawntest_provider_call((%0), (%1), (%2), (%3))
