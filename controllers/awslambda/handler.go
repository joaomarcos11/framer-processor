package awslambda

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/filipeandrade6/framer-processor/domain/service"
	"github.com/filipeandrade6/framer-processor/domain/usecases"
)

type Handler struct {
	strg usecases.Storage
	msgr usecases.Messager
}

func New(strg usecases.Storage, msgr usecases.Messager) *Handler {
	return &Handler{strg: strg, msgr: msgr}
}

func (hdlr *Handler) Handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey

		if err := service.Process(hdlr.strg, hdlr.msgr, bucket, key); err != nil {
			log.Fatal(err)
		}
	}
}
