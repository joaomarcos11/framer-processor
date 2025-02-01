package awslambda

import (
	"github.com/filipeandrade6/framer-processor/domain/usecases"

	"github.com/aws/aws-lambda-go/lambda"
)

func Start(strg usecases.Storage, msgr usecases.Messager) {
	lambdaHndlr := New(strg, msgr)
	lambda.Start(lambdaHndlr.Handler)
}
