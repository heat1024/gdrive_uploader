package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// list files
func (u *gdUploader) listFiles() error {
	var list *drive.FileList
	var err error

	query := fmt.Sprintf("name contains '%s'", u.filter)
	fields := []googleapi.Field{"nextPageToken, files(id, name, md5Checksum,mimeType,parents)"}

	if u.pageLimit == 0 {
		list, err = u.driveSvc.Files.List().SupportsAllDrives(u.sharedUse).Q(query).Fields(fields...).Do()
	} else {
		list, err = u.driveSvc.Files.List().SupportsAllDrives(u.sharedUse).Q(query).Fields(fields...).PageSize(u.pageLimit).Do()
	}
	if err != nil {
		return err
	}

	if len(list.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		fmt.Printf("%-50s %-10s %s\n", "ID", "Type", "Name")
		for _, file := range list.Files {
			fmt.Printf("%-50s %-10s %s\n", file.Id, getType(file), file.Name)
		}
	}
	return nil
}

// get file type
func getType(f *drive.File) string {
	if f.MimeType == FolderMimeType {
		return "dir"
	} else if len(f.Md5Checksum) != 0 {
		return "bin"
	}

	return "doc"
}

// create folder
func (u *gdUploader) createFolder(dir string) (*drive.File, error) {
	if u == nil {
		return nil, errors.New("drive service not defined")
	}

	d := &drive.File{
		Name:     dir,
		Parents:  u.parents,
		MimeType: FolderMimeType,
	}

	return u.driveSvc.Files.Create(d).SupportsAllDrives(u.sharedUse).Do()
}

// upload file
func (u *gdUploader) uploadFile(file string) (*drive.File, error) {
	if u == nil {
		return nil, errors.New("drive service not defined")
	}

	// get clean filename without path
	cleanName := filepath.Base(file)

	fd, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %s", err.Error())
	}

	f := &drive.File{
		Name:    cleanName,
		Parents: u.parents,
	}

	return u.driveSvc.Files.Create(f).SupportsAllDrives(u.sharedUse).Media(fd).Do()
}
