package awso

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

//
type Awso struct {
	s3Service *s3.S3
}

//
func New(s3Service *s3.S3) *Awso {
	return &Awso{s3Service: s3Service}
}

////
//
// Signed url from filename
//
////
func (o *Awso) GetSignedUrl(usrId string, filename string) string {
	bucket := "docculi-image"
	key := bucket + "/" + usrId + "/" + filename
	req, _ := (*o).s3Service.GetObjectRequest(&s3.GetObjectInput{
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
