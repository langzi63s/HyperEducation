package controller

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"encoding/json"
)

var (
	dbConn *sql.DB
	err error
)
/* 数据库操作 */
func MySqlInit(){
	dbConn, err = sql.Open("mysql", "root:123qwe321ewq@tcp(localhost:3306)/education?allowNativePasswords=true")
	dbConn.Ping()
	if err != nil{
		fmt.Println("数据库连接失败")
		return
	}
	dbConn.Exec("truncate table EduProposals")
	dbConn.Exec("truncate table CetProposals")
}
func MySqlLoginCheck(loginName string,password string) (bool,int){
	query := "select * from Users where LoginName="+"\""+loginName+"\""
	rows, err := dbConn.Query(query)
	if err != nil {
		fmt.Println(err)
		return false,-1
	}
	defer func(){
		if rows != nil{
			rows.Close()
		}
	}()
	for rows.Next(){
		rows.Scan(&user.LoginName,&user.Password,&user.Identity,&user.IdentificationCode,&user.StatusCode)
	}
	if user.LoginName != loginName{
		return false,-1
	}
	if user.Password == password{
		return true,user.StatusCode
	}
	return false,user.StatusCode
}
func MySqlUnconfirmedUserInit(){
	query := "select * from Users where StatusCode = 0"
	rows, err := dbConn.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(){
		if rows != nil{
			rows.Close()
		}
	}()
	for rows.Next(){
		rows.Scan(&user.LoginName,&user.Password,&user.Identity,&user.IdentificationCode,&user.StatusCode)
		UserWaitingToApproveList = append(UserWaitingToApproveList,user)
	}
}
func MySqlIdentificationCodeExist(IdentificationCode string) bool{
	query := "select IdentificationCode from Users where IdentificationCode="+"\""+IdentificationCode+"\""
	rows, err := dbConn.Query(query)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func(){
		if rows != nil{
			rows.Close()
		}
	}()
	var ic string
	for rows.Next(){
		rows.Scan(&ic)
	}
	if ic == IdentificationCode{
		return true
	}
	return false
}
func MySqlInsertUsers(user *User){
	var insertSql string
	insertSql = "insert into Users values(?,?,?,?,?)"
	stm, err := dbConn.Prepare(insertSql)
	defer func(){
		if stm != nil{
			stm.Close()
		}
	}()
	if err != nil{
		fmt.Println("预处理失败")
		return
	}
	stm.Exec(user.LoginName,user.Password,user.Identity,user.IdentificationCode,user.StatusCode)
}
func MySqlUpdateUsers(LoginName string,No int){
	var updateSql string
	updateSql = "update Users set StatusCode = ? where LoginName = ?"
	stm, err := dbConn.Prepare(updateSql)
	defer func(){
		if stm != nil{
			stm.Close()
		}
	}()
	if err != nil{
		fmt.Println("预处理失败")
		return
	}
	_, err = stm.Exec(No, LoginName)
}
func MySqlInsertEduProposal(edu *EduWaitingToApproveStruct){
	stm, err := dbConn.Prepare("insert into EduProposals values(?,?,?,?,?,default)")
	defer func(){
		if stm != nil{
			stm.Close()
		}
	}()
	if err != nil{
		fmt.Println("sql预处理失败:InsertEduProposal")
		return
	}
	b, err := json.Marshal(edu.EduItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stm.Exec(edu.Proposer.ProNo,edu.Proposer.LoginName,
						edu.Proposer.ProTime,edu.Proposer.StatusCode,b)
	if err != nil {
		fmt.Printf("sql执行失败:%s",err)
		return
	}
}
func MySqlInsertCetProposal(cet *CetWaitingToApproveStruct){
	stm, err := dbConn.Prepare("insert into CetProposals values(?,?,?,?,?,default)")
	defer func(){
		if stm != nil{
			stm.Close()
		}
	}()
	if err != nil{
		fmt.Println("sql预处理失败:InsertCetProposal")
		return
	}
	b, err := json.Marshal(cet.CetItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stm.Exec(cet.Proposer.ProNo,cet.Proposer.LoginName,
						cet.Proposer.ProTime,cet.Proposer.StatusCode,b)
	if err != nil {
		fmt.Printf("sql执行失败:%s",err)
		return
	}
}
func MySqlUpdateProposal(stc int,Cname string,Pno string,typ string){
	sheetName := typ + "Proposals"
	bytes := typ + "Bytes"
	var updateSql string
	if stc != 0{
		updateSql = "update "+ sheetName+ " set Stcode = ?, Cname = ?,"+bytes+" = NULL where Pno = ?"
	}
	stm, err := dbConn.Prepare(updateSql)
	defer func(){
		if stm != nil{
			stm.Close()
		}
	}()
	if err != nil{
		fmt.Println("预处理失败")
		return
	}
	_, err = stm.Exec(stc,Cname,Pno)
}
