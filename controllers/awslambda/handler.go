package awslambda

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/filipeandrade6/framer-processor/domain/ports"
	"github.com/filipeandrade6/framer-processor/domain/service"
)

type Handler struct {
	strg   ports.Storager
	msgr   ports.Messager
	filer  ports.Filer
	framer ports.Framer
}

func New(strg ports.Storager, msgr ports.Messager, filer ports.Filer, framer ports.Framer) *Handler {
	return &Handler{strg: strg, msgr: msgr, filer: filer, framer: framer}
}

func (hdlr *Handler) Handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey
		size := record.S3.Object.Size

		if err := service.Process(hdlr.strg, hdlr.msgr, hdlr.filer, hdlr.framer, bucket, key, size); err != nil {
			log.Fatal(err)
		}
	}
}
