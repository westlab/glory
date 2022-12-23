package main

import (
	"fmt"
	"testing"
)

const (
	testDir   = "../test_dir"
	docxDir   = "../test_dir/thesis"
	noDocxDir = "../test_dir/no_docx_dir"
	newDocx   = "../test_dir/thesis/new.docx"
	oldDocx   = "../test_dir/thesis/old.docx"
)

func TestDirWalk(t *testing.T) {

	paths := DirWalk(testDir)

	if len(paths) != 4 {
		t.Errorf("DirWalk returns invalid slice: %v", paths)
	}

}

func TestFetchLatestDocx(t *testing.T) {

	// test in no docx dir
	fileName, err := FetchLatestDocx(noDocxDir)
	if err == NoDocxError {
	} else if err != nil {
		t.Errorf("FetchLatestDocx failed in '%s' : '%v'", noDocxDir, err)
	}

	if fileName != "" {
		t.Errorf("FetchLatestDocx returns incorrect file '%s' in '%s'", fileName, noDocxDir)
	}

	// test in docx dir
	fileName, err = FetchLatestDocx(docxDir)
	if err == NoDocxError {
		t.Errorf("FetchLatestDocx returns NoDocxError in docxDir: '%s'", docxDir)
	} else if err != nil {
		t.Errorf("FetchLatestDocx failed in '%s' : '%v'", docxDir, err)
	}

	if fileName != "../test_dir/thesis/new.docx" {
		t.Errorf("FetchLatestDocx returns incorrect docx in docxDir: '%s'", fileName)
	}

}

func TestCountCharsInDocx(t *testing.T) {

	var count int
	var err error

	//test with new docx
	count, _, err = CountCharsInDocx(newDocx)
	if err != nil {
		t.Errorf("CountCharsInDocx failed for %s", newDocx)
	}

	fmt.Printf("CountCharsInDocx returns '%d' in %s. This FAIL is not incorrect.", count, newDocx)

	// test with old docx
	count, _, err = CountCharsInDocx(oldDocx)
	if err != nil {
		t.Errorf("CountCharsInDocx failed for %s", oldDocx)
	}

	fmt.Printf("CountCharsInDocx returns '%d' in '%s'. This FAIL is not correct.", count, oldDocx)

}
