package main

import (
	//"encoding/json"
	"fmt"
	"education/sdkInit"
	"education/service"
	"education/web"
	"education/web/controller"
	"os"
)

const (
	cc_name = "simplecc"
	cc_version = "1.0.0"
)
func main() {
	// init orgs information
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/education/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},

	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/education/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    os.Getenv("GOPATH")+"/src/education/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")
	//提交申请
	
	edu := service.Education{
		Name: "刘嘉楷",
		Gender: "男",
		Nation: "汉",
		EntityID: "32243120011221001X",
		Place: "江西吉安",
		BirthDay: "2001年12月21日",
		EnrollDate: "2019年9月",
		GraduationDate: "2023年7月",
		SchoolName: "华东理工大学",
		Major: "计算机科学与技术",
		QuaType: "普通",
		Length: "四年",
		Mode: "普通全日制",
		Level: "本科",
		Graduation: "毕业",
		CertNo: "108294189014681257",
		Photo: "/static/photo/11.png",
	}
	cet := service.CertificateObj{
		Name:"刘嘉楷",
		Gender:"男",
		EntityID:"32243120011221001X",
		Level:"四",
		CertNo:"211231005003443",
		TestTime:"2021年09月",
		TestNo:"310052221201818",
		Score:"520",
	}
	cet2 := service.CertificateObj{
		Name:"刘嘉楷",
		Gender:"男",
		EntityID:"32243120011221001X",
		Level:"六",
		CertNo:"221231005003442",
		TestTime:"2022年12月",
		TestNo:"310063222201815",
		Score:"510",
	}
	controller.AddEduProposal(&edu,"bob")
	controller.AddCetProposal(&cet,"bob")
	
	
	serviceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err!=nil{
		fmt.Println()
		os.Exit(-1)
	}
	serviceSetup.SaveEdu(edu)
	controller.EduWaitingToApproveList[0].UpdateStatusCode(1,"school")
	serviceSetup.SaveCet(cet)
	controller.CetWaitingToApproveList[0].UpdateStatusCode(1,"school")
	controller.AddCetProposal(&cet2,"bob")
	app := controller.Application{
		Setup: serviceSetup,
	}
	web.WebStart(app)
}