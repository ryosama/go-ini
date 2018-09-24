/*Package ini permits interaction with ini files (configuration file format for windows).
You can parse, read values, set values and save your ini files

Basic Usage :

import "github.com/ryosama/go-ini"

myIni := new(ini.Ini)
if err := myIni.LoadFromFile("config.ini") ; err != nil {
	panic("Unable to load configuration : " + err.Error())
}

myHost, _ := myIni.Get("Server","host")

if myIni.Exists("Server","port") { // test if key exists
	myPort := int( myIni.Get("Server","port") )
} else {
	myPort := 80
}

myIni.Set("Server","host","127.0.0.1")
myIni.Save()
*/
package ini

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

/*
Ini Create a new Ini object

Example :

	myIni := new(ini.Ini)
*/
type Ini struct {
	data map[string]Section

	// Last filename pass to Load or Save
	Filename string

	// Add this string before section (default is "")
	SectionPrefix string

	// Add this string before every items (default is "  ")
	ItemPrefix string

	// Add this string after every items (default is " ")
	ItemSuffix string

	// Add this string before every values (default is " ")
	ValuePrefix string

	// Add this string before every new sections (except the first one) (default is "\r\n")
	SectionSeparator string

	// Add this string before every new item (except the first one) (default is "\r\n")
	ItemSeparator string

	// If set to false, remove all the comments while saving (default is true)
	WithComments bool

	// Caractere(s) from prefixing a comment (default is "; ")
	CommentPrefix string
}

// Section has items and comments
type Section struct {
	items    map[string]Item
	comments []string
}

// Item has value and comments
type Item struct {
	value    string
	comments []string
}

/*
LoadFromFile reads ini format from a file

Example :

	err := myIni.LoadFromFile("config.ini")
*/
func (ini *Ini) LoadFromFile(filename string) error {
	ini.Filename = filename
	content, err := ioutil.ReadFile(ini.Filename)
	if err != nil {
		return err
	}
	tmp := string(content)
	ini.LoadFromString(&tmp)
	return nil
}

/*
LoadFromString loads data from a string pointer

Example :

	content := `

	[section1]

	item1=value1`

	myini.LoadFromString( &content )
*/
func (ini *Ini) LoadFromString(content *string) {

	// default value for formating
	ini.WithComments = true
	ini.CommentPrefix = "; "
	ini.SectionSeparator = "\r\n"
	ini.ItemSeparator = "\r\n"
	ini.ItemPrefix = "  "
	ini.ItemSuffix = " "
	ini.ValuePrefix = " "

	currentSection := ""
	comments := make([]string, 0)

	contentArray := strings.Split(*content, "\n") // split lines into array

	reComment := regexp.MustCompile("\\s*[;#]\\s*(.*)")
	reSection := regexp.MustCompile("\\s*\\[\\s*(.+?)\\s*\\]")
	reItem := regexp.MustCompile("\\s*(.+?)\\s*=\\s*(.+)\\s*")

	ini.data = make(map[string]Section)

	for _, value := range contentArray { // for each line

		if matches := reComment.FindStringSubmatch(value); matches != nil { // a comment
			comments = append(comments, strings.TrimSpace(matches[1]))

		} else if matches := reSection.FindStringSubmatch(value); matches != nil { // a section
			section := strings.TrimSpace(matches[1])
			currentSection = section // set active section

			var d Section
			d.comments = comments
			ini.data[currentSection] = d
			comments = make([]string, 0) // clears comments

		} else if matches := reItem.FindStringSubmatch(value); matches != nil { // an item
			name := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])

			var tmp Item
			tmp.comments = comments
			tmp.value = value

			s := ini.data[currentSection]
			if s.items == nil { // create structure for the first time
				s.items = make(map[string]Item)
			}
			s.items[name] = tmp
			ini.data[currentSection] = s
			comments = make([]string, 0) // clears comments

		}
	}
}

/*GetSections returns all the sections of the ini file */
func (ini *Ini) GetSections() []string {
	sections := make([]string, len(ini.data))
	i := 0
	for key := range ini.data {
		sections[i] = key
		i++
	}
	return sections
}

/*GetItems returns all the items of the ini file for a given section */
func (ini *Ini) GetItems(section string) []string {
	items := make([]string, len(ini.data[section].items))
	i := 0
	for key := range ini.data[section].items {
		items[i] = key
		i++
	}
	return items
}

/*SectionExists returns true or false if the section exists */
func (ini *Ini) SectionExists(section string) bool {
	_, exists := ini.data[section]
	return exists
}

/*ItemExists returns true or false if the item exists for the given section */
func (ini *Ini) ItemExists(section string, item string) bool {
	_, exists := ini.data[section]
	if exists {
		_, exists := ini.data[section].items[item]
		return exists
	}
	return false
}

/*Exists is an alias for ItemExists */
func (ini *Ini) Exists(section string, item string) bool {
	return ini.ItemExists(section, item)
}

/*
GetItem returns the items value of the ini file for a given section and item (and true as second return value)

If the item does not exists, return false as second return value

Example :

	value, success := myini.GetItem("section1","item1")
*/
func (ini *Ini) GetItem(section string, item string) (string, bool) {
	if ini.ItemExists(section, item) {
		return ini.data[section].items[item].value, true
	}
	return "", false
}

/*Get is an alias for GetItem */
func (ini *Ini) Get(section string, item string) (string, bool) {
	return ini.GetItem(section, item)
}

/*
SetItem sets the value of an item

Returns true if success, false section or item does not exists
*/
func (ini *Ini) SetItem(section string, item string, value string) bool {
	if ini.ItemExists(section, item) {
		i := ini.data[section].items[item]
		i.value = value
		ini.data[section].items[item] = i
		return true
	}
	return false
}

/*Set is an alias for SetItem */
func (ini *Ini) Set(section string, item string, value string) bool {
	return ini.SetItem(section, item, value)
}

/*
RenameSection renames a section

Returns true if success, false if section does not exists
*/
func (ini *Ini) RenameSection(oldName string, newName string) bool {
	if ini.SectionExists(oldName) {
		ini.data[newName] = ini.data[oldName]
		delete(ini.data, oldName)
		return true
	}
	return false
}

/*
RenameItem renames an item

Returns true if success, false if section or item does not exists
*/
func (ini *Ini) RenameItem(section, oldName string, newName string) bool {
	if ini.ItemExists(section, oldName) {
		ini.data[section].items[newName] = ini.data[section].items[oldName]
		delete(ini.data[section].items, oldName)
		return true
	}
	return false
}

/*
AddSection adds a section

Returns true if success, false if section already exists
*/
func (ini *Ini) AddSection(section string) bool {
	if ini.SectionExists(section) {
		return false
	}
	var s Section
	s.items = make(map[string]Item)
	s.comments = make([]string, 0)
	ini.data[section] = s
	return true
}

/*
AddItem adds an item

Returns true if success, false if item already exists
*/
func (ini *Ini) AddItem(section string, item string, value string) bool {
	if ini.SectionExists(section) {
		if ini.ItemExists(section, item) {
			return false
		}

		var tmp Item
		tmp.value = value
		tmp.comments = make([]string, 0)
		ini.data[section].items[item] = tmp

	} else {
		// section does not exist --> create it
		ini.AddSection(section)
		ini.AddItem(section, item, value)
	}
	return true
}

// SetOrCreate sets a value for an item, create section and item if needed
func (ini *Ini) SetOrCreate(section string, item string, value string) {
	ini.AddItem(section, item, value)
	ini.Set(section, item, value)
}

// DeleteItem deletes an item, return true if succes, false if the item does not exists
func (ini *Ini) DeleteItem(section string, item string) bool {
	if ini.ItemExists(section, item) {
		delete(ini.data[section].items, item)
		return true
	}
	return false
}

// DeleteSection deletes a section, return true if succes, false if the section does not exists
func (ini *Ini) DeleteSection(section string) bool {
	if ini.SectionExists(section) {
		delete(ini.data, section)
		return true
	}
	return false
}

/*
GetSectionComments returns the comments, one per line, just before the section, return empty slice of string if section does not exists.

Example :

	for _, com := range sectionExists := myIni.GetSectionComments("mySection") {

		print("Comment for mySection", com, "\n")

	}
*/
func (ini *Ini) GetSectionComments(section string) []string {
	if ini.SectionExists(section) {
		return ini.data[section].comments
	}
	return make([]string, 0)
}

/*
GetItemComments returns the comments, one per line, just before the item, return empty slice of string if item does not exists.

Example :

	for _, com := range myIni.GetItemComments("mySection","myItem") {

		print("Comment for myItem", com, "\n")

	}
*/
func (ini *Ini) GetItemComments(section string, item string) []string {
	if ini.ItemExists(section, item) {
		return ini.data[section].items[item].comments
	}
	return make([]string, 0)
}

// AddItemComment adds a comment to an item, return true if succeed, false otherwise
func (ini *Ini) AddItemComment(section string, item string, comment string) bool {
	if ini.SectionExists(section) && ini.ItemExists(section, item) {
		tmp := ini.data[section].items[item]
		tmp.comments = append(tmp.comments, comment) // add the comment
		ini.data[section].items[item] = tmp
		return true
	}
	return false
}

// DeleteItemComments deletes all the comments of an item, return true if succeed, false otherwise
func (ini *Ini) DeleteItemComments(section string, item string) bool {
	if ini.SectionExists(section) && ini.ItemExists(section, item) {
		tmp := ini.data[section].items[item]
		tmp.comments = make([]string, 0) // clear comments
		ini.data[section].items[item] = tmp
		return true
	}
	return false
}

// DeleteItemComment deletes the comment number id, return true if succeed, false otherwise
func (ini *Ini) DeleteItemComment(section string, item string, id int) bool {
	if ini.SectionExists(section) && ini.ItemExists(section, item) {
		comments := ini.GetItemComments(section, item)
		ini.DeleteItemComments(section, item) // delete all the comment
		for i, com := range comments {
			if i != id {
				ini.AddItemComment(section, item, com) // add new comments except the deleted one
			}
		}
		return true
	}
	return false
}

// AddSectionComment adds a comment to a setion, return true if succeed, false otherwise
func (ini *Ini) AddSectionComment(section string, comment string) bool {
	if ini.SectionExists(section) {
		tmp := ini.data[section]
		tmp.comments = append(tmp.comments, comment) // add the comment
		ini.data[section] = tmp
		return true
	}
	return false
}

// DeleteSectionComments deletes all the comments of a section, return true if succeed, false otherwise
func (ini *Ini) DeleteSectionComments(section string) bool {
	if ini.SectionExists(section) {
		tmp := ini.data[section]
		tmp.comments = make([]string, 0) // clear comments
		ini.data[section] = tmp
		return true
	}
	return false
}

// DeleteSectionComment deletes the comment number id, return true if succeed, false otherwise
func (ini *Ini) DeleteSectionComment(section string, id int) bool {
	if ini.SectionExists(section) {
		comments := ini.GetSectionComments(section)
		ini.DeleteSectionComments(section) // delete all the comment
		for i, com := range comments {
			if i != id {
				ini.AddSectionComment(section, com) // add new comments except the deleted one
			}
		}
		return true
	}
	return false
}

/*
Save saves the ini format to a file

Example :

	err := myIni.Save() // use myIni.Filename to save

	err := myIni.Save("new_config.ini") // use new_config.ini and set myIni.Filename
*/
func (ini *Ini) Save(params ...string) error {
	if len(params) > 0 {
		ini.Filename = params[0]
	}

	if ini.Filename == "" {
		return errors.New("You must specify a filename before saving")
	}
	return ioutil.WriteFile(ini.Filename, []byte(ini.Sprint()), os.ModePerm)
}

/*
Sprint returns the ini format into a formatted string

TIPS :
You can set SectionPrefix,ItemPrefix, ItemSuffix, ValuePrefix, SectionSeparator, ItemSeparator, WithComments, CommentPrefix to tweak format aspect
*/
func (ini *Ini) Sprint() string {
	cr := "\r\n"
	s := ""

	sections := ini.GetSections()
	for i := 0; i < len(sections); i++ {
		if ini.WithComments { // add the sections comments
			for _, com := range ini.GetSectionComments(sections[i]) {
				s += ini.SectionPrefix + ini.CommentPrefix + com + cr
			}
		}

		s += ini.SectionPrefix + "[" + sections[i] + "]" + cr

		items := ini.GetItems(sections[i])
		for j := 0; j < len(items); j++ {
			if ini.WithComments { // add the item comments
				for _, com := range ini.GetItemComments(sections[i], items[j]) {
					s += ini.ItemPrefix + ini.CommentPrefix + com + cr
				}
			}

			value, _ := ini.GetItem(sections[i], items[j])
			s += ini.ItemPrefix + items[j] + ini.ItemSuffix + "=" + ini.ValuePrefix + value + cr

			if j != len(items)-1 {
				s += ini.ItemSeparator
			}

			if j == len(items)-1 && i != len(sections)-1 { // add section separator if last item
				s += ini.SectionSeparator
			}
		}
	}

	return s
}

/*
Print prints the ini format into a formatted string

TIPS :
You can set SectionPrefix,ItemPrefix, ItemSuffix, ValuePrefix, SectionSeparator, ItemSeparator, WithComments, CommentPrefix to tweak format aspect
*/
func (ini *Ini) Print() {
	print(ini.Sprint())
}
