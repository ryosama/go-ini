// This package permits interaction with ini files (configuration file format for windows)
// You can parse, read values, set values and save your ini files
package ini

import (
	"io/ioutil"
	"strings"
	"regexp"
	"os"
	"errors"
)


/*
Create a new Ini object

Example :

myini := new(ini.Ini)
*/
type Ini struct {
	filename string
	sections map[string]map[string] string

	// default "" : Add this string before every items
	ItemPrefix string

	// default "" : Add this string before every new sections (except the first one)
	SectionSeparator string
}

/*
filename : Read ini format from this file

Example :

err := myini.LoadFromFile("config.ini")
*/
func (this *Ini) LoadFromFile(filename string) error {
	this.filename = filename
	content, err := ioutil.ReadFile(this.filename)
	if err != nil {
       return err
    }
    t := string(content)
    this.LoadFromString(&t)
    return nil
}

/*
Load data from a string pointer

Example :

content := "[section1] \n item1=value1"

myini.LoadFromString( &content )
*/
func (this *Ini) LoadFromString(content *string) {

	current_section := ""
    content_array := strings.Split( *content, "\n")

    re_comment 	:= regexp.MustCompile("\\s*;")
    re_section 	:= regexp.MustCompile("\\s*\\[\\s*(.+?)\\s*\\]")
	re_item 	:= regexp.MustCompile("\\s*(.+?)\\s*=\\s*(.+)\\s*")

   	this.sections = make(map[string]map[string]string)

    for _, value := range content_array {

    	if re_comment.MatchString( value ) { // un commentaire
    		
		} else if 	matches := re_section.FindStringSubmatch( value ) ; matches != nil { // un section
			this.sections[matches[1]] = map[string]string{}
			current_section = matches[1]

    	} else if 	matches := re_item.FindStringSubmatch( value ) ; matches != nil { // un item
			this.sections[current_section][matches[1]] = strings.TrimSpace(matches[2])
    	}
    }
}

/*
Return the ini format into a formatted string

TIPS : You can set myIniObj.ItemPrefix and myIniObj.SectionSeparator to tweak format aspect
*/
func (this *Ini) Sprint() string {
	s := ""

	sections := this.GetSections()
	for i:=0 ; i<len(sections) ; i++ {
		s += "[" + sections[i] + "]\r\n"
		items := this.GetItems(sections[i])
		for j:=0 ; j<len(items) ; j++ {
			item, _ := this.GetItem(sections[i],items[j])
			s += this.ItemPrefix + items[j] + "=" + item + "\r\n"
			if j==len(items)-1 {
				s += this.SectionSeparator
			}
		}
	}

	return s
}

/*
Print the ini format into a formatted string

TIPS : You can set myIniObj.ItemPrefix and myIniObj.SectionSeparator to tweak format aspect
*/
func (this *Ini) Print() {
	print(this.Sprint())
}

/* Returns all the sections of the ini file */
func (this *Ini) GetSections() []string {
	sections := make( []string, len(this.sections) )
	i := 0
	for key,_ := range this.sections {
		sections[i] = key
		i++
	}
	return sections
}

/* Returns all the items of the ini file for a given section */
func (this *Ini) GetItems(section string) []string {
	items := make( []string, len(this.sections[section]) )
	i := 0
	for key,_ := range this.sections[section] {
		items[i] = key
		i++
	}
	return items
}

/*
Returns the items value of the ini file for a given section and item (and true as second return value)

If the item does not exists, return false as second return value

Example :

value, success := myini.GetItem("section1","item1")
*/
func (this *Ini) GetItem(section string, item string) (string, bool) {
	if this.ItemExists(section, item) {
		return this.sections[section][item], true
	} else {
		return "", false
	}
}

/* Alias for GetItem */
func (this *Ini) Get(section string, item string) (string, bool) {
	return this.GetItem(section, item)
}

/* Returns true or false if the section exists */
func (this *Ini) SectionExists(section string) bool {
	_, exists := this.sections[section]
	return exists
}

/* Returns true or false if the item exists for the given section */
func (this *Ini) ItemExists(section string, item string) bool {
	_, exists := this.sections[section]
	if exists {
		_, exists := this.sections[section][item]
		return exists
	} else {
		return false
	}
}

/* Alias for ItemExists */
func (this *Ini) Exists(section string, item string) bool {
	return this.ItemExists(section, item)
}

/*
Set the value of an item

Returns true if success, false section or item does not exists
*/
func (this *Ini) SetItem(section string, item string, value string) bool {
	if this.ItemExists(section,item) {
		this.sections[section][item] = value
		return true
	} else {
		return false
	}
}

/* Alias for SetItem */
func (this *Ini) Set(section string, item string, value string) bool {
	return this.SetItem(section, item, value)
}

/*
Rename a section

Returns true if success, false if section does not exists
*/
func (this *Ini) RenameSection(oldName string, newName string) bool {
	newSections := make(map[string]map[string]string)
	var state bool = false
	for _, section := range this.GetSections() {
		if section == oldName { // rename
			newSections[newName] = this.sections[oldName]
			state = true
		} else { // only copy
			newSections[section] = this.sections[section]
		}
	}
	this.sections = newSections
	return state
}

/*
Rename an item

Returns true if success, false if section or item does not exists
*/
func (this *Ini) RenameItem(section, oldName string, newName string) bool {
	newItems := make(map[string]string)
	var state bool = false
	for _, item := range this.GetItems(section) {
		if item == oldName { // rename
			newItems[newName] = this.sections[section][oldName]
			state = true
		} else { // only copy
			newItems[item] = this.sections[section][item]
		}
	}
	this.sections[section] = newItems
	return state
}

/*
Add a section

Returns true if success, false if section already exists
*/
func (this *Ini) AddSection(section string) bool {
	if this.SectionExists(section) {
		return false
	}
	this.sections[section] = make(map[string]string)
	return true
}

/* 
Add an item

Returns true if success, false if item already exists
*/
func (this *Ini) AddItem(section string, item string, value string) bool {
	if this.SectionExists(section) {
		if this.ItemExists(section,item) {
			return false
		}
		this.sections[section][item] = value
		return true

	} else {
		// section does not exist --> create it
		this.AddSection(section)
		this.AddItem(section, item, value)
		return true
	}
}

// Set a value for an item, create section and item if needed
func (this *Ini) SetOrCreate(section string, item string, value string) {
	this.AddItem(section,item,value)
	this.Set(section,item,value)
}

/*
Save() error

Save(filename string) error

Save the ini format to a file

Example :

err := myini.Save("new_config.ini")
*/
func (this *Ini) Save(params ...string) error {
	if len(params)>0 {
		this.filename = params[0]
	}

	if this.filename == "" {
		return errors.New("You must specify a filename before saving")
	}

	return ioutil.WriteFile(this.filename,  []byte( this.Sprint() ), os.ModePerm)
}