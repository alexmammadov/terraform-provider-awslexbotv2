terraform {
  required_providers {
    awslexbotv2 = {
      source  = "alexmammadov/aws/awslexbotv2"
      version = "~> 0.0.1"
    }
  }
}

provider "awslexbotv2" {
}

resource "awslexbotv2_uploadurl" "example11" {
  bot_name         = "example"
  idle_session_ttl = 300
  role_arn         = "your-role-arn-here"

  file_path = "CallCenterBot.zip"
  etag      = filemd5("CallCenterBot.zip")
}

output "awslexbotv2_uploadurl" {
  value = awslexbotv2_uploadurl.example11
}

# resource "awslexbotv2_bot" "example" {
#   name             = "example"
#   description      = "Example bot v2"
#   idle_session_ttl = 300
# }

# output "example_bot" {
#   value = awslexbotv2_bot.example
# }

# resource "awslexbotv2_bot" "example4" {
#   name             = "example4"
#   description      = "Example bot v2"
#   idle_session_ttl = 300
# }


# resource "awslexbotv2_bot" "example5" {
#   name             = "example5"
#   description      = "Example bot v2"
#   idle_session_ttl = 300
# }

# resource "aws_s3_bucket_object" "object" {
#   bucket = "your_bucket_name"
#   key    = "new_object_key"
#   source = "path/to/file"

#   # The filemd5() function is available in Terraform 0.11.12 and later
#   # For Terraform 0.11.11 and earlier, use the md5() function and the file() function:
#   # etag = "${md5(file("path/to/file"))}"
#   etag = filemd5("path/to/file")
# }