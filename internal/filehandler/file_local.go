package filehandler

import (
	"errors"
	"io/fs"
	"os"
)

type FileLocal struct {
}

func (f *FileLocal) DoesFileExist(filepath string) bool {
	// https://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-exists
	_, err := os.Stat(filepath)
	if err == nil {
		return true
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return false
}

func (f *FileLocal) CreateFolder(folderpath string) error {
	if _, err := os.Stat(folderpath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(folderpath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileLocal) WriteFile(filepath string, body []byte) error {
	return os.WriteFile(filepath, []byte(body), 0644)
}

func (f *FileLocal) ReadFile(filepath string) (*[]byte, error) {
	v, err := os.ReadFile(filepath)
	return &v, err
}

func (f *FileLocal) MoveFile(original_path string, new_path string) error {
	return os.Rename(original_path, new_path)
}

func (f *FileLocal) FolderScan(folderpath string, glob string) (*[]string, error) {
	root := os.DirFS(folderpath)
	stlFiles, err := fs.Glob(root, glob)
	if err != nil {
		return nil, err
	}
	return &stlFiles, nil
}
