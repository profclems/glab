package config

import "testing"

func BenchmarkGetEnv(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetKeyValueInFile(configFile, "foo")
	}
}
