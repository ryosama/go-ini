package ini

import "testing"

func Test (t *testing.T) {
	result 	:= ""
	success := false
	myini := new(Ini)
	content := `
	[section1]
		item1=value1
	# comment with a sharp
	item2=value2
; comment with a dot-comma
 [ section2  ]
	item1=value1
	; comment 3
	item2=value2
`
	myini.LoadFromString(&content)

	// GetSections
	if len(myini.GetSections()) != 2 {
		t.Error("For", "GetSections()","expected", 2 ,"got", len(myini.GetSections()) )
	}

	// GetItems
	if len(myini.GetItems("section1")) != 2 {
		t.Error("For", "GetItems(section1)","expected", 2 ,"got", len(myini.GetItems("section1")) )
	}

	// SectionExists
	if myini.SectionExists("section1") == false {
		t.Error("For", "SectionExists(section1)","expected", true ,"got", myini.SectionExists("section1") )
	}

	// SectionExists
	if myini.SectionExists("does not exists") == true {
		t.Error("For", "SectionExists(does not exists)","expected", false ,"got", myini.SectionExists("does not exists") )
	}

	// Exists
	if myini.Exists("section1","item1") == false {
		t.Error("For", "Exists(section1,item1)","expected", true ,"got", myini.Exists("section1","item1") )
	}

	// Exists
	if myini.Exists("section1","does not exists") == true {
		t.Error("For", "Exists(section1,does not exists)","expected", false ,"got", myini.Exists("section1","does not exists") )
	}

	// Exists
	if myini.Exists("does not exists","does not exists") == true {
		t.Error("For", "Exists(does not exists,does not exists)","expected", false ,"got", myini.Exists("does not exists","does not exists") )
	}

	// GET
	result, _ 	= myini.Get("section1","item1")
	if result != "value1" {
		t.Error("For", "Get(section1,item1)","expected", "value1","got", result)
	}

	// SET
	success 	= myini.Set("section1","item2","edit value")
	if !success {
		t.Error("For", "Set(section1,item2,edit value)","expected", true,"got", false)
	}
	// GET again
	result, _ 	= myini.Get("section1","item2")
	if result != "edit value" {
		t.Error("For", "Get(section1,item2)","expected", "edit value", "got", result)
	}

	// RenameSection
	success 	= myini.RenameSection("section2","section3")
	if !success {
		t.Error("For", "RenameSection(section2,section3)","expected", true, "got", success)
	}
	// SectionExists
	success 	= myini.SectionExists("section3")
	if !success {
		t.Error("For", "SectionExists(section3)","expected", true, "got", success)
	}

	// RenameItem
	success 	= myini.RenameItem("section3","item1","edit item")
	if !success {
		t.Error("For", "RenameItem(section3,item1,edit item)","expected", true, "got", success)
	}
	// Exists
	success 	= myini.Exists("section3","edit item")
	if !success {
		t.Error("For", "Exists(section3,edit item)","expected", true, "got", success)
	}

	// AddSection
	success 	= myini.AddSection("added section")
	if !success {
		t.Error("For", "AddSection(added section)","expected", true, "got", success)
	}
	// AddSection
	success 	= myini.AddSection("added section")
	if success != false {
		t.Error("For", "AddSection(added section)","expected", false, "got", success)
	}

	// AddItem
	success 	= myini.AddItem("added section", "added item", "add value")
	if !success {
		t.Error("For", "AddSection(added section, added item, add value)","expected", true, "got", success)
	}
	// AddItem
	success 	= myini.AddItem("added section", "added item", "add value")
	if success != false {
		t.Error("For", "AddSection(added section, added item, add value)","expected", false, "got", success)
	}

	// SetOrCreate
	myini.SetOrCreate("new section", "new item", "new value")
	result, _ 	= myini.Get("new section","new item")
	if result != "new value" {
		t.Error("For", "Get(new section,new item)","expected", "new value","got", result)
	}

	// SetOrCreate
	myini.SetOrCreate("new section", "new item", "another value")
	result, _ 	= myini.Get("new section","new item")
	if result != "another value" {
		t.Error("For", "Get(new section,new item)","expected", "another value","got", result)
	}
}