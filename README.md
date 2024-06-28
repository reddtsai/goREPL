# REPL

This is a command-line program developed by Golang, allowing users to operate the file system through commands.

# Commands

## User Management

### Register

`register [username]`

| Parameter | Type   | Lenght | Desc                                                   |
| --------- | ------ | ------ | ------------------------------------------------------ |
| username  | string | 3 - 20 | case insensitive, can only contain letters and numbers |

| Response | Desc                                 |
| -------- | ------------------------------------ |
| Success  | Add [username] successfully          |
| Error    | invalid command                      |
| Error    | the [username] invalid length        |
| Error    | the [username] contain invalid chars |
| Error    | the [username] has already existed   |

## Folder Management

### Create Folder

`create-folder [username] [foldername] [description]?`

| Parameter   | Type   | Lenght  | Desc                                                                                                                                                                                                               |
| ----------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| username    | string | 3 - 20  | case insensitive, contain only contain letters and numbers                                                                                                                                                         |
| foldername  | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |
| description | string | 500     |                                                                                                                                                                                                                    |

| Response | Desc                                   |
| -------- | -------------------------------------- |
| Success  | Create [foldername] successfully       |
| Error    | invalid command                        |
| Error    | the [username] doesn't exist           |
| Error    | the [foldername] invalid length        |
| Error    | the [foldername] contain invalid chars |
| Error    | the [foldername] has already existed   |
| Error    | the [description] invalid length       |

### Delete Folder

`delete-folder [username] [foldername]`

| Parameter  | Type   | Lenght  | Desc                                                                                                                                                                                                               |
| ---------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| username   | string | 3 - 20  | case insensitive, contain only contain letters and numbers                                                                                                                                                         |
| foldername | string | 1 - 100 | case insensitive, contain only the following characters: uppercase letters (A-Z), lowercase letters (a-z), numbers (0-9), periods (.), hyphens (-), tildes (~), underscores (\_), equal signs (=), and colons (:). |

| Response | Desc                             |
| -------- | -------------------------------- |
| Success  | Delete [foldername] successfully |
| Error    | invalid command                  |
| Error    | the [username] doesn't exist     |
| Error    | the [foldername] doesn't exist   |

### List Folders

`list-folders [username] [--sort-name|--sort-created] [asc|desc]`

| Parameter      | Type   | Lenght | Desc                                                   |
| -------------- | ------ | ------ | ------------------------------------------------------ |
| username       | string | 3 - 20 | case insensitive, can only contain letters and numbers |
| --sort-name    | string | 5      | asc or desc                                            |
| --sort-created | string | 5      | asc or desc                                            |

| Response | Desc                                             |
| -------- | ------------------------------------------------ |
| Success  | List {foldername description createat username } |
| Error    | invalid command                                  |
| Error    | the [username] doesn't exist                     |
| Error    | unknown flag                                     |
| Error    | the [asc desc] invalid                           |
