provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "framer_videos" {
  bucket = "fiap44-framer-videos"
}

resource "aws_s3_bucket" "framer_images" {
  bucket = "fiap44-framer-images"
}