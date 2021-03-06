package primitives_test

import (
	"encoding/xml"
	"fmt"
	"github.com/plandem/xlsx/internal/ml/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateLinks(t *testing.T) {
	type Entity struct {
		Attribute primitives.UpdateLinksType `xml:"attribute,attr"`
	}

	list := map[string]primitives.UpdateLinksType{
		"":        primitives.UpdateLinksType(0),
		"userSet": primitives.UpdateLinksTypeUserSet,
		"never":   primitives.UpdateLinksTypeNever,
		"always":  primitives.UpdateLinksTypeAlways,
	}

	for s, v := range list {
		t.Run(s, func(tt *testing.T) {
			entity := Entity{Attribute: v}
			encoded, err := xml.Marshal(&entity)

			require.Empty(tt, err)
			if s == "" {
				require.Equal(tt, `<Entity></Entity>`, string(encoded))
			} else {
				require.Equal(tt, fmt.Sprintf(`<Entity attribute="%s"></Entity>`, s), string(encoded))
			}

			var decoded Entity
			err = xml.Unmarshal(encoded, &decoded)
			require.Empty(tt, err)

			require.Equal(tt, entity, decoded)
			require.Equal(tt, s, decoded.Attribute.String())
		})
	}
}
