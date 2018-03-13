package aws

import (
	"fmt"
	//"github.com/aws/aws-sdk-go/aws/awsutil"
	"log"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsSesNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSesNotificationSet,
		Read:   resourceAwsSesNotificationRead,
		Update: resourceAwsSesNotificationSet,
		Delete: resourceAwsSesNotificationDelete,

		Schema: map[string]*schema.Schema{
			"identity": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIdentity,
			},

			"bounce_topic": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateArn,
			},

			"complaint_topic": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateArn,
			},

			"delivery_topic": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateArn,
			},

			"forwarding_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAwsSesNotificationSet(d *schema.ResourceData, meta interface{}) error {
	topics := []string{ses.NotificationTypeBounce, ses.NotificationTypeComplaint, ses.NotificationTypeDelivery}
	identity := d.Get("identity").(string)
	conn := meta.(*AWSClient).sesConn

	d.Partial(true)
	d.SetId(identity)

	for _, topic := range topics {
		schema_name := strings.ToLower(topic) + "_topic"

		if d.HasChange(schema_name) {
			sns_topic := aws.String("")

			if v, ok := d.GetOk(schema_name); ok {
				sns_topic = aws.String(v.(string))
			} else {
				sns_topic = nil
			}

			setOpts := &ses.SetIdentityNotificationTopicInput{
				Identity:         aws.String(identity),
				NotificationType: aws.String(topic),
				SnsTopic:         sns_topic,
			}

			log.Printf("[DEBUG] Setting SES Identity Notification: %+v", setOpts)

			_, err := conn.SetIdentityNotificationTopic(setOpts)

			if err != nil {
				return fmt.Errorf("Error setting SES Identity Notification: %s", err)
			}

			d.SetPartial(schema_name)
		}
	}

	d.Partial(false)

	return nil
}

func resourceAwsSesNotificationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn
	identity := d.Id()
	topics := []string{ses.NotificationTypeBounce, ses.NotificationTypeComplaint, ses.NotificationTypeDelivery}

	getOpts := &ses.GetIdentityNotificationAttributesInput{
		Identities: []*string{aws.String(identity)},
	}

	log.Printf("[DEBUG] Reading SES Identity Notification Attributes: %#v", getOpts)

	response, err := conn.GetIdentityNotificationAttributes(getOpts)

	if err != nil {
		return fmt.Errorf("Error reading SES Identity Notification: %s", err)
	}
	notificationAttributes := response.NotificationAttributes[identity]
	r := reflect.ValueOf(notificationAttributes)

	log.Printf("[DEBUG] notificationAttributes: %+v", notificationAttributes)

	for _, topic := range topics {
		attribute_name := topic + "Topic"
		schema_name := strings.ToLower(topic) + "_topic"
		topic_arn_ptr := *(reflect.Indirect(r).FieldByName(attribute_name).Addr().Interface().(**string))

		if topic_arn_ptr == nil {
			if err := d.Set(schema_name, nil); err != nil {
				return err
			}
			continue
		}

		topic_arn := *topic_arn_ptr

		if err := d.Set(schema_name, topic_arn); err != nil {
			return err
		}
	}

	return nil
}

func resourceAwsSesNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn
	notification := d.Get("notification_type").(string)
	identity := d.Get("identity").(string)

	setOpts := &ses.SetIdentityNotificationTopicInput{
		Identity:         aws.String(identity),
		NotificationType: aws.String(notification),
		SnsTopic:         nil,
	}

	log.Printf("[DEBUG] Deleting SES Identity Notification: %#v", setOpts)

	if _, err := conn.SetIdentityNotificationTopic(setOpts); err != nil {
		return fmt.Errorf("Error deleting SES Identity Notification: %s", err)
	}

	return resourceAwsSesNotificationRead(d, meta)
}

func validateIdentity(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	if value != "" {
		return
	}

	errors = append(errors, fmt.Errorf("%q must not be empty", k))
	return
}
