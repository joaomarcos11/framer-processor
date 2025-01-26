provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "fiap44_framer_videos_bucket" {
  bucket = "fiap44-framer-videos"
}

resource "aws_s3_bucket" "fiap44_framer_images_bucket" {
  bucket = "fiap44-framer-images"
}