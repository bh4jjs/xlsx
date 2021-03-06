package types

import (
	"errors"
	"fmt"
	sharedML "github.com/plandem/ooxml/ml"
	"github.com/plandem/xlsx/format"
	"github.com/plandem/xlsx/internal"
	"github.com/plandem/xlsx/internal/ml"
	"github.com/plandem/xlsx/internal/validator"
	"regexp"
	"strings"
)

type HyperlinkInfo struct {
	hyperlink *ml.Hyperlink
	styleID   format.DirectStyleID
	linkType  hyperlinkType
}

type hyperlinkOption func(o *HyperlinkInfo)
type hyperlinkType byte

const (
	hyperlinkTypeUnknown hyperlinkType = iota
	hyperlinkTypeWeb
	hyperlinkTypeEmail
	hyperlinkTypeFile
)

//Hyperlink is a 'namespace' for all possible settings for hyperlink
var Hyperlink hyperlinkOption

//NewHyperlink creates and returns a new HyperlinkInfo object that holds settings for hyperlink and related styles
func NewHyperlink(options ...hyperlinkOption) *HyperlinkInfo {
	i := &HyperlinkInfo{
		hyperlink: &ml.Hyperlink{},
	}
	i.Set(options...)
	return i
}

//Set sets new options for hyperlink
func (i *HyperlinkInfo) Set(options ...hyperlinkOption) {
	for _, o := range options {
		o(i)
	}
}

//Validate validates hyperlink info and return error in case of invalid settings
func (i *HyperlinkInfo) Validate() error {
	switch i.linkType {
	case hyperlinkTypeUnknown:
		if len(i.hyperlink.Location) == 0 {
			return errors.New("unknown type of hyperlink")
		}
	case hyperlinkTypeWeb:
		if len(i.hyperlink.RID) > internal.UrlLimit {
			return errors.New(fmt.Sprintf("url exceeded maximum allowed length (%d chars)", internal.UrlLimit))
		}

		if len(i.hyperlink.RID) <= 3 {
			return errors.New("url is too short")
		}

		if strings.Contains(string(i.hyperlink.RID), "#") {
			return errors.New("url contains a pound sign (#)")
		}

		if !validator.IsURL(string(i.hyperlink.RID)) {
			return errors.New("url is not valid")
		}
	case hyperlinkTypeEmail:
		if len(i.hyperlink.RID) > internal.UrlLimit {
			return errors.New(fmt.Sprintf("email exceeded maximum allowed length (%d chars)", internal.UrlLimit))
		}

		if !validator.IsEmail(string(i.hyperlink.RID)) {
			if ok, info := validator.IsMailTo(string(i.hyperlink.RID)); ok && validator.IsEmail(info["email"]) {
				break
			}

			return errors.New("email is not valid")
		}
	case hyperlinkTypeFile:
		if len(i.hyperlink.RID) > internal.UrlLimit {
			return errors.New(fmt.Sprintf("link to file exceeded maximum allowed length (%d chars)", internal.UrlLimit))
		}

		if len(i.hyperlink.RID) <= 3 {
			return errors.New("filename is too short")
		}

		if strings.Contains(string(i.hyperlink.RID), "#") {
			return errors.New("filename contains a pound sign (#)")
		}
	}

	return nil
}

//Formatting returns style that will be used by hyperlink
func (i *HyperlinkInfo) Formatting() format.DirectStyleID {
	return i.styleID
}

//String returns text version of hyperlink info
func (i *HyperlinkInfo) String() string {
	target := string(i.hyperlink.RID)
	location := i.hyperlink.Location
	//
	//if len(location) > 0 && location[0] == '#' {
	//	location = location[1:]
	//}

	if len(location) > 0 {
		return fmt.Sprintf("%s#%s", target, location)
	}

	return target
}

func (o *hyperlinkOption) Formatting(styleID format.DirectStyleID) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		i.styleID = styleID
	}
}

func (o *hyperlinkOption) Tooltip(tip string) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		i.hyperlink.Tooltip = tip
	}
}

func (o *hyperlinkOption) Display(display string) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		i.hyperlink.Display = display
	}
}

//ToMail sets target to email
func (o *hyperlinkOption) ToMail(address, subject string) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		if len(subject) > 0 {
			i.hyperlink.RID = sharedML.RID(fmt.Sprintf("mailto:%s?subject=%s", address, subject))
		} else {
			i.hyperlink.RID = sharedML.RID(fmt.Sprintf("mailto:%s", address))
		}

		i.linkType = hyperlinkTypeEmail
	}
}

//ToUrl sets target to web site
func (o *hyperlinkOption) ToUrl(address string) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		i.hyperlink.RID = sharedML.RID(escapeTarget(strings.TrimRight(address, `/`)))
		i.linkType = hyperlinkTypeWeb
	}
}

//ToFile sets target to external file
func (o *hyperlinkOption) ToFile(fileName string) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		//change the directory separator from Unix to DOS
		fileName = strings.Replace(fileName, "/", "\\", -1)

		//add the file:/// URI to the url for Windows style "C:/" link and network shares
		if validator.IsWinPath(fileName) {
			fileName = "file:///" + fileName
		}

		//convert a '.\dir\filename' link to 'dir\filename'
		re := regexp.MustCompile(`^\.\\`)
		fileName = re.ReplaceAllString(fileName, "")

		i.hyperlink.RID = sharedML.RID(escapeTarget(fileName))
		i.linkType = hyperlinkTypeFile
	}
}

//ToRef sets target to ref of sheet with sheetName. Omit sheetName to set location to ref of active sheet
func (o *hyperlinkOption) ToRef(ref Ref, sheetName string) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		if len(ref) > 0 {
			if len(sheetName) > 0 {
				//sheet + ref
				i.hyperlink.Location = fmt.Sprintf("%s!%s", escapeLocation(sheetName), ref)
			} else {
				//ref only, can be cell or bookmark
				i.hyperlink.Location = fmt.Sprintf("%s", ref)
			}
		}
	}
}

//ToBookmark sets target to bookmark, that can be named region in xlsx, bookmark of remote file or even site
func (o *hyperlinkOption) ToBookmark(location string) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		if len(location) > 0 {
			if location[0] == '#' {
				location = location[1:]
			}

			//ref only, can be cell or bookmark
			i.hyperlink.Location = fmt.Sprintf("%s", escapeLocation(location))
		}
	}
}

/*
ToTarget is very close to HYPERLINK function of Excel
 https://support.office.com/en-us/article/hyperlink-function-333c7ce6-c5ae-4164-9c47-7de9b76f577f

	a) to target: "target" or "[target]"
	b) to location at target: "[target]location" or "target#location"

Here are some examples of supported values:
	- same file, same sheet
	=HYPERLINK("#A1", "Reference to same sheet")

	- same file, other sheet
	=HYPERLINK("#SheetName!A1", "Reference to sheet without space in name")
	=HYPERLINK("#'Sheet Name'!A1", "Reference to sheet with space in name")

	- other local file
	=HYPERLINK("D:\Folder\File.docx","Word file")
	=HYPERLINK("D:\Folder\File.docx#Bookmark","Local Word file with bookmark")
	=HYPERLINK("D:\Folder\File.xlsx#SheetName!A1","Local Excel file with reference")
	=HYPERLINK("D:\Folder\File.xlsx#'Sheet Name'!A1","Local Excel file with reference")

	=HYPERLINK("[D:\Folder\File.docx]","Word file")
	=HYPERLINK("[D:\Folder\File.docx]Bookmark","Local Word file with bookmark")
	=HYPERLINK("[D:\Folder\File.xlsx]SheetName!A1","Local Excel file with reference")
	=HYPERLINK("[D:\Folder\File.xlsx]'Sheet Name'!A1","Local Excel file with reference")

	- other remote file
	=HYPERLINK("\\SERVER\Folder\File.doc", "Remote Word file")
	=HYPERLINK("\\SERVER\Folder\File.xlsx#SheetName!A1", "Remote Excel file with reference")
	=HYPERLINK("\\SERVER\Folder\File.xlsx#'Sheet Name'!A1", "Remote Excel file with reference")
	=HYPERLINK("[\\SERVER\Folder\File.xlsx]SheetName!A1", "Remote Excel file with reference")
	=HYPERLINK("[\\SERVER\Folder\File.xlsx]'Sheet Name'!A1", "Remote Excel file with reference")

	- url
	=HYPERLINK("https://www.spam.it","Website without bookmark")
	=HYPERLINK("https://www.spam.it/#bookmark","Website with bookmark")
	=HYPERLINK("[https://www.spam.it/]bookmark","Website with bookmark")

	-email
	=HYPERLINK("mailto:spam@spam.it","Email without subject")
	=HYPERLINK("mailto:spam@spam.it?subject=topic","Email with subject")
*/
func (o *hyperlinkOption) ToTarget(target string) hyperlinkOption {
	return func(i *HyperlinkInfo) {
		var location string

		//location is set using pound sign (#)
		if i := strings.LastIndexByte(target, '#'); i != -1 {
			location = target[i+1:]
			target = target[:i]
		} else if i = strings.LastIndexByte(target, ']'); target[0] == '[' && i != -1 {
			location = target[i+1:]
			target = target[1:i]
		}

		if len(location) > 0 {
			//TODO: potential corrupted location. Ideally it should be parsed and set via 'ToBookmark' or 'ToRef'
			i.hyperlink.Location = location
		}

		//detect type of link and call related method to set proper info
		if len(target) > 0 {
			if validator.IsURL(target) {
				i.Set(Hyperlink.ToUrl(target))
			} else if ok, mail := validator.IsMailTo(target); ok {
				i.Set(Hyperlink.ToMail(mail["email"], mail["subject"]))
			} else if validator.IsEmail(target) {
				i.Set(Hyperlink.ToMail(target, ""))
			} else if validator.IsFilePath(target) {
				i.Set(Hyperlink.ToFile(target))
			} else {
				panic(fmt.Sprintf("Can't detect type of hyperlink for target: %s", target))
			}
		}
	}
}

//private method used by hyperlinks manager to unpack HyperlinkInfo
func fromHyperlinkInfo(info *HyperlinkInfo) (hyperlink *ml.Hyperlink, styleID format.DirectStyleID, err error) {
	if err = info.Validate(); err != nil {
		return
	}

	styleID = info.styleID
	hyperlink = info.hyperlink
	return
}

//private method used by hyperlinks manager to pack HyperlinkInfo
func toHyperlinkInfo(link *ml.Hyperlink, target string, styleID format.DirectStyleID) *HyperlinkInfo {
	//normalize location
	location := link.Location
	if len(location) > 0 && location[0] != '#' {
		location = "#" + location
	}

	return NewHyperlink(
		Hyperlink.Formatting(styleID),
		Hyperlink.Display(link.Display),
		Hyperlink.Tooltip(link.Tooltip),
		Hyperlink.ToTarget(target+location),
	)
}

func escapeLocation(location string) string {
	// TODO: escape location (research what kind of escaping Excel is expecting)
	return `'` + strings.Replace(location, `'`, `\'`, -1) + `'`
}

func escapeTarget(target string) string {
	//pound symbol (#) is not allowed in target
	return strings.Replace(target, `#`, `%23`, -1)
}
