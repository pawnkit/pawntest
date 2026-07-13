#if defined _pt_sc_http_inc
    #endinput
#endif
#define _pt_sc_http_inc

#define MOCK_HTTP_RESPONSE(%0,%1,%2) \
    __pt_http_response((%0), (%1), (%2))

#define MOCK_HTTP_RESPONSE_FOR(%0,%1,%2,%3) \
    __pt_http_method_response((%0), (%1), (%2), (%3))

#define ASSERT_HTTP_REQUESTS(%0) \
    if (!__pt_http_requests((%0), __file, __line)) return 0

#define ASSERT_HTTP_REQUEST(%0,%1,%2) \
    if (!__pt_http_request((%0), (%1), (%2), __file, __line)) return 0
