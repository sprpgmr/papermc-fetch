package files

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	testFileName := ".test"

	cleanupTestFile(testFileName)

	fileService := GetFileService()
	exists := fileService.FileExists(testFileName)
	if exists {
		t.Error("expected file to not exist!")
	}

	err := createTestFile(testFileName)
	if err != nil {
		t.Error(err)
	}

	defer cleanupTestFile(testFileName)

	exists = fileService.FileExists(testFileName)
	if !exists {
		t.Error("expected file to exist!")
	}
}

func TestDeleteIfExists(t *testing.T) {
	testFileName := ".test"

	fileService := GetFileService()

	err := createTestFile(testFileName)
	if err != nil {
		t.Error(err)
	}

	defer cleanupTestFile(testFileName)

	err = fileService.DeleteIfExists(testFileName)
	if err != nil {
		t.Error(err)
	}

	if _, err = os.Stat(testFileName); err == nil {
		t.Error("File shouldn't exist after trying to delete!!")
	}
}

func createTestFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func cleanupTestFile(filename string) {
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}
}
