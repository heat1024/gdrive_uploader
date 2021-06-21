# gdrive_uploader

upload to gdrive include shared drive

## How to get

```
go get github.com/heat1024/gdrive_uploader
```

## About credential

You must set your client_id and client_secret to `credential.json` file.
Default authenticate config directory is `${HOME}/.gdrive_uploader/`
There is refer about create client_id and client_secret.
visit [Create credentials](https://developers.google.com/workspace/guides/create-credentials) and create desktop app OAuth 2.0 Client ID.

If you use config option, you can specific your own auth directory path.

When you have no token file `token.json` in your auth directory, you must exchange token follow the URL in your browser.

```Shell
$ cat test_creds/credential.json
{
    "client_id": "xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.com",
    "client_secret": "some_secrets"
}

$ ./gdrive_uploader list -c ./test_creds
Go to the following link in your browser then type the authorization code:
https://accounts.google.com/o/oauth2/auth?access_type=offline&xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.com&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive&state=state-token
{ENTER AUTH CODE HERE}
Saving credential file to: ./test_creds/token.json
ID                                                 Type       Name
1S2MDtDx1C_QhIpUzpbId3ehSH8vb6JnZ                  bin        car.jpg
1X87w93ji-m47ixgoZRW4u2aFpbtvVbqw                  dir        gdrive_test
1e8kOo3r51b2BWtTs_1uADIA5djfXhPT36s6eHVRIvaU       doc        Go 1.4 "Internal" Packages

$ ls test_creds
credential.json token.json
```

## Usage

```Shell
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
```

## Example


### Get file list with grep option

```Shell
$ ./gdrive_uploader list -c ./test_creds --grep gdrive
ID                                                 Type       Name
1X87w93ji-m47ixgoZRW4u2aFpbtvVbqw                  dir        gdrive_test
```

### Make folder

### Make one folder and multi folders

```Shell
$ ./gdrive_uploader mkdir -c ./test_creds example
[1/1] Folder created: (ID: 19ilu7lA_ivLgH1sNil0VMrhtOYXCbaOu Name: example)

$ ./gdrive_uploader mkdir -c ./test_creds example1 example2 example3
[1/3] Folder created: (ID: 1nhHxfp_F7xqT6gq7XewUWcC1bpaSXsq1 Name: example1)
[2/3] Folder created: (ID: 1s8O449fI1kxAW8EASpBuIDh5TtMXj6Wj Name: example2)
[3/3] Folder created: (ID: 1ZzIx-__XQuKlUw2timJvd8wKLLX3fz8R Name: example3)
```

### Make folder under the other folder

```Shell
$ ./gdrive_uploader mkdir -c ./test_creds --parents 1nhHxfp_F7xqT6gq7XewUWcC1bpaSXsq1 sub_folder
[1/1] Folder created: (ID: 1XFx7PS4ZrQcfXJzcYbYOei0LpTJpqank Name: sub_folder)
```

## Upload file

### Upload multi files under the folder

```Shell
❯ ./gdrive_uploader upload -c ./test_creds --parents 1XFx7PS4ZrQcfXJzcYbYOei0LpTJpqank ./gdrive_uploader test_creds/credential.json test_creds/token.json
[1/3] Upload finished: (ID: 1lnhjOIWZo_gMflF9bDrzdeDO7OGsqh8W Name: gdrive_uploader)
[2/3] Upload finished: (ID: 1JFOjEtAZg9sQTYE7SAdjBYmwfgdvsZtK Name: credential.json)
[3/3] Upload finished: (ID: 1OjkJWCrGrmq8-ifBsVucq-69ytbKRmp4 Name: token.json)

❯ ./gdrive_uploader list -c ./test_creds --grep json
ID                                                 Type       Name
1OjkJWCrGrmq8-ifBsVucq-69ytbKRmp4                  bin        token.json
1JFOjEtAZg9sQTYE7SAdjBYmwfgdvsZtK                  bin        credential.json
```

## Require

- Go 1.15 or later

## Author

@heat1024
