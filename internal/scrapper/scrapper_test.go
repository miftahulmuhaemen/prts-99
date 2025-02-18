// filepath: /home/sahana/Code_Repository/back-go/internal/scrapper_test.go
package internal

import (
	"chat-ak-wikia/internal/model"
	"testing"

	"github.com/gocolly/colly"
	"github.com/stretchr/testify/assert"
)

func TestScrapper(t *testing.T) {

	url := "https://arknights.wiki.gg/wiki/Operator/6-star"
	domain := "arknights.wiki.gg"
	collector := colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.CacheDir("./cache"),
	)

	// Call the Scrapper function
	operators, err := Scrapper(1, url, collector)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, operators, 1)

	// Check the first operator
	assert.Equal(t, "Aak", operators[0].OperatorName)
	assert.Equal(t, "Specialist", operators[0].Class)
	assert.Equal(t, "Geek", operators[0].Branch)
	assert.Equal(t, "Lee's Detective Agency", operators[0].Faction)
	assert.Equal(t, "Ranged", operators[0].Position)
	assert.ElementsMatch(t, []string{"Support", "DPS"}, operators[0].Tags)
	assert.Equal(t, "Continually loses HP over time", operators[0].Trait)

}

func TestUpdateOrCreateAttribute(t *testing.T) {
	tests := []struct {
		name              string
		attributes        map[string]model.Attribute
		attribute         string
		field             string
		value             interface{}
		expectedAttribute map[string]model.Attribute
	}{
		{
			name:       "Test with existing attribute",
			attributes: map[string]model.Attribute{},
			attribute:  "Base",
			field:      "HP",
			value:      100,
			expectedAttribute: map[string]model.Attribute{
				"Base": {
					HP: 100,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			updateOrCreateAttribute(tc.attributes, tc.attribute, tc.field, tc.value)
			assert.Equal(t, tc.expectedAttribute, tc.attributes)
		})
	}
}
