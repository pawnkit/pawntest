#define PAWNTEST_STRICT_SCENARIOS
#include <pawntest>

#define HTTP_GET (1)

native bool:HTTP(index, method, const url[], const data[], const callback[]);

forward OnStrictHTTP(index, response_code, data[]);
public OnStrictHTTP(index, response_code, data[])
{
    ASSERT_EQ(index, 7);
    ASSERT_EQ(response_code, 200);
    ASSERT_STR_EQ(data, "ok");
    return 1;
}

TEST(strict_http_scenario)
{
    MOCK_HTTP_RESPONSE_FOR(HTTP_GET, "api.example.test/health", 200, "ok");
    new bool:requested = HTTP(7, HTTP_GET, "api.example.test/health", "", "OnStrictHTTP");
    ASSERT_TRUE(requested);
    ASSERT_HTTP_REQUESTS(1);
}
