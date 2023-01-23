package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	idxenrontgz "github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/index/enrontgz"
	idxtypes "github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/index/types"
	"github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/logger"
	"github.com/WitsoftGroup/index-enron-mail-zincsearch/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func getInputParams() idxenrontgz.InputParams {

	fields := os.Getenv(idxtypes.ENV_FIELDS)
	format := utils.IfThenElse(len(os.Getenv(idxtypes.ENV_FORMAT)) > 0, os.Getenv(idxtypes.ENV_FORMAT), "bulk")
	indexName := os.Getenv(idxtypes.ENV_IDXNAME)
	separator := os.Getenv(idxtypes.ENV_SEPARATOR)
	terminator := os.Getenv(idxtypes.ENV_TERMINATOR)
	tokenZinc := os.Getenv(idxtypes.ENV_TOKENZINC)
	urlZinc := os.Getenv(idxtypes.ZINC_URL)

	return idxenrontgz.InputParams{
		Fields:     &fields,
		Format:     &format,
		IndexName:  &indexName,
		Separator:  &separator,
		Terminator: &terminator,
		TokenZinc:  &tokenZinc,
		UrlZinc:    &urlZinc,
	}
}

func HandleRequest(ctx context.Context, events events.S3Event) (string, error) {
	bucketName := events.Records[0].S3.Bucket.Name
	bucketKey := events.Records[0].S3.Object.Key
	logger.Info.Println("Bucket name:", bucketName, ", bucket key:", bucketKey)

	tmpfile, err := os.Create(filepath.Join(os.TempDir(), bucketKey))
	utils.CheckError(err)
	tmpFilepath := tmpfile.Name()
	defer os.Remove(tmpFilepath)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(
		tmpfile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(bucketKey),
		},
	)
	utils.CheckError(err)

	logger.Info.Println("Downloaded:", tmpFilepath, numBytes, "bytes")

	inputParams := getInputParams()
	inputParams.InputPath = &tmpFilepath
	done := make(chan bool)

	go idxenrontgz.ProcessFile(inputParams, done)

	<-done

	return fmt.Sprintf("bucket name: %s, bucket key: %s", events.Records[0].S3.Bucket.Name, events.Records[0].S3.Object.Key), nil
}

func main() {
	logger.Info.Println("Lambda function started.")
	lambda.Start(HandleRequest)
	logger.Info.Println("Lambda function finished.")
}
