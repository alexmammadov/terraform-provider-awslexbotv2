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

resource "awslexbotv2_bot" "example" {
  name             = "example"
  description      = "Example bot v2"
  idle_session_ttl = 300
}

output "example_bot" {
  value = awslexbotv2_bot.example
}
