package fileutils

import (
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from source to dest and returns
// an error if any.
func CopyFile(source string, dest string) error {
	// Open the source file.
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	// Makes the directory needed to create the dst
	// file.
	err = os.MkdirAll(filepath.Dir(dest), 0766)
	if err != nil {
		return err
	}

	// Create the destination file.
	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the contents of the file.
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	// Copy the mode if the user can't
	// open the file.
	info, err := os.Stat(source)
	if err != nil {
		err = os.Chmod(dest, info.Mode())
		if err != nil {
			return err
		}
	}

	return nil
}
