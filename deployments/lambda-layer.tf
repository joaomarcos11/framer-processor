provider "aws" {
  region = "us-east-1"
}

resource "aws_lambda_layer_version" "ffmpeg" {
  filename    = "${path.module}/ffmpeg.zip"
  layer_name  = "ffmpeg"
  description = "ffmpeg executable"

  compatible_runtimes = ["provided.al2023"]
}