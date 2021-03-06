package xlsx

import (
	"errors"
	"fmt"
	"github.com/plandem/xlsx/format"
	"github.com/plandem/xlsx/internal"
	"github.com/plandem/xlsx/internal/ml"
	"github.com/plandem/xlsx/internal/ml/primitives"
	_ "unsafe"
)

//go:linkname toRichFont github.com/plandem/xlsx/format.toRichFont
func toRichFont(f *format.StyleFormat) *ml.RichFont

func toRichText(parts ...interface{}) (*ml.StringItem, error) {
	si := &ml.StringItem{}
	length := 0

	if len(parts) > 0 {
		//if last part is format, then remove it
		if _, lastIsFormat := parts[len(parts)-1].(*format.StyleFormat); lastIsFormat {
			parts = parts[:len(parts)-1]
		}

		richText := make([]*ml.RichText, 0)
		fontPart := true

		for i, p := range parts {
			switch v := p.(type) {
			case string:
				length += len(v)

				if !fontPart || i == 0 {
					//previous part was string or it's first part - add new block with a string and 'default format'
					richText = append(richText, &ml.RichText{
						Text: primitives.Text(v),
					})
				} else {
					//previous part was a format, so attach a string to prev block
					richText[len(richText)-1].Text = primitives.Text(v)
				}

				fontPart = false

			case *format.StyleFormat:
				if fontPart && i > 0 {
					return nil, errors.New("two styles in row is not allowed")
				}

				richText = append(richText, &ml.RichText{
					Font: toRichFont(v),
				})

				fontPart = true
			}
		}

		si.RichText = &richText
	}

	if length > internal.ExcelCellLimit {
		return nil, errors.New(fmt.Sprintf("text exceeds allowed length for cell value = %d", internal.ExcelCellLimit))
	}

	return si, nil
}

func fromRichText(text *ml.StringItem) (s string) {
	if text == nil {
		return
	}

	s += string(text.Text)

	if text.RichText != nil {
		for _, part := range *text.RichText {
			s += string(part.Text)
		}
	}

	return
}
