resource "aws_sqs_queue" "fiap44_framer_sqs_status" {
  name                      = "framer-status.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
}

resource "aws_sqs_queue" "fiap44_framer_sqs_notification" {
  name                      = "framer-notification"
}