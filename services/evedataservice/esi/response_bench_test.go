package esi

import (
	"github.com/andrewapj/arcturus/config"
	"testing"
)

func BenchmarkNewResponse(b *testing.B) {

	config.SetTestConfig()

	for i := 0; i < b.N; i++ {
		httpResponse := buildTestHttpResponse()
		_, err := newResponse(httpResponse)
		if err != nil {
			b.Fatal(err.Error())
		}
	}

	b.ReportAllocs()
}
