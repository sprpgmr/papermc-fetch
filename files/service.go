package files

import (
	"fmt"
	"os"
)

// Service contains method for determining if files exist on the system.
type Service interface {
	FileExists(filepath string) bool
	DeleteIfExists(filepath string) error
}

// GetFileService returns the default file service
func GetFileService() Service {
	return &serviceImpl{}
}

type serviceImpl struct{}

func (s *serviceImpl) FileExists(filepath string) bool {
	_, err := os.Stat(filepath)

	if err == nil {
		return true
	}

	if !os.IsNotExist(err) {
		fmt.Println("Couldn't check if file exists: ", err)
	}

	return false
}

func (s *serviceImpl) DeleteIfExists(filepath string) error {
	if s.FileExists(filepath) {
		err := os.Remove(filepath)
		return err
	}

	return nil
}
