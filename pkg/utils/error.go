package utils

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/labring/cloud-kernel/pkg/dingding"
	"github.com/labring/cloud-kernel/pkg/logger"
)

func ProcessError(err error) error {
	//_ = os.Stderr.Close()
	//os.Exit(0)
	logger.Error(err.Error())
	return err
}

func ProcessCloudError(err error) error {
	switch err.(type) {
	case *errors.ServerError:
		e := err.(*errors.ServerError)
		logger.Error("SDK.ServerError")
		logger.Error("ErrorCode: ", e.ErrorCode())
		logger.Error("Recommend: ", e.Recommend())
		logger.Error("RequestId: ", e.RequestId())
		logger.Error("Message: ", e.Message())
		dingding.DingdingLink("离线包打包失败", "错误码:"+e.ErrorCode()+",详细信息: "+e.Message(), e.Recommend(), false)
	case *errors.ClientError:
		e := err.(*errors.ClientError)
		logger.Error("SDK.ClientError")
		logger.Error("ErrorCode: ", e.ErrorCode())
		logger.Error("Message: ", e.Message())
		dingding.DingdingText("离线包打包失败,错误码:"+e.ErrorCode()+",详细信息: "+e.Message(), false)
	default:
		logger.Error(err.Error())
		dingding.DingdingText("离线包打包失败,详细信息: "+err.Error(), false)
	}
	//_ = os.Stderr.Close()
	//os.Exit(0)
	return err
}
