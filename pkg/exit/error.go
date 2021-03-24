package exit

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/sealyun/cloud-kernel/pkg/dingding"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"os"
)

func ProcessError(err error) {
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
	}
	//_ = os.Stderr.Close()
	os.Exit(0)
}
