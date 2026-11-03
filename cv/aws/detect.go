package aws

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

// FaceDetect detects face and return json data
func FaceDetect(reader io.Reader) (ret []byte, err error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		return ret, err
	}
	svc := rekognition.NewFromConfig(cfg)

	// Read buf to byte[]
	b, err := io.ReadAll(reader)
	if err != nil {
		return b, err
	}

	// Send the request to Rekognition
	input := &rekognition.DetectFacesInput{
		Image: &types.Image{
			Bytes: b,
		},
		Attributes: []types.Attribute{types.AttributeAll},
	}

	result, err := svc.DetectFaces(ctx, input)
	if err != nil {
		log.Print(err)
		return ret, err
	}
	return []byte(fmt.Sprintf("AWS Faces: %d", len(result.FaceDetails))), nil
}
