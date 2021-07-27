terraform {
  required_providers {
    awslexbotv2 = {
      source = "alexmammadov/awslexbotv2"
      version = "~> 0.1.6"
    }
  }
}

provider "awslexbotv2" {
}

# data "awslexbotv2_instance" "first" {
#   instance_id = "ebdfba29-ca78-41ea-b980-b328ab640a71"
# }

# output "first_order" {
#   value = data.awslexbotv2_instance.first
# }

resource "awslexbotv2_instance" "second" {
  instance_alias           = "sample-instance-2"
  identity_management_type = "CONNECT_MANAGED"
  inbound_calls_enabled    = true
  outbound_calls_enabled   = true
}

output "second_order" {
  value = awslexbotv2_instance.second
}

resource "awslexbotv2_instance_lex_bot" "second" {
  instance_id    = awslexbotv2_instance.second.instance_id
  lex_bot_region = "us-east-1"
  lex_bot_name   = "LexBotForLexBotV2"
}

# resource "awslexbotv2_instance_contact_flow" "second" {
#   instance_id = awslexbotv2_instance.second.instance_id
#   name        = "LexBotLexBotV2 contact flow"
#   type        = "CONTACT_FLOW"
#   description = "Contact flow description"
#   content     = file("contact_flow.json")
# }