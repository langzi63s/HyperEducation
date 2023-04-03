
package controller

import "education/service"

type Application struct {
	Setup *service.ServiceSetup
}
type Userinfo struct {
	EntityID	string
	
}
type User struct {
	LoginName	string
	Password	string
	Identity 	string
	IsAdmin		string  //具有学历修改权限
	Info		*Userinfo
}


var users []User
const(
	Admin = "Admin"
	Member = "Member"
	Enterprises = "Enterprises"
	Individual = "Individual"
)
func init() {

	admin := User{LoginName:"admin", Password:"123456", Identity:Admin,IsAdmin:"T"}
	school := User{LoginName:"school", Password:"123456",Identity:Member, IsAdmin:"F"}
	bob := User{LoginName:"bob", Password:"123456", Identity:Individual,IsAdmin:"F"}
	huawei := User{LoginName:"huawei", Password:"123456", Identity:Enterprises,IsAdmin:"F"}

	users = append(users, admin)
	users = append(users, school)
	users = append(users, bob)
	users = append(users, huawei)

}