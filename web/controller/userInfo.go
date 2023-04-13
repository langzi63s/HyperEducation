
package controller

import "education/service"

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName	string
	Password	string
	Identity 	string
	IdentificationCode	string
}

type Proposer struct {
	Userinfo User
	ProTime string
}

type CetWaitingToApproveStruct struct{
	proposer Proposer
	cetItem service.CertificateObj
}
type EduWaitingToApproveStruct struct{
	proposer Proposer
	eduItem service.Education
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

	admin := User{LoginName:"admin", Password:"123456", Identity:Admin}
	school := User{LoginName:"school", Password:"123456",Identity:Member}
	bob := User{LoginName:"bob", Password:"123456", Identity:Individual}
	huawei := User{LoginName:"huawei", Password:"123456", Identity:Enterprises}

	users = append(users, admin)
	users = append(users, school)
	users = append(users, bob)
	users = append(users, huawei)

}