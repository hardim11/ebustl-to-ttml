package filehandler

// the idea is to move all file operations to here so we can support remote URLs etc
// for now, assume file!

var file_handler = FileLocal{}

func DoesFileExist(filepath string) bool {
	return file_handler.DoesFileExist(filepath)
}

func CreateFolder(folderpath string) error {
	return file_handler.CreateFolder(folderpath)
}

func WriteFile(filepath string, body []byte) error {
	return file_handler.WriteFile(filepath, body)
}

func ReadFile(filepath string) (*[]byte, error) {
	return file_handler.ReadFile(filepath)
}

func MoveFile(original_path string, new_path string) error {
	return file_handler.MoveFile(original_path, new_path)
}

func FolderScan(folderpath string, glob string) (*[]string, error) {
	return file_handler.FolderScan(folderpath, glob)
}
