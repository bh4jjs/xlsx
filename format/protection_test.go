package format

import (
	"github.com/plandem/xlsx/internal/ml"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProtection(t *testing.T) {
	style := NewStyles(
		Protection.Hidden,
		Protection.Locked,
	)

	require.IsType(t, &StyleFormat{}, style)
	require.Equal(t, createStylesAndFill(func(f *StyleFormat) {
		f.styleInfo.Protection = &ml.CellProtection{
			Locked: true,
			Hidden: true,
		}
	}), style)
}
