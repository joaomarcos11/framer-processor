package main

import (
	"log"

	"framer-proc/adapters/message/awssqs"
	"framer-proc/adapters/storage/awss3"
	"framer-proc/controllers/awslambda"
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

	awslambda.Start(storage, messager)
}
