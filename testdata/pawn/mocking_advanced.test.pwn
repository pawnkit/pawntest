#include <pawntest>

native mock_sequence();
native mock_scalar_output(&value);
native mock_string_output(value[], size);
native mock_callback(value);

new callback_value;

forward on_mock_callback(value);
public on_mock_callback(value)
{
    callback_value = value;
}

TEST(mock_return_sequence)
{
    MOCK_RETURN(mock_sequence, 30);
    MOCK_RETURN_ONCE(mock_sequence, 10);
    MOCK_RETURN_ONCE(mock_sequence, 20);

    ASSERT_EQ(mock_sequence(), 10);
    ASSERT_EQ(mock_sequence(), 20);
    ASSERT_EQ(mock_sequence(), 30);
}

TEST(mock_output_parameters)
{
    new value;
    new text[16];

    MOCK_RETURN(mock_scalar_output, 1);
    MOCK_OUTPUT(mock_scalar_output, 0, 42);
    ASSERT_EQ(mock_scalar_output(value), 1);
    ASSERT_EQ(value, 42);

    MOCK_RETURN(mock_string_output, 1);
    MOCK_OUTPUT_STRING(mock_string_output, 0, sizeof text, "hello");
    ASSERT_EQ(mock_string_output(text, sizeof text), 1);
    ASSERT_STR_EQ(text, "hello");
}

TEST(mock_callback)
{
    callback_value = 0;
    MOCK_RETURN(mock_callback, 1);
    MOCK_CALLBACK(mock_callback, on_mock_callback);

    ASSERT_EQ(mock_callback(77), 1);
    ASSERT_EQ(callback_value, 77);
}
