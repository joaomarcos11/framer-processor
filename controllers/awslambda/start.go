package awslambda

import (
	"github.com/filipeandrade6/framer-processor/domain/ports"

	"github.com/aws/aws-lambda-go/lambda"
)

func Start(strg ports.Storager, msgr ports.Messager, filer ports.Filer, framer ports.Framer) {
	lambdaHndlr := New(strg, msgr, filer, framer)
	lambda.Start(lambdaHndlr.Handler)
}
