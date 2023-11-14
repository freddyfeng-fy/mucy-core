package google

import (
	recaptcha "cloud.google.com/go/recaptchaenterprise/v2/apiv1"
	"cloud.google.com/go/recaptchaenterprise/v2/apiv1/recaptchaenterprisepb"
	"context"
	"fmt"
	"github.com/freddyfeng-fy/mucy-core/core"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

var (
	google *Google
)

func initReCaptchaConfig(config *Google) {
	google = config
}

func ReCaptchaAssessment(token string) bool {
	// 创建客户端
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: google.AccessToken,
	})
	client, err := recaptcha.NewClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		core.App.Log.Error("err:", zap.Any("err:", err))
		return false
	}
	defer client.Close()

	// 创建评估请求
	assessment := &recaptchaenterprisepb.Assessment{
		Event: &recaptchaenterprisepb.Event{
			Token:   token,
			SiteKey: google.ReCaptCha.RecaptchaKey,
		},
	}

	captCheReq := &recaptchaenterprisepb.CreateAssessmentRequest{
		Parent:     fmt.Sprintf("projects/%s", google.ReCaptCha.ProjectID),
		Assessment: assessment,
	}

	// 调用服务
	resp, err := client.CreateAssessment(ctx, captCheReq)
	if err != nil {
		core.App.Log.Error("err:", zap.Any("err:", err))
		return false
	}

	// 检查评估结果
	if !resp.TokenProperties.Valid {
		return false
	} else if resp.RiskAnalysis.Score < 0.5 {
		return false
	} else {
		return true
	}
}
