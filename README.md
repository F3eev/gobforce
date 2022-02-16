

## introduction

结合rustscan实现快速扫描爆破


## setup

1，安装nmap,rustscan
2，运行scan.sh开始扫描并保持扫描结果Xml
3，运行Goscanpro读取XMl 并进行扫描

### start

运行脚本使用rustscan扫描./ip.txt文件中ip，out目录为nmap输出结果目录

```
./scan.sh

```
扫描服务弱口令

```
扫面目标
go run main.go -nFile out/123.56.102.89.xml
go run main.go -nDir out
使用自定义字典
go run main.go -nDir out -CustomDict
```

- [x] FTP
- [x] mongodb
- [x] mysql
- [x] postgres
- [x] rdp
- [x] redis
- [x] ssh
- [x] vnc


## to do
- rdp优化 

## update log

### 2022.2.11
int
### 2022.2.16
update


