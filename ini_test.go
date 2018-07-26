package ini

import (
	"testing"
	"fmt"
	"sort"
	_ "github.com/davecgh/go-spew/spew"
)

func Test (t *testing.T) {
	//var result string = ""
	var a []string
	var s string	= ""
	var b bool 		= false

	myIni := new(Ini)
	content := `
# quelques comment en debut de fichier
 ; un autre commentaire en debut
	[section1]
# ce com est pour section1.item1
# ce 2eme com est pour section1.item1

		item1=value1

	# comment with a sharp

	item2=value2
; comment with a dot-comma
; second line of comments
 [ section2  ]
	item1=value1
	; comment 3
# second line of comments for comment 3
	item2=value2
`
	myIni.LoadFromString(&content)

	//myIni.SectionSeparator 	= "\r\n;-----------------------\r\n"
	//myIni.ItemSeparator 	= "\r\n"
	//myIni.SectionPrefix 	= " "
	//myIni.ItemPrefix 		= " "
	//myIni.ItemSuffix 		= " "
	//myIni.ValuePrefix 	= " "
	//myIni.CommentPrefix  	= "# "


	// GetSections
	a = myIni.GetSections() ; sort.Strings(a) ;	s = fmt.Sprintf("%v", a )
	exceptedSections := fmt.Sprintf("%v", []string{"section1","section2"} )
	if s != exceptedSections {
		t.Error("For", "GetSections()", "expected",	exceptedSections, "got", s )
	}

	// GetItems
	a = myIni.GetItems("section1") ; sort.Strings(a) ;	s = fmt.Sprintf("%v", a )
	exceptedItems := fmt.Sprintf("%v", []string{"item1","item2"} )
	if s != exceptedItems {
		t.Error("For", "GetItems(section1)", "expected", exceptedItems , "got", s )
	}

	// SectionExists
	if b = myIni.SectionExists("section1") ; !b {
		t.Error("For", "SectionExists(section1)", "expected", true , "got", b)
	}

	// SectionExists
	if b = myIni.SectionExists("does not exists") ; b {
		t.Error("For", "SectionExists(does not exists)", "expected", false , "got", b)
	}

	// Exists
	if b = myIni.Exists("section1","item1") ; !b {
		t.Error("For", "Exists(section1,item1)","expected", true , "got", b )
	}

	// Exists
	if b = myIni.Exists("section1","does not exists"); b {
		t.Error("For", "Exists(section1,does not exists)","expected", false , "got", b )
	}

	// Exists
	if b = myIni.Exists("does not exists","does not exists") ; b {
		t.Error("For", "Exists(does not exists,does not exists)","expected", false , "got", b )
	}

	// GET
	expectedValue := "value1"
	if s, _ = myIni.Get("section1","item1") ; s != expectedValue {
		t.Error("For", "Get(section1,item1)","expected", expectedValue, "got", s)
	}

	// SET
	if b = myIni.Set("section1","item2","edit value") ; !b {
		t.Error("For", "Set(section1,item2,edit value)","expected", true, "got", b)
	}
	// GET again to test if SET is OK
	expectedValue = "edit value"
	if s, _ = myIni.Get("section1","item2") ; s != expectedValue {
		t.Error("For", "Get(section1,item2)","expected", expectedValue, "got", s)
	}

	// RenameSection
	if b = myIni.RenameSection("section2","rename section") ; !b {
		t.Error("For", "RenameSection(section2,rename section)","expected", true, "got", b)
	}
	// SectionExists
	if b = myIni.SectionExists("rename section") ; !b {
		t.Error("For", "SectionExists(rename section)","expected", true, "got", b)
	}
	// SectionExists
	if b = myIni.SectionExists("section2") ; b {
		t.Error("For", "SectionExists(section2)","expected", false, "got", b)
	}

	// RenameItem
	if b = myIni.RenameItem("rename section","item1","renamed item") ; !b {
		t.Error("For", "RenameItem(rename section,item1,renamed item)","expected", true, "got", b)
	}
	// Exists
	if b = myIni.Exists("rename section","renamed item") ; !b {
		t.Error("For", "Exists(rename section,renamed item)","expected", true, "got", b)
	}
	// Exists
	if b = myIni.Exists("rename section","item1") ; b {
		t.Error("For", "Exists(rename section,item1)","expected", false, "got", b)
	}

	// AddSection
	if b = myIni.AddSection("added section") ; !b {
		t.Error("For", "AddSection(added section)","expected", true, "got", b)
	}
	// AddSection
	if b = myIni.AddSection("added section") ; b {
		t.Error("For", "AddSection(added section)","expected", false, "got", b)
	}
	// GetSections
	a = myIni.GetSections() ; sort.Strings(a) ;	s = fmt.Sprintf("%v", a )
	exceptedSections = fmt.Sprintf("%v", []string{"added section", "rename section", "section1"} )
	if s != exceptedSections {
		t.Error("For", "GetSections()", "expected",	exceptedSections, "got", s )
	}

	// AddItem
	if b = myIni.AddItem("added section", "added item", "added value") ; !b {
		t.Error("For", "AddItem(added section, added item, added value)","expected", true, "got", b)
	}
	if b = myIni.AddItem("added section", "added item", "added value") ; b {
		t.Error("For", "AddItem(added section, added item, added value)","expected", false, "got", b)
	}
	// GetItems
	a = myIni.GetItems("added section") ; sort.Strings(a) ;	s = fmt.Sprintf("%v", a )
	exceptedItems = fmt.Sprintf("%v", []string{"added item"} )
	if s != exceptedItems {
		t.Error("For", "GetItems(added section)", "expected", exceptedItems , "got", s )
	}
	// Get check if value is correct
	expectedValue = "added value"
	if s, _ = myIni.Get("added section","added item") ; s != expectedValue {
		t.Error("For", "Get(added section, added item)","expected", expectedValue, "got", s)
	}
	// AddItem to a existing section
	if b = myIni.AddItem("rename section", "added item", "add value") ; !b {
		t.Error("For", "AddItem(rename section, added item, add value)","expected", true, "got", b)
	}
	// AddItem to a non existing section
	if b = myIni.AddItem("another new section", "added item", "added value") ; !b {
		t.Error("For", "AddItem(another new section, added item, another value)","expected", true, "got", b)
	}

	// SetOrCreate Create a new section with new values
	myIni.SetOrCreate("new section", "new item", "new value")
	expectedValue = "new value"
	if s, _ = myIni.Get("new section","new item") ; s != expectedValue {
		t.Error("For", "Get(new section,new item)","expected", expectedValue,"got", s)
	}

	// SetOrCreate Set an existing value
	myIni.SetOrCreate("new section", "new item", "another value")
	expectedValue = "another value"
	if s, _ = myIni.Get("new section","new item") ; s != expectedValue {
		t.Error("For", "Get(new section,new item)","expected", expectedValue,"got", s)
	}

	// DeteleItem
	if b = myIni.DeleteItem("new section","new item") ; !b {
		t.Error("For", "DeleteItem(new section, new item)","expected", true,"got", b)
	}
	// Exists
	if b = myIni.Exists("new section", "new item") ; b {
		t.Error("For", "Exists(new section, new item)","expected", false, "got", b)
	}
	// DeteleItem try to redelete a not existing item
	if b = myIni.DeleteItem("new section","new item") ; b {
		t.Error("For", "DeleteItem(new section, new item)","expected", false,"got", b)
	}

	// DeleteSection
	if b = myIni.DeleteSection("new section") ; !b {
		t.Error("For", "DeleteSection(new section)","expected", true,"got", b)
	}
	// Exists
	if b = myIni.SectionExists("new section") ; b {
		t.Error("For", "SectionExists(new section)","expected", false, "got", b)
	}
	// DeleteSection try to redelete a not existing section
	if b = myIni.DeleteSection("new section") ; b {
		t.Error("For", "DeleteSection(new section)","expected", false,"got", b)
	}

	// GetSectionComments
	a = myIni.GetSectionComments("section1") ;	s = fmt.Sprintf("%v", a )
	exceptedComments := fmt.Sprintf("%v", []string{"quelques comment en debut de fichier", "un autre commentaire en debut"} )
	if s != exceptedComments {
		t.Error("For", "GetSectionComments(section1)", "expected", exceptedComments , "got", s )
	}

	// GetSectionComments with no comments
	a = myIni.GetSectionComments("another new section") ; s = fmt.Sprintf("%v", a )
	exceptedComments = fmt.Sprintf("%v", []string{} )
	if s != exceptedComments {
		t.Error("For", "GetSectionComments(another new section)", "expected", exceptedComments , "got", s )
	}

	// GetItemComments
	a = myIni.GetItemComments("section1","item1") ; s = fmt.Sprintf("%v", a )
	exceptedComments = fmt.Sprintf("%v", []string{"ce com est pour section1.item1", "ce 2eme com est pour section1.item1"} )
	if s != exceptedComments {
		t.Error("For", "GetItemComments(section1,item1)", "expected", exceptedComments , "got", s )
	}

	// GetItemComments with no comments
	a = myIni.GetItemComments("another new section","added item") ; s = fmt.Sprintf("%v", a )
	exceptedComments = fmt.Sprintf("%v", []string{} )
	if s != exceptedComments {
		t.Error("For", "GetItemComments(another new section,added item)", "expected", exceptedComments , "got", s )
	}

	//myIni.Print()
}