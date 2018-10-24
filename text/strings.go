/*
字符串在我们平常的Web开发中经常用到，包括用户的输入，数据库读取的数据等，我们经常需要对字符串进行分割、连接、转换等操作，本小节将通过Go标准库中的strings和strconv两个包中的函数来讲解如何进行有效快速的操作。
*/
package main

//字符串操作
//func Contains(s, substr string) bool  字符串s中是否包含substr，返回bool值
//func Join(a []string, sep string) string 字符串链接，把slice a通过sep链接起来
//func Index(s, sep string) int 在字符串s中查找sep所在的位置，返回位置值，找不到返回-1
//func Repeat(s string, count int) string  重复s字符串count次，最后返回重复的字符串
//func Replace(s, old, new string, n int) string  在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
//func Split(s, sep string) []string 把s字符串按照sep分割，返回slice
//func Trim(s string, cutset string) string  在s字符串的头部和尾部去除cutset指定的字符串
//func Fields(s string) []string 去除s字符串的空格符，并且按照空格分割返回slice

//字符串转换
//Append 系列函数将整数等转换为字符串后，添加到现有的字节数组中。
//Format 系列函数把其他类型的转换为字符串
//Parse 系列函数把字符串转换为其他类型
