package main

import (
	"log"
	"os"
)

/*
目录操作
文件操作的大多数函数都是在os包里面，下面列举了几个目录操作的：

func Mkdir(name string, perm FileMode) error
创建名称为name的目录，权限设置是perm，例如0777

func MkdirAll(path string, perm FileMode) error
根据path创建多级子目录，例如astaxie/test1/test2。

func Remove(name string) error
删除名称为name的目录，当目录下有文件或者其他目录时会出错

func RemoveAll(path string) error
根据path删除多级子目录，如果path是单个名称，那么该目录下的子目录全部删除。
*/
func dir() {
	os.Mkdir("Zou", 0777)
	os.MkdirAll("Zou/Hello/world/nihao", 0777)
	err := os.Remove("Zou")
	if err != nil {
		log.Println(err)
		os.RemoveAll("Zou")
	}
}

/*
 文件操作
建立与打开文件
新建文件可以通过如下两个方法

func Create(name string) (file *File, err Error)
根据提供的文件名创建新的文件，返回一个文件对象，默认权限是0666的文件，返回的文件对象是可读写的。

func NewFile(fd uintptr, name string) *File
根据文件描述符创建相应的文件，返回一个文件对象

通过如下两个方法来打开文件：

func Open(name string) (file *File, err Error)
该方法打开一个名称为name的文件，但是是只读方式，内部实现其实调用了OpenFile。

func OpenFile(name string, flag int, perm uint32) (file *File, err Error)
打开名称为name的文件，flag是打开的方式，只读、读写等，perm是权限
*/
/*
写文件
写文件函数：

func (file *File) Write(b []byte) (n int, err Error)
写入byte类型的信息到文件

func (file *File) WriteAt(b []byte, off int64) (n int, err Error)
在指定位置开始写入byte类型的信息

func (file *File) WriteString(s string) (ret int, err Error)
写入string信息到文件
*/
func writeFile() {
	userFile := "zou.txt"
	fout, err := os.Create(userFile)
	if err != nil {
		log.Println(userFile, err)
		return
	}
	defer fout.Close()
	for i := 0; i < 10; i++ {
		fout.WriteString("Just a testb!\r\n")
		fout.Write([]byte("Just a testb!\r\n"))
	}

}

/*
 读文件
读文件函数：

func (file *File) Read(b []byte) (n int, err Error)
读取数据到b中

func (file *File) ReadAt(b []byte, off int64) (n int, err Error)
从off开始读取数据到b中
*/
func readFile() {
	userFile := "zou.txt"
	fl, err := os.Open(userFile)
	if err != nil {
		log.Println(userFile, err)
		return
	}
	defer fl.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := fl.Read(buf)
		if 0 == n {
			break
		}
		os.Stdout.Write(buf[:n])
	}
}

func main() {
	dir()
	writeFile()
	readFile()
}
