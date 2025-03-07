package main

import (
	"log"

	"github.com/filipeandrade6/framer-processor/adapters/file"
	"github.com/filipeandrade6/framer-processor/adapters/framer"
	"github.com/filipeandrade6/framer-processor/adapters/message/awssqs"
	"github.com/filipeandrade6/framer-processor/adapters/storage/awss3"
	"github.com/filipeandrade6/framer-processor/controllers/awslambda"
)

func main() {
	storage, err := awss3.New()
	if err != nil {
		log.Fatalf("failed to configure storage: %s", err)
	}

	messager, err := awssqs.New()
	if err != nil {
		log.Fatalf("failed to configure messager: %s", err)
	}

	filer := file.New()
	framer := framer.New()

	awslambda.Start(storage, messager, filer, framer)
}
