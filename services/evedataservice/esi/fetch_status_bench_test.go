package esi

import (
	"context"
	"github.com/andrewapj/arcturus/config"
	"testing"
)

func BenchmarkClient_FetchStatus(b *testing.B) {

	config.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()

	for i := 0; i < b.N; i++ {
		_, err := client.FetchStatus(context.Background())
		if err != nil {
			b.Fatal(err.Error())
		}
	}
	b.ReportAllocs()
}
