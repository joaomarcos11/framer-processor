provider "aws" {
  region = "us-east-1"
}

resource "aws_lambda_function" "fiap44_framer_processor" {
  filename      = "${path.module}/deployment.zip"
  function_name = "fiap44_framer_processor"
  role          = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/LabRole"
  handler       = "handleRequest"
  runtime       = "provided.al2023"
  layers        = [aws_lambda_layer_version.ffmpeg.arn]
  timeout       = 600
  memory_size   = 512
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket2"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fiap44_framer_processor.arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.fiap44_framer_videos_bucket.arn
}

resource "aws_s3_bucket_notification" "fiap44_framer_processor_notification" {
  bucket = aws_s3_bucket.fiap44_framer_videos_bucket.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.fiap44_framer_processor.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [ aws_lambda_permission.allow_bucket ]
}

data "aws_caller_identity" "current" {}

data "aws_s3_bucket" "fiap44_framer_videos_bucket" {
  bucket = "fiap44-framer-videos"
}
