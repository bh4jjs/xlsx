package primitives_test

import (
	"encoding/xml"
	"fmt"
	"github.com/plandem/xlsx/format"
	"github.com/plandem/xlsx/internal/ml/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFontScheme(t *testing.T) {
	type Element struct {
		Property primitives.FontSchemeType `xml:"property,omitempty"`
	}

	list := map[string]primitives.FontSchemeType{
		"none":     format.FontSchemeNone,
		"major":    format.FontSchemeMajor,
		"minor":    format.FontSchemeMinor,
		"schema-a": primitives.FontSchemeType("schema-a"),
	}

	for s, v := range list {
		t.Run(s, func(tt *testing.T) {
			entity := Element{Property: v}
			encoded, err := xml.Marshal(&entity)

			require.Empty(tt, err)
			require.Equal(tt, fmt.Sprintf(`<Element><property val="%s"></property></Element>`, s), string(encoded))

			var decoded Element
			err = xml.Unmarshal(encoded, &decoded)
			require.Empty(tt, err)

			require.Equal(tt, entity, decoded)
		})
	}
}
