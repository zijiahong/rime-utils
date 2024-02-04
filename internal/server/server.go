package server

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/storage"
)

type Server struct {
	app             *gin.Engine
	storage         *storage.Storage
	routerWhiteList map[string]struct{}
}

func New() *Server {
	st := storage.InitStorage()

	return &Server{
		storage: st,
		app:     gin.Default(),
		routerWhiteList: map[string]struct{}{
			"Run": {},
		},
	}
}

func (s *Server) Run() {
	s.app.Use(gin.Recovery())
	err := s.autoBindRouter()
	if err != nil {
		panic(err)
	}

	fmt.Printf("service on %s \n", config.CONFIG.Port)
	fmt.Println(s.app.Run("0.0.0.0:" + config.CONFIG.Port))
}

// autoBindRouter 自动绑定
func (s *Server) autoBindRouter() error {

	serverTypeOf := reflect.TypeOf(s)

	// 获取方法数量
	methodNum := serverTypeOf.NumMethod()
	// 遍历获取所有与公开方法
	for i := 0; i < methodNum; i++ {
		method := serverTypeOf.Method(i)
		// 排除私有方法和Run方法
		_, ok := s.routerWhiteList[method.Name]
		if !method.IsExported() || ok {
			continue
		}

		methodFundType := method.Func.Type()

		//// 检查入参格式 ////
		if methodFundType.NumIn() != 3 {
			return errors.New("the input parameter format is incorrect")
		}
		// 获取第一个参数的类型
		firstInParamType := methodFundType.In(1)

		// 判断第一个参数是否为 *gin.Context 类型
		if firstInParamType != reflect.TypeOf(&gin.Context{}) {
			return errors.New("first input parameter must be *gin.Context")
		}

		// 获取第二个参数的类型
		secondInParamType := methodFundType.In(2)
		// 判断第二个参数是否为结构体类型
		if secondInParamType.Kind() != reflect.Struct {
			return errors.New("second input parameter must be a struct")
		}

		//// 检查出参格式 ////
		if methodFundType.NumOut() != 2 {
			return errors.New("the output parameter format is incorrect")
		}

		// 获取第一个返回值的类型
		firstOutReturnType := methodFundType.Out(0)

		// 判断第一个返回值是否为切片或结构体
		switch firstOutReturnType.Kind() {
		case reflect.Slice, reflect.Struct, reflect.Ptr:
			// 第一个返回值是切片或结构体，继续检查第二个返回值
		default:
			return errors.New("first return value must be a slice or struct")
		}

		// 获取第二个返回值的类型
		secondOutReturnType := methodFundType.Out(1)

		// 判断第二个返回值是否为 error 类型
		if secondOutReturnType != reflect.TypeOf((*error)(nil)).Elem() {
			return errors.New("second return value must be of type error")
		}

		// 创建 Gin 路由处理函数
		handlerFunc := func(c *gin.Context) {
			// 创建参数值的切片
			paramValues := make([]reflect.Value, 3)
			paramValues[0] = reflect.ValueOf(s).Elem().Addr()
			paramValues[1] = reflect.ValueOf(c).Elem().Addr()
			paramValues[2] = reflect.New(methodFundType.In(2)).Elem()

			// 绑定请求参数到结构体
			if c.Request.ContentLength > 0 {
				if err := c.ShouldBind(paramValues[2].Addr().Interface()); err != nil {
					c.JSON(http.StatusOK, CommonResponse{
						Code: 500,
						Msg:  "请求参数错误",
						Data: nil,
					})
					return
				}
			}
			// 绑定uri参数
			if err := c.ShouldBindQuery(paramValues[2].Addr().Interface()); err != nil {
				c.JSON(http.StatusOK, CommonResponse{
					Code: 500,
					Msg:  "请求参数错误",
					Data: nil,
				})
				return
			}

			// toto 调用函数
			returnValues := method.Func.Call(paramValues)

			// 处理返回值
			var resultValue interface{}
			if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Ptr && returnValues[0].Elem().IsValid() {
				resultValue = returnValues[0].Elem().Interface()
			} else if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Slice {
				resultValue = returnValues[0].Interface()
			}
			errValue, ok := returnValues[1].Interface().(error)
			commonResponse := CommonResponse{}
			if ok && errValue != nil {
				// 如果有错误，使用 CommonResponse 中的 Msg 返回，并设置 code 为 500
				commonResponse.Code = 500
				commonResponse.Msg = errValue.Error()
			} else {
				// 如果无错误，将数据给 CommonResponse 中的 Data
				commonResponse.Code = 200
				commonResponse.Data = resultValue
			}

			c.JSON(http.StatusOK, commonResponse)
		}
		// 添加路由
		s.app.Handle("POST", method.Name, handlerFunc)
	}
	return nil
}

type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// getFunctionName 获取f的方法名
func (s *Server) getFunctionName(i interface{}) string {
	// 获取函数的指针
	ptr := reflect.ValueOf(i).Pointer()

	// 获取函数的反射信息
	funcName := runtime.FuncForPC(ptr).Name()

	// 使用 strings 包提取方法名
	lastDotIndex := strings.LastIndex(funcName, ".")
	if lastDotIndex != -1 {
		return strings.Replace(funcName[lastDotIndex+1:], "-fm", "", -1)
	}

	return strings.Replace(funcName, "-fm", "", -1)
}
