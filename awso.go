package awso

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

////
//
// Signed url from filename
//
////
func GetSignedUrl(s3Service *s3.S3, usrId string, filename string) string {
	bucket := "docculi-image"
	key := bucket + "/" + usrId + "/" + filename
	req, _ := s3Service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	signedUrl, err := req.Presign(24 * time.Hour)
	if err != nil {
		fmt.Printf("S3 error: unable to upload %q to %q, %v\n\n", key, bucket, err)
		return ""
	}
	return signedUrl
}
