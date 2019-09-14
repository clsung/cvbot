package aws

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

// FaceDetect detects face and return json data
func FaceDetect(reader io.Reader) (ret []byte, err error) {
	sess := session.New(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	svc := rekognition.New(sess)

	// Read buf to byte[]
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return b, err
	}

	// Send the request to Rekognition
	input := &rekognition.DetectFacesInput{
		Image: &rekognition.Image{
			Bytes: b,
		},
		Attributes: []*string{aws.String("ALL")},
	}

	result, err := svc.DetectFaces(input)
	if err != nil {
		log.Print(err)
		return ret, err
	}
	return []byte(fmt.Sprintf("AWS Faces: %d", len(result.FaceDetails))), nil
}
