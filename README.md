# REPL

This is a command-line program developed by Golang, allowing users to operate the file system through commands.

[![Build Status][ci-badge]][ci-runs]

# Installation

`go install github.com/reddtsai/goREPL`

# Commands

## User Management

### Register

`register [username]`

Register a user.

| Parameter | Type   | Lenght | Desc                                                    |
| --------- | ------ | ------ | ------------------------------------------------------- |
| username  | string | 3 - 20 | case insensitive, can only contain letters and numbers. |

| Response | Content                              |
| -------- | ------------------------------------ |
| Success  | add [username] successfully          |
| Error    | unrecognized argument                |
| Error    | the [username] invalid length        |
| Error    | the [username] contain invalid chars |
| Error    | the [username] has already existed   |

## Folder Management

### Create Folder

`create-folder [username] [foldername] [description]?`

Create a folder for a user.

| Parameter   | Type   | Lenght  | Desc                                                                                                                                                                                                               |
| ----------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| username    | string | 3 - 20  | case insensitive, contain only contain letters and numbers.                                                                                                                                                        |
| foldername  | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |
| description | string | 500     | you can use either double quotes (“”) or single quotes(‘’).                                                                                                                                                        |

| Response | Content                                |
| -------- | -------------------------------------- |
| Success  | create [foldername] successfully       |
| Error    | unrecognized argument                  |
| Error    | the [username] doesn't exist           |
| Error    | the [foldername] invalid length        |
| Error    | the [foldername] contain invalid chars |
| Error    | the [foldername] has already existed   |
| Error    | the [description] invalid length       |

### Delete Folder

`delete-folder [username] [foldername]`

Delete a folder for a user.

| Parameter  | Type   | Lenght  | Desc                                                                                                                                                                                                               |
| ---------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| username   | string | 3 - 20  | case insensitive, contain only contain letters and numbers.                                                                                                                                                        |
| foldername | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |

| Response | Content                          |
| -------- | -------------------------------- |
| Success  | delete [foldername] successfully |
| Error    | unrecognized argument            |
| Error    | the [username] doesn't exist     |
| Error    | the [foldername] doesn't exist   |

### List Folders

`list-folders [username] [--sort-name|--sort-created] [asc|desc]`

List user folders.

| Parameter | Type   | Lenght | Desc                                                    |
| --------- | ------ | ------ | ------------------------------------------------------- |
| username  | string | 3 - 20 | case insensitive, can only contain letters and numbers. |

| Option         | Argument  | Memo                                |
| -------------- | --------- | ----------------------------------- |
| --sort-name    | asc, desc | `--sort-name asc` is defaule option |
| --sort-created | asc, desc |                                     |

| Response | Content                                          |
| -------- | ------------------------------------------------ |
| Success  | List {foldername description createat username } |
| Warning  | the [username] doesn't have any folders          |
| Error    | unrecognized argument                            |
| Error    | the [username] doesn't exist                     |

### Rename Folder

`rename-folder [username] [foldername] [newfoldername]`

Rename a folder for a user.

| Parameter     | Type   | Lenght  | Desc                                                                                                                                                                                                               |
| ------------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| username      | string | 3 - 20  | case insensitive, contain only contain letters and numbers.                                                                                                                                                        |
| foldername    | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |
| newfoldername | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |

| Response | Content                                             |
| -------- | --------------------------------------------------- |
| Success  | rename [foldername] to [newfoldername] successfully |
| Error    | unrecognized argument                               |
| Error    | the [username] doesn't exist                        |
| Error    | the [foldername] doesn't exist                      |
| Error    | the [newfoldername] invalid length                  |
| Error    | the [newfoldername] contain invalid chars           |
| Error    | the [newfoldername] has already existed             |

## File Management

### Create File

`create-file [username] [foldername] [filename] [description]?`

Create a file for a user.

| Parameter   | Type   | Lenght  | Desc                                                                                                                                                                                                               |
| ----------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| username    | string | 3 - 20  | case insensitive, contain only contain letters and numbers.                                                                                                                                                        |
| foldername  | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |
| filename    | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |
| description | string | 500     | you can use either double quotes (“”) or single quotes(‘’).                                                                                                                                                        |

| Response | Content                                                   |
| -------- | --------------------------------------------------------- |
| Success  | create [filename] in [username]/[foldername] successfully |
| Error    | unrecognized argument                                     |
| Error    | the [username] doesn't exist                              |
| Error    | the [foldername] doesn't exist                            |
| Error    | the [filename] invalid length                             |
| Error    | the [filename] contain invalid chars                      |
| Error    | the [filename] has already existed                        |
| Error    | the [description] invalid length                          |

### Delete File

`delete-file [username] [foldername] [filename]`

Delete a file for a user.

| Parameter  | Type   | Lenght  | Desc                                                                                                                                                                                                               |
| ---------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| username   | string | 3 - 20  | case insensitive, contain only contain letters and numbers.                                                                                                                                                        |
| foldername | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |
| filename   | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |

| Response | Content                                                   |
| -------- | --------------------------------------------------------- |
| Success  | delete [filename] in [username]/[foldername] successfully |
| Error    | unrecognized argument                                     |
| Error    | the [username] doesn't exist                              |
| Error    | the [foldername] doesn't exist                            |
| Error    | the [filename] doesn't exist                              |

### List Files

`list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]`

| Parameter  | Type   | Lenght  | Desc                                                                                                                                                                                                               |
| ---------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| username   | string | 3 - 20  | case insensitive, can only contain letters and numbers.                                                                                                                                                            |
| foldername | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |

| Option         | Argument  | Memo                                |
| -------------- | --------- | ----------------------------------- |
| --sort-name    | asc, desc | `--sort-name asc` is defaule option |
| --sort-created | asc, desc |                                     |

| Response | Content                                                   |
| -------- | --------------------------------------------------------- |
| Success  | List {filename description createat foldername username } |
| Warning  | the [foldername] is empty                                 |
| Error    | unrecognized argument                                     |
| Error    | the [username] doesn't exist                              |
| Error    | the [foldername] doesn't exist                            |

## Help

`help`

List all command.

## Exit

`exit`

Close command prompt.
