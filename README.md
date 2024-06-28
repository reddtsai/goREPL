# REPL

這是個由 Golang 開發的命令列程式，讓使用者透過指令來操作檔案系統。

# 指令

## User Registration

### register

`register [username]`

| Parameter | Type   | Lenght | Desc                                                   |
| --------- | ------ | ------ | ------------------------------------------------------ |
| username  | string | 3-20   | case insensitive, can only contain letters and numbers |

| Response | Desc                           |
| -------- | ------------------------------ |
| Success  | Add [username] successfully    |
| Error    | invalid command                |
| Error    | The [%s] invalid length        |
| Error    | The [%s] contain invalid chars |
| Error    | The [%s] has already existed   |
