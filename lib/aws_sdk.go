package lib

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"identity-token-relayer/config"
	"sync"
)

var (
	awsSdkInit = sync.Once{}
	awsSession *session.Session
	awsKmsInit = sync.Once{}
	kmsClient  *kms.KMS
)

func GetAwsSession() *session.Session {
	awsSdkInit.Do(func() {
		// init AWS sdk
		sessionOption := session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Profile:           config.Get().Aws.Profile,
		}
		awsRegion := config.Get().Aws.Region
		sessionOption.Config.Region = &awsRegion

		awsSession = session.Must(session.NewSessionWithOptions(sessionOption))
	})
	return awsSession
}

func GetKmsClient() *kms.KMS {
	awsKmsInit.Do(func() {
		kmsClient = kms.New(GetAwsSession())
	})
	return kmsClient
}
