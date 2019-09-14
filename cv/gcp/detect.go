package gcp

import (
	"context"
	"fmt"
	"io"
	"log"

	vision "cloud.google.com/go/vision/apiv1"
)

/*
2019/06/05 18:31:57 Failed to detect labels: rpc error: code = PermissionDenied desc = Cloud Vision API has not been used in project 1053435683258 before or it is disabled. Enable it by visiting https://console.developers.google.com/apis/api/vision.googleapis.com/overview?project=1053435683258 then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry.
*/
// export GOOGLE_APPLICATION_CREDENTIALS="[PATH]"
// https://cloud.google.com/vision/docs/face-tutorial?hl=zh-tw
// FaceDetect detects face and return json data
func FaceDetect(reader io.Reader) (ret []byte, err error) {
	ctx := context.Background()
	//	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(jsonPath))
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return ret, err
	}
	defer client.Close()
	image, err := vision.NewImageFromReader(reader)
	if err != nil {
		log.Printf("Failed to create image: %v", err)
		return ret, err
	}

	faces, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		log.Printf("Failed to detect labels: %v", err)
		return ret, err
	}

	return []byte(fmt.Sprintf("GCP Faces: %d", len(faces))), nil
}
