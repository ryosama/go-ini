Go-ini
======

Go library to deals with ini configuration files

Install
=======

```bash
$ go get -u github.com/ryosama/go-ini
```

Quick Start
===========

```Go
// Load the library
import "github.com/ryosama/go-ini"

// Create object
myIni := new(ini.Ini)

// Load config file
if err := myIni.LoadFromFile("config.ini") ; err != nil {
	panic("Unable to load configuration : " + err.Error())
}

// Get the value
myHost, _ := myIni.Get("Server","host")

// Set another value and save
myIni.Set("Server","host","127.0.0.1")
myIni.Save()
```

Documentation
=======
```bash
$ godoc github.com/ryosama/go-ini
```
Creating object
	
    myIni := new(ini.Ini)
--------------------
Object properties

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
    // contains filtered or unexported fields
--------------------
- AddItem

Add an item. Returns true if success, false if item already exists

	func (this *Ini) AddItem(section string, item string, value string) bool
--------------------
- AddItemComment

Add a comment to an item, return true if succeed, false otherwise

	func (this *Ini) AddItemComment(section string, item string, comment string) bool
--------------------
- AddSection

Add a section. Returns true if success, false if section already exists
    
	func (this *Ini) AddSection(section string) bool
--------------------
- AddSectionComment

Add a comment to a setion, return true if succeed, false otherwise

	func (this *Ini) AddSectionComment(section string, comment string) bool
--------------------
- DeleteItem

Delete an item, return true if succes, false if the item does not exists

	func (this *Ini) DeleteItem(section string, item string) bool
--------------------
- DeleteItemComment

Delete the comment number id, return true if succeed, false otherwise

	func (this *Ini) DeleteItemComment(section string, item string, id int) bool
--------------------  
- DeleteItemComments

Delete all the comments of an item, return true if succeed, false otherwise

	func (this *Ini) DeleteItemComments(section string, item string) bool
--------------------
- DeleteSection

Delete a section, return true if succes, false if the section does not exists

	func (this *Ini) DeleteSection(section string) bool
--------------------
- DeleteSectionComment

Delete the comment number id, return true if succeed, false otherwise

	func (this *Ini) DeleteSectionComment(section string, id int) bool
--------------------
- DeleteSectionComments

Delete all the comments of a section, return true if succeed, false
    otherwise

	func (this *Ini) DeleteSectionComments(section string) bool
    
--------------------
- Exists

Alias for ItemExists

	func (this *Ini) Exists(section string, item string) bool
--------------------
- Get

Alias for GetItem

	func (this *Ini) Get(section string, item string) (string, bool)
--------------------
- GetItem

Returns the items value of the ini file for a given section and item
    (and true as second return value)

If the item does not exists, return false as second return value

	func (this *Ini) GetItem(section string, item string) (string, bool)

Example :

    value, success := myini.GetItem("section1","item1")
--------------------
- GetItemComments

Return the comments, one per line, just before the item, return empty slice of string if item does not exists.

	func (this *Ini) GetItemComments(section string, item string) []string

Example :

    for _, com := range myIni.GetItemComments("mySection","myItem") {
		print("Comment for myItem", com, "\n")
    }
--------------------
- GetItems

Returns all the items of the ini file for a given section

	func (this *Ini) GetItems(section string) []string
--------------------
- GetSectionComments


Return the comments, one per line, just before the section, return empty slice of string if section does not exists.

	func (this *Ini) GetSectionComments(section string) []string
    
Example :

    for _, com := range sectionExists := myIni.GetSectionComments("mySection") {
		print("Comment for mySection", com, "\n")
    }
--------------------
- GetSections

Returns all the sections of the ini file

	func (this *Ini) GetSections() []string
--------------------
- ItemExists

Returns true or false if the item exists for the given section

	func (this *Ini) ItemExists(section string, item string) bool
--------------------
- LoadFromFile

Read ini format from a file

	func (this *Ini) LoadFromFile(filename string) error
    
Example :

    err := myIni.LoadFromFile("config.ini")
--------------------
- LoadFromString

Load data from a string pointer

	func (this *Ini) LoadFromString(content *string)

Example :   

    content := `

    [section1]

    item1=value1`

    myini.LoadFromString( &content )
--------------------
- Print

Print the ini format into a formatted string

TIPS : You can set _SectionPrefix,ItemPrefix, ItemSuffix, ValuePrefix, SectionSeparator, ItemSeparator, WithComments, CommentPrefix_ to tweak format aspect

	func (this *Ini) Print()
--------------------
- RenameItem

Rename an item. Returns true if success, false if section or item does not exists

	func (this *Ini) RenameItem(section, oldName string, newName string) bool
--------------------
- RenameSection

Rename a section. Returns true if success, false if section does not exists

	func (this *Ini) RenameSection(oldName string, newName string) bool
--------------------
- Save

Save the ini format to a file

	func (this *Ini) Save(params ...string) error
    
Example :

	err := myIni.Save() // use myIni.Filename to save
    err := myIni.Save("new_config.ini") // use new_config.ini and set
    myIni.Filename
--------------------
- SectionExists

Returns true or false if the section exists

	func (this *Ini) SectionExists(section string) bool
--------------------
- Set

Alias for SetItem	

	func (this *Ini) Set(section string, item string, value string) bool
--------------------
- SetItem

Set the value of an item. Returns true if success, false section or item does not exists

	func (this *Ini) SetItem(section string, item string, value string) bool
--------------------
- SetOrCreate

Set a value for an item, create section and item if needed

	func (this *Ini) SetOrCreate(section string, item string, value string)
--------------------
- Sprint

Return the ini format into a formatted string

TIPS : You can set _SectionPrefix,ItemPrefix, ItemSuffix, ValuePrefix, SectionSeparator, ItemSeparator, WithComments, CommentPrefix_ to tweak format aspect

	func (this *Ini) Sprint() string
    
