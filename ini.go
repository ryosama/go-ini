/*
This package permits interaction with ini files (configuration file format for windows).
You can parse, read values, set values and save your ini files

Basic Usage :

myIni := new(ini.Ini)

if err := myIni.LoadFromFile("config.ini") ; err {
	panic("Unable to load configuration : ", err)
}

myHost := myIni.Get("Server","host")

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
	"io/ioutil"
	"strings"
	"regexp"
	"errors"
	"os"
	_ "github.com/davecgh/go-spew/spew"
)


/*
Create a new Ini object

Example :

myIni := new(ini.Ini)
*/
type Ini struct {
	data 	 	map[string]Section

	// Last filename pass to Load or Save
	Filename 	string

	// Add this string before every items (default is "")
	ItemPrefix 	string

	// Add this string before every new sections (except the first one) (default is "")
	SectionSeparator string

	// If true Remove the comments while saving (default is false)
	NoComments bool
}


// A section has items and comments
type Section struct {
	items 	 map[string]Item
	comments []string
}

// An item has value and comments
type Item struct {
	value 	   string
	comments []string
}

/*
filename : Read ini format from this file

Example :

err := myIni.LoadFromFile("config.ini")
*/
func (this *Ini) LoadFromFile(filename string) error {
	this.Filename = filename
	content, err := ioutil.ReadFile(this.Filename)
	if err != nil {
       return err
    }
    tmp := string(content)
    this.LoadFromString(&tmp)
    return nil
}

/*
Load data from a string pointer

Example :

content := `

[section1]

item1=value1`

myini.LoadFromString( &content )
*/
func (this *Ini) LoadFromString(content *string) {

	current_section := ""
	comments 	    := make([]string,0)

    content_array := strings.Split( *content, "\n") // split lines into array

    re_comment 	:= regexp.MustCompile("\\s*[;#](.*)")
    re_section 	:= regexp.MustCompile("\\s*\\[\\s*(.+?)\\s*\\]")
	re_item 	:= regexp.MustCompile("\\s*(.+?)\\s*=\\s*(.+)\\s*")

   	this.data = make(map[string]Section)

    for _, value := range content_array { // for each line

    	if matches := re_comment.FindStringSubmatch( value ) ; matches != nil { // a comment
    		comments = append(comments, strings.TrimSpace(matches[1]) )

		} else if 	matches := re_section.FindStringSubmatch( value ) ; matches != nil { // a section
			section := strings.TrimSpace(matches[1])
			current_section = section 		// set active section

			var d Section
				d.comments = comments
			this.data[current_section] = d
			comments = make([]string,0) // clears comments
			
    	} else if 	matches := re_item.FindStringSubmatch( value ) ; matches != nil { // an item
    		name  := strings.TrimSpace(matches[1])
    		value := strings.TrimSpace(matches[2])

    		var tmp Item
				tmp.comments = comments
				tmp.value 	 = value

			s := this.data[current_section]
			if 	s.items == nil { // create structure for the first time
				s.items = make(map[string]Item)
			}
			s.items[name] = tmp
			this.data[current_section] = s
			comments = make([]string,0) // clears comments
		
    	}
    }
}


/* Returns all the sections of the ini file */
func (this *Ini) GetSections() []string {
	sections := make( []string, len(this.data) )
	i := 0
	for key,_ := range this.data {
		sections[i] = key ; i++
	}
	return sections
}

/* Returns all the items of the ini file for a given section */
func (this *Ini) GetItems(section string) []string {
	items := make( []string, len(this.data[section].items) )
	i := 0
	for key,_ := range this.data[section].items {
		items[i] = key ; i++
	}
	return items
}

/* Returns true or false if the section exists */
func (this *Ini) SectionExists(section string) bool {
	_, exists := this.data[section]
	return exists
}


/* Returns true or false if the item exists for the given section */
func (this *Ini) ItemExists(section string, item string) bool {
	_, exists := this.data[section]
	if exists {
		_, exists := this.data[section].items[item]
		return exists
	}
	return false
}

/* Alias for ItemExists */
func (this *Ini) Exists(section string, item string) bool {
	return this.ItemExists(section, item)
}

/*
Returns the items value of the ini file for a given section and item (and true as second return value)

If the item does not exists, return false as second return value

Example :

value, success := myini.GetItem("section1","item1")
*/
func (this *Ini) GetItem(section string, item string) (string, bool) {
	if this.ItemExists(section, item) {
		return this.data[section].items[item].value, true
	}
	return "", false
}

/* Alias for GetItem */
func (this *Ini) Get(section string, item string) (string, bool) {
	return this.GetItem(section, item)
}

/*
Set the value of an item

Returns true if success, false section or item does not exists
*/
func (this *Ini) SetItem(section string, item string, value string) bool {
	if this.ItemExists(section,item) {
		var i Item = this.data[section].items[item]
		i.value = value
		this.data[section].items[item] = i
		return true
	}
	return false
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
	if this.SectionExists(oldName) {
		this.data[newName] = this.data[oldName]
		delete(this.data, oldName)
		return true
	}
	return false
}

/*
Rename an item

Returns true if success, false if section or item does not exists
*/
func (this *Ini) RenameItem(section, oldName string, newName string) bool {
	if this.Exists(section, oldName) {
		this.data[section].items[newName] = this.data[section].items[oldName]
		delete(this.data[section].items, oldName)
		return true
	}
	return false
}

/*
Add a section

Returns true if success, false if section already exists
*/
func (this *Ini) AddSection(section string) bool {
	if this.SectionExists(section) {
		return false
	}
	var s Section
	s.items 	= make(map[string]Item)
	s.comments 	= make([]string,0)
	this.data[section] = s
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

		var tmp Item
			tmp.value = value
			tmp.comments = make([]string,0)
		this.data[section].items[item] = tmp
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

// Delete an item, return true if succes, false if the item does not exists
func (this *Ini) DeleteItem(section string,item string) bool {
	if this.Exists(section,item) {
		delete(this.data[section].items,item)
		return true
	}
	return false
}

// Delete a section, return true if succes, false if the section does not exists
func (this *Ini) DeleteSection(section string) bool {
	if this.SectionExists(section) {
		delete(this.data,section)
		return true
	}
	return false
}

/*
Save the ini format to a file

Example :

err := myIni.Save() // use myIni.Filename to save

err := myIni.Save("new_config.ini") // use new_config.ini and set myIni.Filename
*/
func (this *Ini) Save(params ...string) error {
	if len(params)>0 {
		this.Filename = params[0]
	}

	if this.Filename == "" {
		return errors.New("You must specify a filename before saving")
	}

	return ioutil.WriteFile(this.Filename,  []byte( this.Sprint() ), os.ModePerm)
}

/*
Return the ini format into a formatted string

TIPS :
You can set myIni.ItemPrefix, myIni.SectionSeparator and myIni.NoComments to tweak format aspect
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

TIPS :
You can set myIni.ItemPrefix, myIni.SectionSeparator and myIni.NoComments to tweak format aspect
*/
func (this *Ini) Print() {
	print(this.Sprint())
}

/* TO DO
GetSectionComments(section string) []strings
GetItemComments(section string,item string) []strings
DeleteSection(section string) bool
DeleteComment(id int) bool
DeleteComments() bool
*/