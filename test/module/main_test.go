package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type ModuleListResponse struct {
	ModuleList []Module `json:"Data"`
}

type Module struct {
	ModuleID             int         `json:"ModuleID"`
	IsMainModule         bool        `json:"IsMainModule"`
	MainModuleType       string      `json:"MainModuleType"`
	ModuleType           string      `json:"ModuleType"`
	NeedCompile          bool        `json:"NeedCompile"`
	IsClosed             bool        `json:"IsClosed"`
	ModuleCode           string      `json:"ModuleCode"`
	Platform             string      `json:"Platform"`
	IsAny                bool        `json:"IsAny"`
	IsX86                bool        `json:"IsX86"`
	IsX64                bool        `json:"IsX64"`
	IsARM32              bool        `json:"IsARM32"`
	IsMIPS32             bool        `json:"IsMIPS32"`
	IsARM64              bool        `json:"IsARM64"`
	IsMIPS64             bool        `json:"IsMIPS64"`
	IsLoongArch32        bool        `json:"IsLoongArch32"`
	IsLoongArch64        bool        `json:"IsLoongArch64"`
	IsOtherArch32        bool        `json:"IsOtherArch32"`
	IsOtherArch64        bool        `json:"IsOtherArch64"`
	WithDebug            bool        `json:"WithDebug"`
	IsSystemModule       bool        `json:"IsSystemModule"`
	IsCommonModule       bool        `json:"IsCommonModule"`
	IsSourceCodeRelease  bool        `json:"IsSourceCodeRelease"`
	CommandName          string      `json:"CommandName"`
	CommandID            int         `json:"CommandID"`
	ModuleName           string      `json:"ModuleName"`
	AppName              string      `json:"AppName"`
	OutDirs              string      `json:"OutDirs"`
	Principal            string      `json:"Principal"`
	ModuleStatus         string      `json:"ModuleStatus"`
	LanType              string      `json:"LanType"`
	LanTypeName          string      `json:"LanTypeName"`
	PrincipalDepartment  string      `json:"PrincipalDepartment"`
	SVNModuleName        string      `json:"SVNModuleName"`
	TypeID               int         `json:"TypeID"`
	NeedDigitalSignature bool        `json:"NeedDigitalSignature"`
	HasSVN               bool        `json:"HasSVN"`
	ModuleVersion        interface{} `json:"ModuleVersion"` // Change this to the actual type if known
	Synced               bool        `json:"Synced"`
	Remark               string      `json:"Remark"`
	ThirdVersion         interface{} `json:"ThirdVersion"` // Change this to the actual type if known
	TestByTester         bool        `json:"TestByTester"`
}

func TestModuleList(t *testing.T) {

	// 准备要发送的数据，这里使用字节数组作为示例
	requestData := []byte(`{"pageSize":20,"pageIndex":1,"IsAll":false,"ModuleName":"","ModuleType":"","OsType":"","DevLang":"","CompileCommand":""}`)

	// 构建请求对象
	req, err := http.NewRequest("POST", "http://10.100.3.138/newModule/Module/GetModuleListWithCondition", bytes.NewBuffer(requestData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头，如果需要的话
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "beegosessionID=33d5b38eda1d11bac5d4262c7dc74c79; auth=true; token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ6amhvbmciLCJleHAiOjE3MDc0NDQ0NjcsImlhdCI6MTcwNzM1ODA2NywiaXNzIjoiZGVwb3QiLCJsYXN0SVAiOiIxMC4yMjAuMjMuMjcifQ.LHA0p_MyNTexyuSv-5DWedEF2UlVEQ1nR1ndWtS_jpBCWSL0XNOW_4Q2Q-UykTIRM4zeqHJ41MyStY9U1vcPpg; .AspNetCore.Cookies=CfDJ8ACAe0sV9F1BgtD8EC0c6eW-q5AJtybqbr3-cKGPOo1IBWO9W6fXYA9-wBtZ8kxnzv4zSnmFzaowirxWXCkSEkTHMgxO0RUbdklfJVsgBOhg5pGIAVqY-w7ey7MjUWFhrNeVv8iYfl7qxDX9_x6YmNvkyzl5Q7nIr21CPeP9-yGOElpzKyUU9L5Ydze2a9y3d3-3-RbbcVD0ljw0LYNEr2ubzuu_4EVvrEyMnBbEzodaJ9M9e6bd8fMMvVYvAV_UY0SIgXwZISROdTonC_8wLCIYJL-wOUEIAZV0uu5Ckis7")

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	var responseBytes bytes.Buffer
	_, err = responseBytes.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var response ModuleListResponse
	err = json.Unmarshal(responseBytes.Bytes(), &response)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}

type Data struct {
	ModuleVer           string        `json:"ModuleVer"`
	IsStable            bool          `json:"IsStable"`
	ModuleVersionStatus string        `json:"ModuleVersionStatus"`
	DocUrl              string        `json:"DocUrl"`
	VerRemark           string        `json:"VerRemark"`
	VerNote             interface{}   `json:"VerNote"`
	BuildDate           string        `json:"BuildDate"`
	BuildTime           interface{}   `json:"BuildTime"`
	SummaryContent      interface{}   `json:"SummaryContent"`
	VerDetail           string        `json:"VerDetail"`
	Version             string        `json:"Version"`
	VersionTitle        string        `json:"VersionTitle"`
	IsPublish           bool          `json:"IsPublish"`
	Children            []interface{} `json:"Children"`
	Remark              interface{}   `json:"Remark"`
	ModuleRunTime       interface{}   `json:"ModuleRunTime"`
	PublishDate         interface{}   `json:"PublishDate"`
	IsDiscard           bool          `json:"IsDiscard"`
}

type Response struct {
	Data         []Data      `json:"Data"`
	ErrorCode    interface{} `json:"ErrorCode"`
	ErrorMessage interface{} `json:"ErrorMessage"`
	Page         interface{} `json:"Page"`
	State        int         `json:"State"`
}

func TestGetModuleVerInfoList(t *testing.T) {

	// 准备要发送的数据，这里使用字节数组作为示例
	requestData := []byte("10964")

	// 构建请求对象
	req, err := http.NewRequest("POST", "http://10.100.3.138/newModule/Module/GetModuleVerInfoList", bytes.NewBuffer(requestData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头，如果需要的话
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "beegosessionID=33d5b38eda1d11bac5d4262c7dc74c79; auth=true; token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ6amhvbmciLCJleHAiOjE3MDc0NDQ0NjcsImlhdCI6MTcwNzM1ODA2NywiaXNzIjoiZGVwb3QiLCJsYXN0SVAiOiIxMC4yMjAuMjMuMjcifQ.LHA0p_MyNTexyuSv-5DWedEF2UlVEQ1nR1ndWtS_jpBCWSL0XNOW_4Q2Q-UykTIRM4zeqHJ41MyStY9U1vcPpg; .AspNetCore.Cookies=CfDJ8ACAe0sV9F1BgtD8EC0c6eW-q5AJtybqbr3-cKGPOo1IBWO9W6fXYA9-wBtZ8kxnzv4zSnmFzaowirxWXCkSEkTHMgxO0RUbdklfJVsgBOhg5pGIAVqY-w7ey7MjUWFhrNeVv8iYfl7qxDX9_x6YmNvkyzl5Q7nIr21CPeP9-yGOElpzKyUU9L5Ydze2a9y3d3-3-RbbcVD0ljw0LYNEr2ubzuu_4EVvrEyMnBbEzodaJ9M9e6bd8fMMvVYvAV_UY0SIgXwZISROdTonC_8wLCIYJL-wOUEIAZV0uu5Ckis7")

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	var responseBytes bytes.Buffer
	_, err = responseBytes.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var tt Response
	err = json.Unmarshal(responseBytes.Bytes(), &tt)
	if err != nil {
		panic(err)
	}

	fmt.Println(tt)
}
