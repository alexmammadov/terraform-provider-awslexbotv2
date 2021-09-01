package provider

import (
	"log"
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
	log.Println("[WARN] === sts session init")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sts.New(sess)

	log.Println("[WARN] === starting to call getcalleridentity")

	reqAcc, respAcc := svc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	err := reqAcc.Send()
	if err != nil {
		log.Println("[WARN] === getcalleridentity error")
		log.Println("[WARN] === error %s\n", err)
	}

	log.Println("[WARN] === getcalleridentity success")
	log.Printf("[WARN] === accountno %s\n", respAcc.Account)

	return svc
}
