package primitives

import (
	"encoding/xml"
	"strings"
)

//Bounds is implementation of RefList
type BoundsList []Bounds

func BoundsListFromRefs(refs ...Ref) BoundsList {
	var list []Bounds

	for _, r := range refs {
		list = append(list, r.ToBounds())
	}

	return list
}

//ToRefList returns refs. Alias of String() method
func (bl *BoundsList) ToRefList() RefList {
	return RefList(bl.String())
}

//String return textual version of BoundsList
func (bl BoundsList) String() string {
	var refs []string

	for _, b := range bl {
		refs = append(refs, b.String())
	}

	return strings.Join(refs, " ")
}

//Add adds a new ref into the BoundsList
func (bl *BoundsList) Add(ref Ref) {
	*bl = append(*bl, ref.ToBounds())
}

//IsEmpty return true if type was not initialized
func (bl BoundsList) IsEmpty() bool {
	return len(bl) == 0
}

//MarshalXMLAttr marshal BoundsList
func (bl *BoundsList) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	attr := xml.Attr{Name: name}

	if bl.IsEmpty() {
		attr = xml.Attr{}
	} else {
		attr.Value = bl.String()
	}

	return attr, nil
}

//UnmarshalXMLAttr unmarshal BoundsList
func (bl *BoundsList) UnmarshalXMLAttr(attr xml.Attr) error {
	*bl = RefList(attr.Value).ToBoundsList()
	return nil
}
