---
layout: "aws"
page_title: "AWS: ses_identity_notification"
sidebar_current: "docs-aws-resource-ses-identity-notification"
description: |-
  Provides a resource for configuring the notifications for the SES domain
---

# ses_identity_notification

Amazon SES can notify you of the status of your emails by email or through [Amazon Simple Notification Service (Amazon SNS)](https://aws.amazon.com/de/sns/). Amazon SES supports the following three types of notifications:

* Bounces – The email is rejected by the recipient's ISP or rejected by Amazon SES because the email address is on the Amazon SES suppression list. For ISP bounces, Amazon SES reports only hard bounces and soft bounces that will no longer be retried by Amazon SES. In these cases, your recipient did not receive your email message, and Amazon SES will not try to resend it. Bounce notifications are available through email and Amazon SNS. You are notified of out-of-the-office (OOTO) messages through the same method as bounces, although they don't count toward your bounce statistics. To see an example of an OOTO bounce notification, you can use the Amazon SES mailbox simulator. For more information, see [Testing Amazon SES Email Sending](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/mailbox-simulator.html).

* Complaints – The email is accepted by the ISP and delivered to the recipient, but the recipient does not want the email and clicks a button such as "Mark as spam." If Amazon SES has a feedback loop set up with the ISP, Amazon SES will send you a complaint notification. Complaint notifications are available through email and Amazon SNS.

* Deliveries – Amazon SES successfully delivers the email to the recipient's mail server. This notification does not indicate that the actual recipient received the email because Amazon SES cannot control what happens to an email after the receiving mail server accepts it. Delivery notifications are available only through Amazon SNS.



## Example Usage

You can directly supply a topic and ARN by hand in the `topic_arn` property along with the queue ARN:

```hcl
resource "aws_ses_identity_notification" "bounces" {
  identity          = "foobar.com"
  notification_type = "Bounce"
  topic_arn         = "arn:aws:sns:us-east-1:123456:ses-bounces"
}

resource "aws_ses_identity_notification" "complaints" {
  identity          = "foobar.com"
  notification_type = "Complaint"
  topic_arn         = "arn:aws:sns:us-east-1:123456:ses-complaints"
}

resource "aws_ses_identity_notification" "deliveries" {
  identity          = "foobar.com"
  notification_type = "Delivery"
  topic_arn         = "arn:aws:sns:us-east-1:123456:ses-deliveries"
}
```


## Argument Reference

The following arguments are supported:

* `identity` - (Required) You can specify an identity by using its name or by using its Amazon Resource Name (ARN). Examples: `user@example.com`, `example.com`, `arn:aws:ses:us-east-1:123456789012:identity/example.com`.
* `notification_type` - (Required) The possible values for this are: `Bounce`, `Complaint`, `Delivery`.
* `topic_arn` - (Optional) The ARN of the SNS topic to subscribe to. If missing, forward to topic will be deleted.
