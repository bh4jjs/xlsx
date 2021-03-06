package format

import (
	"github.com/plandem/xlsx/internal/ml"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNumberFormat(t *testing.T) {
	style := NewStyles(
		NumberFormatID(8),
	)

	require.IsType(t, &StyleFormat{}, style)
	require.Equal(t, createStylesAndFill(func(f *StyleFormat) {
		f.styleInfo.NumberFormat = &ml.NumberFormat{
			ID:   8,
			Code: "($#,##0.00_);[RED]($#,##0.00_)",
		}
	}), style)

	style = NewStyles(
		NumberFormat(`$0.00" usd"`),
	)

	require.Equal(t, createStylesAndFill(func(f *StyleFormat) {
		f.styleInfo.NumberFormat = &ml.NumberFormat{
			ID:   -1,
			Code: `$0.00" usd"`,
		}
	}), style)
}
