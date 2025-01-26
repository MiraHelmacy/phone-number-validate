package config

import (
	"io"
	"os"
)

const mAX_READ_SIZE = 1024 // 1 KB

// read a files contents
func ReadFileContents(filePath string) (contents string, err error) {
	//attempts to open the config file
	if file, openFileErr := openFile(filePath); openFileErr == nil {
		//close file after finished
		defer file.Close()

		//full bytes of file
		fullBytes := make([]byte, 0)

		//read file in chunks at a time
		bytes := make([]byte, mAX_READ_SIZE)

		//n bytes read
		var n int = -1

		//read error detected
		var readErr error = nil

		//read file in chunks until an error is encountered or EOF reached.
		for n, readErr = file.Read(bytes); n > 0 && readErr == nil; n, readErr = file.Read(bytes) {

			//append n bytes to fullBytes array
			fullBytes = append(fullBytes, bytes[0:n]...)

			//recreate bytes array
			bytes = make([]byte, mAX_READ_SIZE)
		}

		//check if an error occurred and the error is not end of file
		if readErr != nil && readErr != io.EOF {

			//something other than EOF error occurred
			err = readErr
		} else {
			//convert fullBytes to string
			contents = string(fullBytes)
		}

	} else {
		//failed to open file
		err = openFileErr
	}
	return
}

// open a file
func openFile(filePath string) (file *os.File, err error) {

	//open a file as read only
	file, err = os.OpenFile(filePath, os.O_RDONLY, 0644)

	return
}
