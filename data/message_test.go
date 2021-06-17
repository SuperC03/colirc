package data

import (
	"testing"

	"github.com/superc03/colirc/types"
)

func TestUnmarshalMessage(t *testing.T) {
	tests := map[string]*Message{
		"JOIN": {
			Command: types.CommandJOIN,
		},
		"JOIN colin :fsdf sfsdf sdfdsfd": {
			Command: types.CommandJOIN,
			Params:  []string{"colin", ":fsdf sfsdf sdfdsfd"},
		},
	}
	for k, v := range tests {
		if result, _ := UnmarshalMessage(k); result != v {
			t.Errorf("Expected %v, got %v", v, result)
		}
	}
}
