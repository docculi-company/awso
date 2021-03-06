package awso

import (
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	signedUrl, err := req.Presign(144 * time.Hour)
	if err != nil {
		return ""
	}
	return signedUrl
}

////
//
// Upload file
//
////
func (o *Awso) UploadFile(usrId string, filename string, file io.Reader) error {
	bucket := "docculi-image"
	key := bucket + "/" + usrId + "/" + filename
	_, err := (*o).s3u.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

////
//
// Delete file
//
////
func (o *Awso) DeleteFile(usrId string, filename string) error {
	bucket := "docculi-image"
	key := bucket + "/" + usrId + "/" + filename
	_, err := (*o).s3svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	} else {
		err := (*o).s3svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

////
//
// Send an email notification
//
////
func (o *Awso) SendEmail(sender string, recipient string, subject string, htmlBody string, textBody string) {
	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The character encoding for the email.
	CharSet := "UTF-8"

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBody),
				},
				/*Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody),
				},*/
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := (*o).s3ses.SendEmail(input)
	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return
	}

	log.Printf("Email Sent to address: %s\n\n", recipient)
	log.Printf("Result: %s\n\n", result)
}
