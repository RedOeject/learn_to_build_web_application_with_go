//https://github.com/RedOeject/build-web-application-with-golang/blob/master/zh/07.4.md
package main

import (
	"html/template"
	"os"
)

type Person struct {
	UserName string
	Hobby    []string
	Friends  []*Friend
}

type Friend struct {
	Fname string
}

func main() {
	ifElse()
}

//只能调用大写开头的导出字段，未导出字段会报错，但是不存在的字段不会报错，而是输出为空
//如果模板中输出{{.}}，这个一般应用于字符串对象，默认会调用fmt包输出字符串的内容。
func output() {
	t := template.New("fieldname example")
	t, _ = t.Parse("hello {{.UserName}}!")
	p := Person{UserName: "zou"}
	t.Execute(os.Stdout, p)
}

//使用{{with …}}…{{end}}和{{range …}}{{end}}来进行数据的输出。
//{{range}} 这个和Go语法里面的range类似，循环操作数据
//{{with}}操作是指当前对象的值，类似上下文的概念
func rangeAndWith() {
	f1 := Friend{"minux.ma"}
	f2 := Friend{"donghaonian"}
	t := template.New("fieldname example")
	t, _ = t.Parse(`hello {{.UserName}}!
	{{range.Hobby}}
		a hobby: {{.}}
	{{end}}
	{{with.Friends}}
	{{range .}}
		my friend name is {{.Fname}}
	{{end}}
	{{end}}
	`)
	p := Person{UserName: "Zou",
		Hobby:   []string{"football", "game"},
		Friends: []*Friend{&f1, &f2}}
	t.Execute(os.Stdout, p)
}

//在Go模板里面如果需要进行条件判断，那么我们可以使用和Go语言的if-else语法类似的方式来处理，如果pipeline为空，那么if就认为是false
//注意：if里面无法使用条件判断，例如.Mail=="astaxie@gmail.com"，这样的判断是不正确的，if里面只能是bool值
func ifElse() {
	tEmpty := template.New("template test")
	//检测模板是否正确，例如大括号是否匹配，注释是否正确的关闭，变量是否正确的书写。
	tEmpty = template.Must(tEmpty.Parse("空 pipeline if demo: {{if ``}} 不会输出. {{end}}\n"))
	tEmpty.Execute(os.Stdout, nil)

	tWithValue := template.New("template test")
	tWithValue = template.Must(tWithValue.Parse("不为空的 pipeline if demo: {{if `anything`}} 我有内容，我会输出. {{end}}\n"))
	tWithValue.Execute(os.Stdout, nil)

	tIfElse := template.New("template test")
	tIfElse = template.Must(tIfElse.Parse("if-else demo: {{if `anything`}} if部分 {{else}} else部分.{{end}}\n"))
	tIfElse.Execute(os.Stdout, nil)
}
