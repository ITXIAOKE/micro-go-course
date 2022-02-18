package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/longjoy/micro-go-course/section08/user/dao"
	"github.com/longjoy/micro-go-course/section08/user/endpoint"
	"github.com/longjoy/micro-go-course/section08/user/service"
	"github.com/longjoy/micro-go-course/section08/user/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	var (
		// 服务地址和服务名
		servicePort = flag.Int("service.port", 10086, "service port")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	//err := dao.InitMysql("127.0.0.1", "3306", "root", "xuan", "user")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = redis.InitRedis("127.0.0.1", "6379", "")
	//if err != nil {
	//	log.Fatal(err)
	//}

	userService := service.MakeUserServiceImpl(&dao.UserDAOImpl{})

	userEndpoints := &endpoint.UserEndpoints{
		RegisterEndpoint: endpoint.MakeRegisterEndpoint(userService),
		LoginEndpoint:    endpoint.MakeLoginEndpoint(userService),
	}

	r := transport.MakeHttpHandler(ctx, userEndpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), r)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	log.Println(error)

}

/**
1, dao层与数据库交互，定义表，定义与数据库交互的方法
2，service主要业务实现接口
3，endpoint主要负责接收请求，并调用service包中的业务接口处理进来的请求，然后返回响应给transport
4，transport主要对外暴露项目的服务接口；比如/login接口等
5，main，整个应用的主入口，项目启动起来

curl -X POST http://localhost:10086/register -H 'content-type: application/x-www-form-urlencoded' -d 'email=aoho%40mail.com&password=aoho&username=aoho'
*/
