
package controller

import (
	"education/service"
)

type Application struct {
	Setup *service.ServiceSetup
}
var PersonalSpaceMap map[string]PersonalSpace //string:LoginName
type User struct {
	LoginName	string
	Password	string
	Identity 	string
	IdentificationCode	string
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
var users []User
var CetWaitingToApproveList []CetWaitingToApproveStruct
var EduWaitingToApproveList []EduWaitingToApproveStruct
const(
	Admin = "Admin"
	Member = "Member"
	Enterprises = "Enterprises"
	Individual = "Individual"
)
func init() {
	PersonalSpaceMapInit()
	UsersInit()
}
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
	CetWaitingToApproveList = append(CetWaitingToApproveList,cet)
	Ps := PersonalSpaceMap[cet.Proposer.LoginName]
	Ps.CetPtrList = append(Ps.CetPtrList,&CetWaitingToApproveList[len(CetWaitingToApproveList)-1])
	PersonalSpaceMap[cet.Proposer.LoginName] = Ps
}
func AddEduProposal(edu EduWaitingToApproveStruct){
	EduWaitingToApproveList = append(EduWaitingToApproveList,edu)
	Ps := PersonalSpaceMap[edu.Proposer.LoginName]
	Ps.EduPtrList = append(Ps.EduPtrList,&EduWaitingToApproveList[len(EduWaitingToApproveList)-1])
	PersonalSpaceMap[edu.Proposer.LoginName] = Ps
}
func (c *CetWaitingToApproveStruct) UpdateStatusCode(statusCode int) (string, bool){
	if statusCode != 0 && statusCode != -1 && statusCode != 1{
		return "传入参数有误,函数执行失败",false
	}
	if c.Proposer.StatusCode == statusCode{
		return "状态一致无需改变",true
	}
	c.Proposer.StatusCode = statusCode
	return "更新成功",true
}
func (e *EduWaitingToApproveStruct) UpdateStatusCode(statusCode int) (string, bool){
	if statusCode != 0 && statusCode != -1 && statusCode != 1{
		return "传入参数有误,函数执行失败",false
	}
	if e.Proposer.StatusCode == statusCode{
		return "状态一致无需改变",true
	}
	e.Proposer.StatusCode = statusCode
	return "更新成功",true
}