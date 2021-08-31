package provider

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lexmodelsv2"
	"github.com/aws/aws-sdk-go/service/sts"
)

func lexModelsV2Service() *lexmodelsv2.LexModelsV2 {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		SharedConfigState: session.SharedConfigEnable,
	}))

	return lexmodelsv2.New(sess)
}

func stsService() *sts.STS {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sts.New(sess)

	fmt.Println("starting to call getcalleridentity")

	reqAcc, respAcc := svc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	err := reqAcc.Send()
	if err != nil {
		fmt.Println("getcalleridentity error")
		fmt.Println(err)
	}

	fmt.Println("getcalleridentity success")
	fmt.Println(respAcc.Account)

	return svc
}
