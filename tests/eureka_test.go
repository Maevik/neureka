package goeurekaclient

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Maevik/neureka"
)

func TestAddress(t *testing.T) {
	addr := neureka.NewAddress("TEST", "http", "127.0.0.1", "8080", "/health")
	fmt.Println("-----------服务地址测试----------------")
	fmt.Println(addr.Check())
	fmt.Println(addr.HealthUrl)
	fmt.Println(addr.Url())
	fmt.Println(addr.AppName)
}

func TestApp(t *testing.T) {
	app := neureka.NewApp("TEST")
	app.AddHost("http", "127.0.0.1", "8080", "/health")
	fmt.Println("-----------应用信息测试----------------")
	fmt.Println(app.Name)
	fmt.Println(app.Hosts)
	fmt.Println(app.GetAnUrl())
	fmt.Println(app.Hosts)
}

func TestEurekaConf(t *testing.T) {
	cnf := neureka.NewEurekaConf("default")
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090
	fmt.Println("-----------Eureka配置测试----------------")
	fmt.Println(cnf)
	fmt.Println(cnf.Id())
	fmt.Println(cnf.HostName())
}

func TestEurekaRegister(t *testing.T) {
	cnf := neureka.NewEurekaConf("default")
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090

	ins := neureka.NewEurekaAppInstance(cnf)
	err := neureka.EurekaRegist(cnf.EurekaServerAddress, cnf.Authorization, ins)
	fmt.Println("-----------Eureka注册测试----------------")
	fmt.Println("注册错误信息 >>", err)
}

func TestEurekaHeartBeat(t *testing.T) {
	cnf := neureka.NewEurekaConf("default")
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090

	fmt.Println("-----------Eureka续命测试----------------")
	fmt.Println(time.Now())
	time.Sleep(time.Second * 5)
	err := neureka.EurekaHeartBeat(cnf.EurekaServerAddress, cnf.Authorization, cnf.AppName, cnf.Id())
	fmt.Println("续命错误信息 >>", err)
}

func TestEurekaApp(t *testing.T) {
	cnf := neureka.NewEurekaConf("default")
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090
	fmt.Println("-----------Eureka拉取应用列表测试----------------")
	resp, err := neureka.EurekaGetApp(cnf.EurekaServerAddress, cnf.Authorization, "DEFAULT-EUREKA-APP0")
	fmt.Println("拉取应用错误信息 >>", err)
	res, _ := json.Marshal(resp)
	fmt.Println("拉取应用信息 >>", string(res))
}

func TestEurekaAppAll(t *testing.T) {
	cnf := neureka.NewEurekaConf("default")
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	fmt.Println("-----------Eureka拉取全量应用列表测试----------------")
	resp, err := neureka.EurekaGetAppAll(cnf.EurekaServerAddress, cnf.Authorization)
	fmt.Println("拉取全量应用错误信息 >>", err)
	res, _ := json.Marshal(resp)
	fmt.Println("拉取全量应用信息 >>", string(res))
}

func TestEurekaDeleteApp(t *testing.T) {
	cnf := neureka.NewEurekaConf("default")
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090
	err := neureka.EurekaDeleteApp(cnf.EurekaServerAddress, cnf.Authorization, cnf.AppName, cnf.Id())
	fmt.Println("-----------Eureka删除应用测试----------------")
	fmt.Println("删除应用错误信息 >>", err)
}

func TestEurekaAppsCache(t *testing.T) {
	cnf := neureka.NewEurekaConf("default")
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090

	err := neureka.Start(cnf, true)
	fmt.Println("-----------Eureka客户端服务启动测试----------------")
	fmt.Println("Eureka客户端服务启动错误信息 >>", err)

	time.Sleep(time.Second * 10)
	ul, err := neureka.GetAppUrl(cnf.EurekaName, "DEFAULT-EUREKA-APP0")
	fmt.Println("获取应用服务地址错误信息 >>", err)
	fmt.Println("获取应用服务地址 >>", ul)
	fmt.Println("获取应用服务地址 >>", ul)
	fmt.Println("获取应用服务地址 >>", ul)
}

func TestEurekaBatch(t *testing.T) {
	cnf := neureka.NewEurekaConf("default")
	cnf.AppName = "DEFAULT-EUREKA-APP9"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8190

	cnfs := []neureka.EurekaClientConfig{cnf}
	err := neureka.StartBatch(cnfs, true)
	fmt.Println("-----------Eureka客户端服务批量启动测试----------------")
	fmt.Println("Eureka客户端服务批量启动错误信息 >>", err)
	time.Sleep(time.Second * 10)
	ul, err := neureka.GetAppUrl(cnf.EurekaName, "DEFAULT-EUREKA-APP0")
	fmt.Println("获取应用服务地址错误信息 >>", err)
	fmt.Println("获取应用服务地址 >>", ul)
	fmt.Println("获取应用服务地址 >>", ul)
	fmt.Println("获取应用服务地址 >>", ul)
}
