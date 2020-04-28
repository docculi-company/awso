package awso

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
)

//
type Awso struct {
	s3svc *s3.S3
	s3u   *s3manager.Uploader
	s3ses *ses.SES
}

//
func New(s3svc *s3.S3, s3u *s3manager.Uploader, s3ses *ses.SES) *Awso {
	return &Awso{
		s3svc: s3svc,
		s3u:   s3u,
		s3ses: s3ses,
	}
}

////
//
// Signed url from filename
//
////
func (o *Awso) GetSignedUrl(usrId string, filename string) string {
	bucket := "docculi-image"
	key := bucket + "/" + usrId + "/" + filename
	req, _ := (*o).s3svc.GetObjectRequest(&s3.GetObjectInput{
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

////
//
// Upload file
//
////
func (o *Awso) UploadFile(usrId string, filename string, file io.Reader) {
	key := "docculi-image/" + usrId + "/" + filename
	_, err := (*o).s3u.Upload(&s3manager.UploadInput{
		Bucket: aws.String("docculi-image"),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		fmt.Printf("s3Uploader error: %s\n\n", err)
	}
}

////
//
// Delete file
//
////
func (o *Awso) DeleteFile(usrId string, filename string) {
	bucket := "docculi-image"
	key := bucket + "/" + usrId + "/" + filename
	_, err := (*o).s3svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("S3 error: unable to delete object %q from bucket %q, %v\n\n", key, bucket, err)
	}
	err = (*o).s3svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("S3 error: %s\n\n", err)
	}
}
