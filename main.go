package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type gdUploader struct {
	driveSvc  *drive.Service
	parents   []string
	command   string
	sharedUse bool
	authPath  string
	filter    string
	pageLimit int64
	files     []string
}

const (
	FolderMimeType       = "application/vnd.google-apps.folder"
	defaultWorkDirPrefix = ".gdrive_uploader"
	tokenName            = "token.json"
	credentialName       = "credential.json"
)

// Usage show simple manual
func Usage(code int) {
	fmt.Println("How to use gdrive_cli")
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("- show file list")
	fmt.Println("   $ gdrive_uploader list [OPTION]")
	fmt.Println("- upload file")
	fmt.Println("   $ gdrive_uploader upload [OPTION] file")
	fmt.Println("- create folder")
	fmt.Println("   $ gdrive_uploader mkdir [OPTION] name")
	fmt.Println("- when show usage")
	fmt.Println("   $ gdrive_uploader help")
	fmt.Println()
	fmt.Println("Global options")
	fmt.Println("  -c, --config [auth configfile path] : set path if want use specific auth files")
	fmt.Println("  -s, --shared : access to shared directories and files")
	fmt.Println()
	fmt.Println("Non global options")
	fmt.Println("- for mkdir|upload command")
	fmt.Println("  -p, --parents [parent ID1,parentID2,...] : set parent directory IDs seperate by comma(,)")
	fmt.Println("- for list command")
	fmt.Println("  -l, --limit [integer] : limit list file counts. The max count is 100(default). If set 0, use default.")
	fmt.Println("  -g, --grep [name] : search filename contain [name]")
	fmt.Println()

	os.Exit(code)
}

func (u *gdUploader) initParams() {
	args := os.Args[1:]
	if len(args) < 1 {
		os.Stderr.WriteString("Too few arguments!")
		Usage(1)
	}

	maxArgs := len(args)

	for i := 0; i < maxArgs; i++ {
		argv := args[i]

		switch argv {
		case "upload":
			if len(u.command) != 0 {
				os.Stderr.WriteString("Too many commands!")
				Usage(1)
			}
			u.command = "upload"
		case "mkdir":
			if len(u.command) != 0 {
				os.Stderr.WriteString("Too many commands!")
				Usage(1)
			}
			u.command = "mkdir"
		case "list":
			if len(u.command) != 0 {
				os.Stderr.WriteString("Too many commands!")
				Usage(1)
			}
			u.command = "list"
		case "--parents", "-p":
			if i+1 >= maxArgs {
				os.Stderr.WriteString("Too few arguments!")
				Usage(1)
			}
			if u.command == "list" {
				os.Stderr.WriteString("parants option cannot use with list command")
				Usage(1)
			}
			if len(u.files) > 0 {
				os.Stderr.WriteString("Options cannot set after target files")
				Usage(1)
			}
			parents := strings.Split(args[i+1], ",")
			u.parents = append(u.parents, parents...)

			i++
		case "--shared", "-s":
			u.sharedUse = true
			if len(u.files) > 0 {
				os.Stderr.WriteString("Options cannot set after target files")
				Usage(1)
			}
		case "--config", "-c":
			if len(u.authPath) != 0 {
				os.Stderr.WriteString("Too many token parameters!")
				Usage(1)
			}
			if i+1 >= maxArgs {
				os.Stderr.WriteString("Too few arguments!")
				Usage(1)
			}
			if len(u.files) > 0 {
				os.Stderr.WriteString("Options cannot set after target files")
				Usage(1)
			}
			u.authPath = args[i+1]
			i++
		case "--limit", "-l":
			if i+1 >= maxArgs {
				os.Stderr.WriteString("Too few arguments!")
				Usage(1)
			}
			if u.command != "list" {
				os.Stderr.WriteString("limit option can use with list command only")
				Usage(1)
			}
			if n, err := strconv.Atoi(args[i+1]); err != nil {
				os.Stderr.WriteString("limit parameter must be integer")
				Usage(1)
			} else {
				u.pageLimit = int64(n)
			}
			i++
		case "--grep", "-g":
			if i+1 >= maxArgs {
				os.Stderr.WriteString("Too few arguments!")
				Usage(1)
			}
			if u.command != "list" {
				os.Stderr.WriteString("grep option can use with list command only")
				Usage(1)
			}
			u.filter = args[i+1]
			i++
		case "help":
			Usage(0)
		default:
			u.files = append(u.files, argv)
		}
	}
}

func main() {
	var err error

	u := &gdUploader{
		sharedUse: false,
		authPath:  "",
		command:   "",
		pageLimit: 0,
	}

	u.initParams()

	ctx := context.Background()

	credential, err := getCredential(u.authPath)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("failed to get credential informations: %s", err.Error()))
		os.Exit(1)
	}

	config := getOauthConfig(credential.ClientID, credential.ClientSecret)
	client := getClient(u.authPath, config)

	u.driveSvc, err = drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Unable to retrieve Drive client: %s", err.Error()))
		os.Exit(1)
	}

	targetFiles := len(u.files)

	if u.command == "list" {
		if err := u.listFiles(); err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}

	for i, pathName := range u.files {
		switch u.command {
		case "mkdir":
			res, err := u.createFolder(pathName)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("[%d/%d] Failed to create Folder %s: %s\n", i+1, targetFiles, pathName, err.Error()))
			} else {
				fmt.Printf("[%d/%d] Folder created: (ID: %s Name: %s)\n", i+1, targetFiles, res.Id, res.Name)
			}
		case "upload":
			res, err := u.uploadFile(pathName)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("[%d/%d] Failed to upload File %s: %s\n", i+1, targetFiles, pathName, err.Error()))
			} else {
				fmt.Printf("[%d/%d] Upload finished: (ID: %s Name: %s)\n", i+1, targetFiles, res.Id, res.Name)
			}
		}
	}
}
