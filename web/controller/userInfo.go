package controller

import (
	"education/service"
	"time"
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
	StatusCode int `db:"StatusCode"`
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
const (
	Individual = "Individual"
	Enterprises = "Enterprises"
	Admin = "Admin"
	Member = "Member"
)
var (
	CetWaitingToApproveList []CetWaitingToApproveStruct
	EduWaitingToApproveList []EduWaitingToApproveStruct
	UserWaitingToApproveList []User
)
//init 初始化
func init() {
	PersonalSpaceMapInit()
	MySqlInit()
}
func PersonalSpaceMapInit(){
	PersonalSpaceMap = make(map[string]PersonalSpace)
}
//添加证书申请
func AddCetProposal(cet *service.CertificateObj,ProLoginName string){
	tStr := time.Now().Format("2006-01-02 15:04:05")
	proNo := service.Sha256(ProLoginName+tStr)
	CetWaitingToApproveObj :=CetWaitingToApproveStruct{
		Proposer:ProposerStruct{
			LoginName:ProLoginName,
			ProTime:tStr,
			ProNo:proNo[0:10],
			StatusCode:0,
		},
		CetItem:*cet,
	}
	MySqlInsertCetProposal(&CetWaitingToApproveObj)
	CetWaitingToApproveList = append(CetWaitingToApproveList,CetWaitingToApproveObj)
	Ps := PersonalSpaceMap[ProLoginName]
	Ps.CetPtrList = append(Ps.CetPtrList,&CetWaitingToApproveList[len(CetWaitingToApproveList)-1])
	PersonalSpaceMap[ProLoginName] = Ps
}
func AddEduProposal(edu *service.Education,ProLoginName string){
	tStr := time.Now().Format("2006-01-02 15:04:05")
	proNo := service.Sha256(ProLoginName+tStr)
	EduWaitingToApproveObj :=EduWaitingToApproveStruct{
		Proposer:ProposerStruct {
			LoginName:ProLoginName,
			ProTime:tStr,
			ProNo:proNo[0:10],
			StatusCode:0,
		},
		EduItem:*edu,
	}
	MySqlInsertEduProposal(&EduWaitingToApproveObj)
	EduWaitingToApproveList = append(EduWaitingToApproveList,EduWaitingToApproveObj)
	Ps := PersonalSpaceMap[ProLoginName]
	Ps.EduPtrList = append(Ps.EduPtrList,&EduWaitingToApproveList[len(EduWaitingToApproveList)-1])
	PersonalSpaceMap[ProLoginName] = Ps
}
func (c *CetWaitingToApproveStruct) UpdateStatusCode(statusCode int,Cname string) (string, bool){
	if statusCode != 0 && statusCode != -1 && statusCode != 1{
		return "传入参数有误,函数执行失败",false
	}
	if c.Proposer.StatusCode == statusCode{
		return "状态一致无需改变",true
	}
	MySqlUpdateProposal(statusCode,Cname,c.Proposer.ProNo,"Cet")
	c.Proposer.StatusCode = statusCode
	c.Proposer.ConfirmLoginName = Cname
	return "更新成功",true
}
func (e *EduWaitingToApproveStruct) UpdateStatusCode(statusCode int,Cname string) (string, bool){
	if statusCode != 0 && statusCode != -1 && statusCode != 1{
		return "传入参数有误,函数执行失败",false
	}
	if e.Proposer.StatusCode == statusCode{
		return "状态一致无需改变",true
	}
	MySqlUpdateProposal(statusCode,Cname,e.Proposer.ProNo,"Edu")
	e.Proposer.StatusCode = statusCode
	e.Proposer.ConfirmLoginName = Cname
	return "更新成功",true
}
//添加用户，注册
func AddUserProposal(user *User){
	MySqlInsertUsers(user)
	UserWaitingToApproveList = append(UserWaitingToApproveList,*user)
}
func (u *User) UpdateStatusCode(statusCode int) (string, bool){
	if statusCode != 0 && statusCode != -1 && statusCode != 1{
		return "传入参数有误,函数执行失败",false
	}
	if u.StatusCode == statusCode{
		return "状态一致无需改变",true
	}
	MySqlUpdateUsers(u.LoginName,statusCode)
	u.StatusCode = statusCode
	return "更新成功",true
}
