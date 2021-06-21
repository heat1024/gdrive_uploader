# gdrive_uploader
upload to gdrive include shared drive

How to use gdrive_cli
--------------------------------------------------------------------
- show file list
   $ gdrive_uploader list [OPTION]
- upload file
   $ gdrive_uploader upload [OPTION] file
- create folder
   $ gdrive_uploader mkdir [OPTION] name
- when show usage
   $ gdrive_uploader help

- Global options
  -c, --config [auth configfile path] : set path if want use specific auth files
  -s, --shared : access to shared directories and files

Non global options
- for mkdir|upload command
  -p, --parents [parent ID1,parentID2,...] : set parent directory IDs seperate by comma(,)
- for list command
  -l, --limit [integer] : limit list file counts. The max count is 100(default). If set 0, use default.
  -g, --grep [name] : search filename contain [name]
