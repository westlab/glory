package main

import (
	"testing"
)

const (
	testDir   = "../test_dir"
	docxDir   = "../test_dir/thesis"
	noDocxDir = "../test_dir/no_docx_dir"
	fiveDocx  = "../test_dir/thesis/new.docx"
	anninDocx = "../test_dir/thesis/old.docx"
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
		t.Errorf("FetchLatestDocx returns NoDocxError in '%s'", docxDir)
	} else if err != nil {
		t.Errorf("FetchLatestDocx failed in '%s' : '%v'", docxDir, err)
	}

	if fileName != "../test_dir/thesis/new.docx" {
		t.Errorf("FetchLatestDocx returns incorrect docx '%s'", fileName)
	}

}

func TestCountCharsInDocx(t *testing.T) {

	// test with new docx
	count, _, err := CountCharsInDocx(fiveDocx)
	if err != nil {
		t.Errorf("CountCharsInDocx failed for %s", fiveDocx)
	}

	if count <= 0 {
		t.Errorf("CountCharsInDocx returns incorrct value '%d'", count)
	}

	// test with old docx
	count, _, err = CountCharsInDocx(anninDocx)
	if err != nil {
		t.Errorf("CountCharsInDocx failed for %s", anninDocx)
	}

	if count <= 0 {
		t.Errorf("CountCharsInDocx returns incorrct value '%d' in '%s'", count, anninDocx)
	}

}
