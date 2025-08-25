package rss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemGenId(t *testing.T) {

	testCases := []struct {
		name string
		item Item
	}{
		{
			name: "Existing GUID",
			item: Item{
				Guid:  "1234567890",
				Links: []string{"https://example.com"},
			},
		},
		{
			name: "No GUID",
			item: Item{
				Guid:  "",
				Links: []string{"https://example.com"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			id := tc.item.GenId()
			assert.NotEmpty(id)
		})
	}

}
