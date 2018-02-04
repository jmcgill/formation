package aws

import (
	"github.com/jmcgill/formation/core"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/hashicorp/terraform/config/configschema"
	"github.com/hashicorp/terraform/terraform"
	"sort"
	"crypto/md5"
	"strings"
	"io"
	"fmt"
)

type AwsAutoscalingNotificationImporter struct {
}

type GroupNotifications map[string][]string

// In order to identify autoscaling groups with the same notification, we generate a unique ID based on the
// types of notification for this group
func groupIdentifier(notificationTypes []string) string {
	sort.Strings(notificationTypes)
	h := md5.New()
	io.WriteString(h, strings.Join(notificationTypes, "-"))
	return string(h.Sum(nil))
}

// Lists all resources of this type
func (*AwsAutoscalingNotificationImporter) Describe(meta interface{}) ([]*core.Instance, error) {
	svc :=  meta.(*AWSClient).autoscalingconn

	notifications := make(map[string]GroupNotifications)

	existingInstances := make([]*autoscaling.NotificationConfiguration, 0)
	err := svc.DescribeNotificationConfigurationsPages(nil, func(o *autoscaling.DescribeNotificationConfigurationsOutput, lastPage bool) bool {
		existingInstances = append(existingInstances, o.NotificationConfigurations...)
		return true // continue paging
	})

	if err != nil {
		return nil, err
	}

	instances := make([]*core.Instance, 0)
	for _, existingInstance := range existingInstances {
		arn := aws.StringValue(existingInstance.TopicARN)
		if _, ok := notifications[arn]; !ok {
			notifications[arn] = make(map[string][]string)
		}

		entry := notifications[arn]
		group := aws.StringValue(existingInstance.AutoScalingGroupName)

		if _, ok := entry[group]; !ok {
			entry[group] = make([]string, 0)
		}

		entry[group] = append(entry[group], aws.StringValue(existingInstance.NotificationType))
	}

	// Each notification ARN must have its own autoscaling_notification resource in Terraform
	for arn, groups := range notifications {
		groupsByNotificationSet := make(map[string][]string)
		notificationsByNotificationSet := make(map[string][]string)

		// Group autoscaling groups by the set of notifications targeted
		for group, notifications := range groups {
			id := groupIdentifier(notifications)
			if _, ok := groupsByNotificationSet[id]; !ok {
				// May not be needed
				notificationsByNotificationSet[id] = notifications
				groupsByNotificationSet[id] = make([]string, 0)
			}
			groupsByNotificationSet[id] = append(groupsByNotificationSet[id], group)
		}

		for id, groups := range groupsByNotificationSet {
			// Emit an instance for each grouping
			instance := &core.Instance{
				// TODO(jimmy): Make each of these unique
				Name: core.Format(arn + "_" + id),
				ID: id,
				CompositeID: map[string]string{
					"group_names": strings.Join(groups, ","),
					"topic_arn": arn,
				},
			}
			instances = append(instances, instance)
		}

		//conn := meta.(*AWSClient).autoscalingconn
		//gl := convertSetToList(d.Get("group_names").(*schema.Set))
		//
		//opts := &autoscaling.DescribeNotificationConfigurationsInput{
		//	AutoScalingGroupNames: gl,
		//}
		//
		//topic := d.Get("topic_arn").(string)

	}

	return instances, nil
}

func (*AwsAutoscalingNotificationImporter) Import(in *core.Instance, meta interface{}) ([]*terraform.InstanceState, bool, error) {
	// TODO(jimmy): This might be cleaner if I use a field writer to serialize, but this would require access
	// to the underlying Schema
	group_names := strings.Split(in.CompositeID["group_names"], ",")
	attributes := make(map[string]string)

	// The AWS provider expects Group Names to be accessible as a set. Encode group names into a serialized
	// Terraform InstanceState set.
	// This has the form:
	// attribute_name.#: {{ count }}
	// attribute_name.{{ unique_key }}: {{ value }}
	// attribute_name.{{ unique_key }}: {{ value }}
	// ....
	attributes["group_names.#"] = fmt.Sprintf("%v", len(group_names))
	for i, v := range group_names {
		key := fmt.Sprintf("group_names.%v", i)
		attributes[key] = v
	}
	attributes["topic_arn"] = in.CompositeID["topic_arn"]

	state := &terraform.InstanceState{
		ID: in.ID,
		Attributes: attributes,
	}

	return []*terraform.InstanceState{
		state,
	}, false, nil
}

func (*AwsAutoscalingNotificationImporter) Clean(in *terraform.InstanceState, meta interface{}) *terraform.InstanceState {
	return in
}

func (*AwsAutoscalingNotificationImporter) AdjustSchema(in *configschema.Block) *configschema.Block {
	return in
}

//conn := meta.(*AWSClient).autoscalingconn
//gl := convertSetToList(d.Get("group_names").(*schema.Set))
//
//opts := &autoscaling.DescribeNotificationConfigurationsInput{
//AutoScalingGroupNames: gl,
//}
//
//topic := d.Get("topic_arn").(string)

// Describes which other resources this resource can reference
func (*AwsAutoscalingNotificationImporter) Links() map[string]string {
	return map[string]string{}
}
