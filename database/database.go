package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8")
	checkErr(err)
	defer db.Close()

	//insert
	stmt, err := db.Prepare("INSERT `userinfo` SET username=?,department=?,created=?")
	checkErr(err)

	res, err := stmt.Exec("zoujiejun", "Research and development department", "2018-10-23")
	checkErr(err)
	res, err = stmt.Exec("zoujiejun", "Research and development department", "2018-10-23")
	checkErr(err)
	res, err = stmt.Exec("zoujiejun", "Research and development department", "2018-10-23")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)
	fmt.Println("-----------------------")

	//update
	stmt, err = db.Prepare("UPDATE userinfo set username=? where uid=?")
	checkErr(err)

	res, err = stmt.Exec("luxingyu", id)
	checkErr(err)

	//rows change
	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)
	fmt.Println("-----------------------")

	//select
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	fmt.Println("-----------------------")
	//delete
	stmt, err = db.Prepare("DELETE FROM userinfo WHERE uid=?")
	checkErr(err)

	res, err = stmt.Exec(id)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

}
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
