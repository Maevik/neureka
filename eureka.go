package neureka

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// AppAllResponse Eureka响应所有应用信息
type AppAllResponse struct {
	Applications Applications `json:"applications"`
}

type Applications struct {
	Versions__delta string          `json:"versions__delta"`
	Apps__hashcode  string          `json:"apps__hashcode"`
	Application     []EurekaAppInfo `json:"application"`
}

// AppResponse Eureka响应应用信息
type AppResponse struct {
	Application EurekaAppInfo `json:"application"`
}

// EurekaAppInfo 应用信息
type EurekaAppInfo struct {
	Name     string              `json:"name"`
	Instance []EurekaAppInstance `json:"instance"`
}

// EurekaAppInstance Eureka应用实例信息
type EurekaAppInstance struct {
	InstanceId                    string                 `json:"instanceId"`
	App                           string                 `json:"app"`
	AppGroupName                  string                 `json:"appGroupName"`
	IpAddr                        string                 `json:"ipAddr"`
	Sid                           string                 `json:"sid"`
	Port                          InstancePort           `json:"port"`
	SecurePort                    InstanceSecurePort     `json:"securePort"`
	HealthCheckUrl                string                 `json:"healthCheckUrl"`
	StatusPageUrl                 string                 `json:"statusPageUrl"`
	HomePageUrl                   string                 `json:"homePageUrl"`
	VipAddress                    string                 `json:"vipAddress"`
	SecureVipAddress              string                 `json:"secureVipAddress"`
	CountryId                     int                    `json:"countryId"`
	DataCenterInfo                DataCenterInfo         `json:"dataCenterInfo"`
	HostName                      string                 `json:"hostName"`
	Status                        string                 `json:"status"`
	Overriddenstatus              string                 `json:"overriddenstatus"`
	LeaseInfo                     LeaseInfo              `json:"leaseInfo"`
	IsCoordinatingDiscoveryServer string                 `json:"isCoordinatingDiscoveryServer"`
	Metadata                      map[string]interface{} `json:"metadata"`
	LastUpdatedTimestamp          string                 `json:"lastUpdatedTimestamp"`
	LastDirtyTimestamp            string                 `json:"lastDirtyTimestamp"`
	ActionType                    string                 `json:"actionType"`
}

// InstancePort 实例端口
type InstancePort struct {
	Enable string `json:"@enabled"`
	Value  int    `json:"$"`
}

// InstanceSecurePort 实例安全端口
type InstanceSecurePort struct {
	Enable string `json:"@enabled"`
	Value  int    `json:"$"`
}

// DataCenterInfo
type DataCenterInfo struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

// LeaseInfo
type LeaseInfo struct {
	RenewalIntervalInSecs int64 `json:"renewalIntervalInSecs"`
	DurationInSecs        int64 `json:"durationInSecs"`
	RegistrationTimestamp int64 `json:"registrationTimestamp"`
	LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp"`
	RenewalTimestamp      int64 `json:"renewalTimestamp"`
	EvictionTimestamp     int64 `json:"evictionTimestamp"`
	ServiceUpTimestamp    int64 `json:"serviceUpTimestamp"`
}

// NewEurekaAppInstance 实例化一个Eureka实例
func NewEurekaAppInstance(cnf EurekaClientConfig) EurekaAppInstance {
	// 毫秒时间戳
	ms := GetMs()
	s_ms := strconv.FormatInt(ms, 10)

	porten := "true"
	sporten := "false"
	if cnf.InstancePort == 443 {
		porten = "false"
		sporten = "true"
	}

	// 实例化eureka信息
	euk := EurekaAppInstance{
		InstanceId:   cnf.Id(),
		App:          cnf.AppName,
		AppGroupName: "",
		IpAddr:       cnf.InstanceIp,
		Sid:          "na",
		Port: InstancePort{
			Enable: porten,
			Value:  cnf.InstancePort,
		},
		SecurePort: InstanceSecurePort{
			Enable: sporten,
			Value:  443,
		},
		HealthCheckUrl:   cnf.InstanceHealthCheckUrl,
		StatusPageUrl:    cnf.InstanceStatusUrl,
		HomePageUrl:      cnf.InstanceHomePageUrl,
		VipAddress:       cnf.AppName,
		SecureVipAddress: cnf.AppName,
		CountryId:        1,
		DataCenterInfo: DataCenterInfo{
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Name:  "MyOwn",
		},
		HostName:         cnf.HostName(),
		Status:           "UP",
		Overriddenstatus: "UNKNOWN",
		LeaseInfo: LeaseInfo{
			RenewalIntervalInSecs: cnf.RenewalIntervalInSecs,
			DurationInSecs:        cnf.DurationInSecs,
			RegistrationTimestamp: 0,
			LastRenewalTimestamp:  0,
			RenewalTimestamp:      0,
			EvictionTimestamp:     0,
			ServiceUpTimestamp:    0,
		},
		IsCoordinatingDiscoveryServer: "false",
		Metadata: map[string]interface{}{
			"@class": "java.util.Collections$EmptyMap",
		},
		LastUpdatedTimestamp: s_ms,
		LastDirtyTimestamp:   s_ms,
		ActionType:           "ADDED",
	}

	return euk
}

// EurekaRegisterRequest Eureka注册请求结构
type EurekaRegisterRequest struct {
	Instance EurekaAppInstance `json:"instance"`
}

// EurekaRegist 注册新的服务
func EurekaRegist(ul string, auth string, e EurekaAppInstance) error {
	req := EurekaRegisterRequest{
		Instance: e,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	header := http.Header{}
	header.Set("Authorization", auth)
	header.Set("Content-type", "application/json")
	header.Set("Accept", "application/json")

	ul = strings.TrimRight(ul, "/") + "/apps/" + e.App
	resp, err := HttpPost(ul, header, body, 3)
	if err != nil {
		return errors.New("Eureka regist failed with http err: " + err.Error() + " ul:" + ul)
	}

	if resp.StatusCode != 204 {
		return errors.New("Eureka regist failed with http code " + strconv.Itoa(resp.StatusCode) + " ul:" + ul)
	}

	return nil
}

// EurekaHeartBeat 心跳续约
func EurekaHeartBeat(ul, auth, name, id string) error {
	ul = strings.TrimRight(ul, "/") + "/apps/" + name + "/" + id

	header := http.Header{}
	header.Set("Authorization", auth)

	resp, err := HttpPut(ul, header, nil, 2)
	if err != nil {
		return errors.New("Eureka heartbeat failed with http err: " + err.Error() + " ul:" + ul)
	}

	if resp.StatusCode != 200 {
		return errors.New("Eureka heartbeat failed with http code: " + strconv.Itoa(resp.StatusCode) + " ul:" + ul)
	}

	return nil
}

// EurekaGetApp 拉取应用
func EurekaGetApp(ul, auth, name string) (AppResponse, error) {
	ul = strings.TrimRight(ul, "/") + "/apps/" + name
	header := http.Header{}
	header.Set("Authorization", auth)
	header.Set("Content-type", "application/json")
	header.Set("Accept", "application/json")

	app := AppResponse{}

	resp, err := HttpGet(ul, header, nil, 3)
	if err != nil {
		return app, errors.New("Eureka app get failed with http err: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return app, errors.New("Eureka app get failed with http code " + strconv.Itoa(resp.StatusCode) + " ul:" + ul)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return app, errors.New("Eureka app get failed with read err: " + err.Error() + " ul:" + ul)
	}

	err = json.Unmarshal(body, &app)
	if err != nil {
		return app, errors.New("Eureka app get failed with json err: " + err.Error() + " ul:" + ul)
	}

	return app, nil
}

// EurekaGetAppAll 拉取所有应用
func EurekaGetAppAll(ul, auth string) (AppAllResponse, error) {
	ul = strings.TrimRight(ul, "/") + "/apps/"
	header := http.Header{}
	header.Set("Authorization", auth)
	header.Set("Content-type", "application/json")
	header.Set("Accept", "application/json")

	appAll := AppAllResponse{}

	resp, err := HttpGet(ul, header, nil, 5)
	if err != nil {
		return appAll, errors.New("Eureka app get all failed with http err: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return appAll, errors.New("Eureka app get all failed with http code " + strconv.Itoa(resp.StatusCode) + " ul:" + ul)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return appAll, errors.New("Eureka app get all failed with read err: " + err.Error() + " ul:" + ul)
	}

	err = json.Unmarshal(body, &appAll)
	if err != nil {
		return appAll, errors.New("Eureka app get all failed with json err: " + err.Error() + " ul:" + ul)
	}

	return appAll, nil
}

// EurekaDeleteApp 删除已注册的应用实例
func EurekaDeleteApp(ul, auth, name, id string) error {
	ul = strings.TrimRight(ul, "/") + "/apps/" + name + "/" + id

	header := http.Header{}
	header.Set("Authorization", auth)

	resp, err := HttpDelete(ul, header, nil, 2)
	if err != nil {
		return errors.New("Eureka detete app failed with http err: " + err.Error() + " ul:" + ul)
	}

	if resp.StatusCode != 200 {
		return errors.New("Eureka detete app failed with http code: " + strconv.Itoa(resp.StatusCode) + " ul:" + ul)
	}

	return nil
}
