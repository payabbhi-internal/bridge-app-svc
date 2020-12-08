package helpers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/paypermint/appkit"
)

// GetS3File gets the file in given path
func GetS3File(ctxlogger appkit.AppLogger, bucketname, buckettag, fileCategory, fileName, prefix string) ([][]string, error) {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	ctxlogger.Info("Establishing connection with aws")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(bucketRegion)},
	)
	if err != nil {
		ctxlogger.Error("Unable to establish connection with aws : %v", err)
		return nil, err
	}
	ctxlogger.Info("Connection to aws established")

	ctxlogger.Info("S3 Client initialization")
	// S3 service client
	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(fileName),
	}
	ctxlogger.Info("S3 Client initialized")

	ctxlogger.Info("Started request for GetCSVFile from S3")
	response, err := svc.GetObject(input)
	if err != nil {
		ctxlogger.Error(fmt.Sprintf("Unable to get file %q from %q, %v", fileName, bucketname, err))
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error in reading file %s: %s\n", fileName, err)
	}

	data := make([][]string, 0)
	reader := csv.NewReader(bytes.NewBuffer(body))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error", err)
	}

	for value, record := range records { // for i:=0; i<len(record)
		data = append(data, record)
		fmt.Println("", records[value])
	}

	return data, nil
}
