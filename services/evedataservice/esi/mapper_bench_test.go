package esi

import (
	"github.com/andrewapj/arcturus/clock"
	"testing"
)

func BenchmarkMapToIds(b *testing.B) {

	responses := []*response{
		{body: []byte("[1,2,3]"), expires: clock.GetTime(), pages: 2, statusCode: 200},
		{body: []byte("[4,5,6]"), expires: clock.GetTime(), pages: 2, statusCode: 200}}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := mapToIds(responses)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapToSingle(b *testing.B) {

	responses := []*response{
		{body: []byte(esiStatusResponse), expires: clock.GetTime(), pages: 1, statusCode: 200}}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := mapToSingle[*Status](responses)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}
