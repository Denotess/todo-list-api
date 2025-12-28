package models

import "testing"

func TestPage(t *testing.T) {
	tests := []struct {
		name   string
		query  TodoQuery
		expect int
	}{
		{name: "zero limit", query: TodoQuery{Limit: 0, Offset: 0}, expect: 0},
		{name: "first page", query: TodoQuery{Limit: 10, Offset: 0}, expect: 1},
		{name: "third page", query: TodoQuery{Limit: 10, Offset: 20}, expect: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Page(&test.query); got != test.expect {
				t.Fatalf("expected page %d, got %d", test.expect, got)
			}
		})
	}
}
