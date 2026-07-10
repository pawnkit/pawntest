package runner

import "testing"

func TestSeedForTestIsStableAndTestSpecific(t *testing.T) {
	first := seedForTest(42, "test_one")
	if first != seedForTest(42, "test_one") {
		t.Fatal("seed was not stable")
	}

	if first == seedForTest(42, "test_two") {
		t.Fatal("different tests received the same derived seed")
	}
}
