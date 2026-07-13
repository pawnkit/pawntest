#define PAWNTEST_STRICT_SCENARIOS
#include <support>

TEST(service_status)
{
    USE_FIXTURE(spawned_player);
    MOCK_HTTP_RESPONSE_FOR(HTTP_GET, "api.example.test/status", 200, "ok");

    ASSERT_TRUE(RequestServiceStatus(testPlayer));
    ASSERT_HTTP_REQUESTS(1);
    ASSERT_HTTP_REQUEST(HTTP_GET, "api.example.test/status", "");
    ASSERT_PLAYER_MESSAGE(testPlayer, "Service online");
}

TEST(account_database)
{
    new DB:database = OpenAccountDatabase();
    new DBResult:result = DB_ExecuteQuery(database, "INSERT INTO accounts VALUES ('Alice')");
    DB_FreeResultSet(result);

    result = DB_ExecuteQuery(database, "SELECT name FROM accounts");
    ASSERT_EQ(DB_GetRowCount(result), 1);
    DB_FreeResultSet(result);
    DB_Close(database);

    ASSERT_DATABASE_CONNECTIONS(0);
    ASSERT_DATABASE_RESULTS(0);
}

TEST_TAG(int_service_status)
TEST_TAG(int_account_database)
