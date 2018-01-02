package aws
//
//import (
//	"github.com/jmcgill/formation/core"
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/service/iam"
//)
//
//type AwsIamServerCertificateImporter struct {
//}
//
//// Lists all resources of this type
//func (*AwsIamServerCertificateImporter) Describe(meta interface{}) ([]*core.Instance, error) {
//	svc :=  meta.(*AWSClient).iamconn
//
//	// Add code to list resources here
//	existingInstances := make([]*iam.ServerCertificate, 0)
//	err := svc.ListServerCertificatesPages(nil, func(o *iam.ListServerCertificatesOutput, lastPage bool) bool {
//		for _, i := range o.ServerCertificateMetadataList {
//			existingInstances = append(existingInstances, i)
//		}
//		return true
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	instances := make([]*core.Instance, len(existingInstances))
//	for i, existingInstance := range existingInstances {
//		instances[i] = &core.Instance{
//			Name: core.Format(aws.StringValue(existingInstance.RoleName)),
//			ID:   aws.StringValue(existingInstance.RoleName),
//		}
//	}
//
//	return instances, nil
//}
//
//// Describes which other resources this resource can reference
//func (*AwsIamServerCertificateImporter) Links() map[string]string {
//	return map[string]string{
//	}
//}