# Specify the provider and access details
provider "aws" {
  region = "us-east-1"
}

variable "domain" {
  default = "dev-stages.global"
}

resource "aws_ses_domain_identity" "ses" {
  domain = "${var.domain}"
}

resource "aws_sns_topic" "ses-notifications" {
  name = "ses-notifications"
}

resource "aws_sns_topic" "ses-notifications2" {
  name = "ses-notifications2"
}

resource "aws_sns_topic" "ses-notifications3" {
  name = "ses-notifications3"
}

resource "aws_ses_identity_notification" "ses" {
  identity           = "${var.domain}"
  bounce_topic       = "${aws_sns_topic.ses-notifications.arn}"
  complaint_topic    = "${aws_sns_topic.ses-notifications2.arn}"
  delivery_topic     = "${aws_sns_topic.ses-notifications3.arn}"
  forwarding_enabled = false
}
