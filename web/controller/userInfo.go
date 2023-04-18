
package controller

import (
	"education/service"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"encoding/json"
)

type Application struct {
	Setup *service.ServiceSetup
}
var PersonalSpaceMap map[string]PersonalSpace //string:LoginName
type User struct {
	LoginName	string `db:"LoginName"`
	Password	string `db:"Password"` 
	Identity 	string `db:"Identity"`
	IdentificationCode	string `db:"IdentificationCode"`
}
type PersonalSpace struct{
	CetPtrList []*CetWaitingToApproveStruct
	EduPtrList []*EduWaitingToApproveStruct
}

type ProposerStruct struct {
	LoginName string
	ProTime string
	ProNo string 
	StatusCode int 
	ConfirmLoginName string
}

type CetWaitingToApproveStruct struct{
	Proposer ProposerStruct
	CetItem service.CertificateObj
}
type EduWaitingToApproveStruct struct{
	Proposer ProposerStruct
	EduItem service.Education
}
var user User
var (
	CetWaitingToApproveList []CetWaitingToApproveStruct
	EduWaitingToApproveList []EduWaitingToApproveStruct
)
var (
	dbConn *sql.DB
	err error
)
const(
	Admin = "Admin"
	Member = "Member"
	Enterprises = "Enterprises"
	Individual = "Individual"
)
func init() {
	PersonalSpaceMapInit()
	//UsersInit()
	MySqlInit()
}
/*
func UsersInit(){
	admin := User{LoginName:"admin", Password:"123456", Identity:Admin}
	school := User{LoginName:"school", Password:"123456",Identity:Member}
	bob := User{LoginName:"bob", Password:"123456", Identity:Individual, IdentificationCode:"32243120011221001X"}
	huawei := User{LoginName:"huawei", Password:"123456", Identity:Enterprises}

	users = append(users, admin)
	users = append(users, school)
	users = append(users, bob)
	users = append(users, huawei)
}
*/
func PersonalSpaceMapInit(){
	PersonalSpaceMap = make(map[string]PersonalSpace)
	for i := 0; i < len(CetWaitingToApproveList); i++{
		ptr := &CetWaitingToApproveList[i]
		Ps := PersonalSpaceMap[ptr.Proposer.LoginName]
		Ps.CetPtrList = append(Ps.CetPtrList,ptr)
		PersonalSpaceMap[ptr.Proposer.LoginName] = Ps
	}
	for i := 0; i < len(EduWaitingToApproveList); i++{
		ptr := &EduWaitingToApproveList[i]
		Ps := PersonalSpaceMap[ptr.Proposer.LoginName]
		Ps.EduPtrList = append(Ps.EduPtrList,ptr)
		PersonalSpaceMap[ptr.Proposer.LoginName] = Ps
	}
}
func AddCetProposal(cet CetWaitingToApproveStruct){
	MySqlInsertCetProposal(&cet)
	CetWaitingToApproveList = append(CetWaitingToApproveList,cet)
	Ps := PersonalSpaceMap[cet.Proposer.LoginName]
	Ps.CetPtrList = append(Ps.CetPtrList,&CetWaitingToApproveList[len(CetWaitingToApproveList)-1])
	PersonalSpaceMap[cet.Proposer.LoginName] = Ps
}
func AddEduProposal(edu EduWaitingToApproveStruct){
	MySqlInsertEduProposal(&edu)
	EduWaitingToApproveList = append(EduWaitingToApproveList,edu)
	Ps := PersonalSpaceMap[edu.Proposer.LoginName]
	Ps.EduPtrList = append(Ps.EduPtrList,&EduWaitingToApproveList[len(EduWaitingToApproveList)-1])
	PersonalSpaceMap[edu.Proposer.LoginName] = Ps
}
func (c *CetWaitingToApproveStruct) UpdateStatusCode(statusCode int,Cname string) (string, bool){
	if statusCode != 0 && statusCode != -1 && statusCode != 1{
		return "传入参数有误,函数执行失败",false
	}
	if c.Proposer.StatusCode == statusCode{
		return "状态一致无需改变",true
	}
	MySqlUpdateCetProposal(statusCode,Cname,c.Proposer.ProNo)
	c.Proposer.StatusCode = statusCode
	return "更新成功",true
}
func (e *EduWaitingToApproveStruct) UpdateStatusCode(statusCode int,Cname string) (string, bool){
	if statusCode != 0 && statusCode != -1 && statusCode != 1{
		return "传入参数有误,函数执行失败",false
	}
	if e.Proposer.StatusCode == statusCode{
		return "状态一致无需改变",true
	}
	MySqlUpdateEduProposal(statusCode,Cname,e.Proposer.ProNo)
	e.Proposer.StatusCode = statusCode
	return "更新成功",true
}

/* 数据库操作 */
func MySqlInit(){
	dbConn, err = sql.Open("mysql", "root:123qwe321ewq@tcp(localhost:3306)/education?allowNativePasswords=true")
	dbConn.Ping()
	if err != nil{
		fmt.Println("数据库连接失败")
		return
	}
}
func MySqlLoginCheck(loginName string,password string) bool{
	query := "select * from Users where LoginName="+"\""+loginName+"\""
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
	for rows.Next(){
		rows.Scan(&user.LoginName,&user.Password,&user.Identity,&user.IdentificationCode)
	}
	if user.Password == password{
		return true
	}else{
		return false
	}
}

func MySqlInsertEduProposal(edu *EduWaitingToApproveStruct){
	stm, err := dbConn.Prepare("insert into EduProposals values(?,?,?,?,?,default)")
	defer func(){
		if stm != nil{
			stm.Close()
		}
	}()
	if err != nil{
		fmt.Println("预处理失败")
		return
	}
	b, err := json.Marshal(edu.EduItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := stm.Exec(edu.Proposer.ProNo,edu.Proposer.LoginName,
						edu.Proposer.ProTime,edu.Proposer.StatusCode,b)
	if err != nil {
		fmt.Println(err)
		fmt.Println("sql执行失败")
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		fmt.Println("结果获取失败")
		return
	}				
	if count > 0{
		fmt.Println("新增成功")
	}else{
		fmt.Println("新增失败")
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
		fmt.Println("预处理失败")
		return
	}
	b, err := json.Marshal(cet.CetItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := stm.Exec(cet.Proposer.ProNo,cet.Proposer.LoginName,
						cet.Proposer.ProTime,cet.Proposer.StatusCode,b)
	if err != nil {
		fmt.Println(err)
		fmt.Println("sql执行失败")
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		fmt.Println("结果获取失败")
		return
	}				
	if count > 0{
		fmt.Println("新增成功")
	}else{
		fmt.Println("新增失败")
	}	
}
func MySqlUpdateEduProposal(stc int,Cname string,Pno string){
	var updateSql string
	if Cname == ""{
		updateSql = "update EduProposals set Stcode = ? where Pno = ?"
	}else{
		updateSql = "update EduProposals set Stcode = ?, Cname = ? where Pno = ?"
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
	if Cname == ""{
		stm.Exec(stc,Pno)
	}else{
		stm.Exec(stc,Cname,Pno)
	}
}
func MySqlUpdateCetProposal(stc int,Cname string,Pno string){
	var updateSql string
	if Cname == ""{
		updateSql = "update CetProposals set Stcode = ? where Pno = ?"
	}else{
		updateSql = "update CetProposals set Stcode = ?, Cname = ? where Pno = ?"
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
	if Cname == ""{
		res, err := stm.Exec(stc,Pno)
		if err != nil {
			fmt.Println(err)
			fmt.Println("sql执行失败")
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			fmt.Println("结果获取失败")
			return
		}				
		if count > 0{
			fmt.Println("新增成功")
		}else{
			fmt.Println("新增失败")
		}
			
	}else{
		res, err := stm.Exec(stc,Cname,Pno)
		if err != nil {
			fmt.Println(err)
			fmt.Println("sql执行失败")
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			fmt.Println("结果获取失败")
			return
		}				
		if count > 0{
			fmt.Println("新增成功")
		}else{
			fmt.Println("新增失败")
		}
	}	
}