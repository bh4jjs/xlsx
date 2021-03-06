package xlsx

import (
	"github.com/plandem/xlsx/format"
	"github.com/plandem/xlsx/internal/ml"
	"github.com/plandem/xlsx/types"
	_ "unsafe"
)

//go:linkname fromConditionalFormat github.com/plandem/xlsx/format.fromConditionalFormat
func fromConditionalFormat(f *format.ConditionalFormat) (*ml.ConditionalFormatting, []*format.StyleFormat)

type conditionals struct {
	sheet *sheetInfo
}

//newConditionals creates an object that implements conditional formatting functionality
func newConditionals(sheet *sheetInfo) *conditionals {
	return &conditionals{sheet: sheet}
}

func (c *conditionals) initIfRequired() {
	//attach conditionals if required
	if c.sheet.ml.ConditionalFormatting == nil {
		var conditionals []*ml.ConditionalFormatting
		c.sheet.ml.ConditionalFormatting = &conditionals
	}
}

//Add adds a conditional formatting with attaching additional refs if required
func (c *conditionals) Add(conditional *format.ConditionalFormat, refs []types.Ref) error {
	c.initIfRequired()

	//attach additional refs, if required
	if len(refs) > 0 {
		conditional.Set(format.Conditions.Refs(refs...))
	}

	if err := conditional.Validate(); err != nil {
		return err
	}

	info, styles := fromConditionalFormat(conditional)
	if info != nil && len(styles) > 0 && len(info.Bounds) > 0 {
		for i, styleInfo := range styles {
			if styleInfo != nil {
				//add a new diff styles
				styleID := c.sheet.workbook.doc.styleSheet.addDiffStyle(styleInfo)
				info.Rules[i].Style = &styleID
			}

			//add a new conditional
			*c.sheet.ml.ConditionalFormatting = append(*c.sheet.ml.ConditionalFormatting, info)
		}
	}

	return nil
}

//Remove deletes a conditional formatting from refs
func (c *conditionals) Remove(refs []types.Ref) {
	panic(errorNotSupported)
}

//Resolve checks if requested cIdx and rIdx related to any conditionals formatting and returns it
func (c *conditionals) Resolve(cIdx, rIdx int) *format.ConditionalFormat {
	//TODO: Populate format.ConditionalFormat with required information
	panic(errorNotSupported)
}

func (c *conditionals) pack() *[]*ml.ConditionalFormatting {
	//conditionals must have at least one object
	if c.sheet.ml.ConditionalFormatting != nil && len(*c.sheet.ml.ConditionalFormatting) == 0 {
		c.sheet.ml.ConditionalFormatting = nil
	}

	return c.sheet.ml.ConditionalFormatting
}
