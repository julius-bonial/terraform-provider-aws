# Specify the provider and access details
provider "aws" {
  region = "us-east-1"
}

variable "domain" {
  default = "dev-stages.global"
}

resource "aws_ses_domain_identity" "foo-bar" {
  domain = "foo-bar.global"
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

resource "aws_ses_identity_notification" "bounces" {
  identity          = "${var.domain}"
  notification_type = "Bounce"
  topic_arn         = "${aws_sns_topic.ses-notifications2.arn}"
}

resource "aws_ses_identity_notification" "complaints" {
  identity          = "${var.domain}"
  notification_type = "Complaint"
  #topic_arn         = "${aws_sns_topic.ses-notifications.arn}"
}

resource "aws_ses_identity_notification" "deliveries" {
  identity          = "${var.domain}"
  notification_type = "Delivery"
  topic_arn         = "${aws_sns_topic.ses-notifications.arn}"
}


