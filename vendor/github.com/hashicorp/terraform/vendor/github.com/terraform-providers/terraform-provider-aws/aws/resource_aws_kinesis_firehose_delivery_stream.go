package aws

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func cloudWatchLoggingOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"log_group_name": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"log_stream_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func s3ConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bucket_arn": {
					Type:     schema.TypeString,
					Required: true,
				},

				"buffer_size": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  5,
				},

				"buffer_interval": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  300,
				},

				"compression_format": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "UNCOMPRESSED",
				},

				"kms_key_arn": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validateArn,
				},

				"role_arn": {
					Type:     schema.TypeString,
					Required: true,
				},

				"prefix": {
					Type:     schema.TypeString,
					Optional: true,
				},

				"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),
			},
		},
	}
}

func processingConfigurationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"processors": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"parameters": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"parameter_name": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
												value := v.(string)
												if value != "LambdaArn" && value != "NumberOfRetries" {
													errors = append(errors, fmt.Errorf(
														"%q must be one of 'LambdaArn', 'NumberOfRetries'", k))
												}
												return
											},
										},
										"parameter_value": {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
												value := v.(string)
												if len(value) < 1 || len(value) > 512 {
													errors = append(errors, fmt.Errorf(
														"%q must be at least one character long and at most 512 characters long", k))
												}
												return
											},
										},
									},
								},
							},
							"type": {
								Type:     schema.TypeString,
								Required: true,
								ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
									value := v.(string)
									if value != "Lambda" {
										errors = append(errors, fmt.Errorf(
											"%q must be 'Lambda'", k))
									}
									return
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAwsKinesisFirehoseDeliveryStream() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsKinesisFirehoseDeliveryStreamCreate,
		Read:   resourceAwsKinesisFirehoseDeliveryStreamRead,
		Update: resourceAwsKinesisFirehoseDeliveryStreamUpdate,
		Delete: resourceAwsKinesisFirehoseDeliveryStreamDelete,

		SchemaVersion: 1,
		MigrateState:  resourceAwsKinesisFirehoseMigrateState,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 64 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be longer than 64 characters", k))
					}
					return
				},
			},

			"kinesis_source_configuration": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kinesis_stream_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateArn,
						},

						"role_arn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateArn,
						},
					},
				},
			},

			"destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					value := v.(string)
					return strings.ToLower(value)
				},
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "s3" && value != "extended_s3" && value != "redshift" && value != "elasticsearch" {
						errors = append(errors, fmt.Errorf(
							"%q must be one of 's3', 'extended_s3', 'redshift', 'elasticsearch'", k))
					}
					return
				},
			},

			"s3_configuration": s3ConfigurationSchema(),

			"extended_s3_configuration": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"s3_configuration"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_arn": {
							Type:     schema.TypeString,
							Required: true,
						},

						"buffer_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},

						"buffer_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  300,
						},

						"compression_format": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "UNCOMPRESSED",
						},

						"kms_key_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateArn,
						},

						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},

						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),

						"processing_configuration": processingConfigurationSchema(),
					},
				},
			},

			"redshift_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_jdbcurl": {
							Type:     schema.TypeString,
							Required: true,
						},

						"username": {
							Type:     schema.TypeString,
							Required: true,
						},

						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},

						"s3_backup_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Disabled",
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if value != "Disabled" && value != "Enabled" {
									errors = append(errors, fmt.Errorf(
										"%q must be one of 'Disabled', 'Enabled'", k))
								}
								return
							},
						},

						"s3_backup_configuration": s3ConfigurationSchema(),

						"retry_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3600,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								if value < 0 || value > 7200 {
									errors = append(errors, fmt.Errorf(
										"%q must be in the range from 0 to 7200 seconds.", k))
								}
								return
							},
						},

						"copy_options": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"data_table_columns": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"data_table_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),
					},
				},
			},

			"elasticsearch_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"buffering_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  300,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								if value < 60 || value > 900 {
									errors = append(errors, fmt.Errorf(
										"%q must be in the range from 60 to 900 seconds.", k))
								}
								return
							},
						},

						"buffering_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								if value < 1 || value > 100 {
									errors = append(errors, fmt.Errorf(
										"%q must be in the range from 1 to 100 MB.", k))
								}
								return
							},
						},

						"domain_arn": {
							Type:     schema.TypeString,
							Required: true,
						},

						"index_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"index_rotation_period": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "OneDay",
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if value != "NoRotation" && value != "OneHour" && value != "OneDay" && value != "OneWeek" && value != "OneMonth" {
									errors = append(errors, fmt.Errorf(
										"%q must be one of 'NoRotation', 'OneHour', 'OneDay', 'OneWeek', 'OneMonth'", k))
								}
								return
							},
						},

						"retry_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  300,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								if value < 0 || value > 7200 {
									errors = append(errors, fmt.Errorf(
										"%q must be in the range from 0 to 7200 seconds.", k))
								}
								return
							},
						},

						"role_arn": {
							Type:     schema.TypeString,
							Required: true,
						},

						"s3_backup_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "FailedDocumentsOnly",
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if value != "FailedDocumentsOnly" && value != "AllDocuments" {
									errors = append(errors, fmt.Errorf(
										"%q must be one of 'FailedDocumentsOnly', 'AllDocuments'", k))
								}
								return
							},
						},

						"type_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								if len(value) > 100 {
									errors = append(errors, fmt.Errorf(
										"%q cannot be longer than 100 characters", k))
								}
								return
							},
						},

						"cloudwatch_logging_options": cloudWatchLoggingOptionsSchema(),
					},
				},
			},

			"arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"destination_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func createSourceConfig(source map[string]interface{}) *firehose.KinesisStreamSourceConfiguration {

	configuration := &firehose.KinesisStreamSourceConfiguration{
		KinesisStreamARN: aws.String(source["kinesis_stream_arn"].(string)),
		RoleARN:          aws.String(source["role_arn"].(string)),
	}

	return configuration
}

func createS3Config(d *schema.ResourceData) *firehose.S3DestinationConfiguration {
	s3 := d.Get("s3_configuration").([]interface{})[0].(map[string]interface{})

	configuration := &firehose.S3DestinationConfiguration{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64(int64(s3["buffer_interval"].(int))),
			SizeInMBs:         aws.Int64(int64(s3["buffer_size"].(int))),
		},
		Prefix:                  extractPrefixConfiguration(s3),
		CompressionFormat:       aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration: extractEncryptionConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func expandS3BackupConfig(d map[string]interface{}) *firehose.S3DestinationConfiguration {
	config := d["s3_backup_configuration"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	s3 := config[0].(map[string]interface{})

	configuration := &firehose.S3DestinationConfiguration{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64(int64(s3["buffer_interval"].(int))),
			SizeInMBs:         aws.Int64(int64(s3["buffer_size"].(int))),
		},
		Prefix:                  extractPrefixConfiguration(s3),
		CompressionFormat:       aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration: extractEncryptionConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func createExtendedS3Config(d *schema.ResourceData) *firehose.ExtendedS3DestinationConfiguration {
	s3 := d.Get("extended_s3_configuration").([]interface{})[0].(map[string]interface{})

	configuration := &firehose.ExtendedS3DestinationConfiguration{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64(int64(s3["buffer_interval"].(int))),
			SizeInMBs:         aws.Int64(int64(s3["buffer_size"].(int))),
		},
		Prefix:                  extractPrefixConfiguration(s3),
		CompressionFormat:       aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration: extractEncryptionConfiguration(s3),
	}

	if _, ok := s3["processing_configuration"]; ok {
		configuration.ProcessingConfiguration = extractProcessingConfiguration(s3)
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func updateS3Config(d *schema.ResourceData) *firehose.S3DestinationUpdate {
	s3 := d.Get("s3_configuration").([]interface{})[0].(map[string]interface{})

	configuration := &firehose.S3DestinationUpdate{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64((int64)(s3["buffer_interval"].(int))),
			SizeInMBs:         aws.Int64((int64)(s3["buffer_size"].(int))),
		},
		Prefix:                   extractPrefixConfiguration(s3),
		CompressionFormat:        aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration:  extractEncryptionConfiguration(s3),
		CloudWatchLoggingOptions: extractCloudWatchLoggingConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func updateS3BackupConfig(d map[string]interface{}) *firehose.S3DestinationUpdate {
	config := d["s3_backup_configuration"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	s3 := config[0].(map[string]interface{})

	configuration := &firehose.S3DestinationUpdate{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64((int64)(s3["buffer_interval"].(int))),
			SizeInMBs:         aws.Int64((int64)(s3["buffer_size"].(int))),
		},
		Prefix:                   extractPrefixConfiguration(s3),
		CompressionFormat:        aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration:  extractEncryptionConfiguration(s3),
		CloudWatchLoggingOptions: extractCloudWatchLoggingConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func updateExtendedS3Config(d *schema.ResourceData) *firehose.ExtendedS3DestinationUpdate {
	s3 := d.Get("extended_s3_configuration").([]interface{})[0].(map[string]interface{})

	configuration := &firehose.ExtendedS3DestinationUpdate{
		BucketARN: aws.String(s3["bucket_arn"].(string)),
		RoleARN:   aws.String(s3["role_arn"].(string)),
		BufferingHints: &firehose.BufferingHints{
			IntervalInSeconds: aws.Int64((int64)(s3["buffer_interval"].(int))),
			SizeInMBs:         aws.Int64((int64)(s3["buffer_size"].(int))),
		},
		Prefix:                   extractPrefixConfiguration(s3),
		CompressionFormat:        aws.String(s3["compression_format"].(string)),
		EncryptionConfiguration:  extractEncryptionConfiguration(s3),
		CloudWatchLoggingOptions: extractCloudWatchLoggingConfiguration(s3),
		ProcessingConfiguration:  extractProcessingConfiguration(s3),
	}

	if _, ok := s3["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(s3)
	}

	return configuration
}

func extractProcessingConfiguration(s3 map[string]interface{}) *firehose.ProcessingConfiguration {
	config := s3["processing_configuration"].([]interface{})
	if len(config) == 0 {
		return nil
	}

	processingConfiguration := config[0].(map[string]interface{})

	return &firehose.ProcessingConfiguration{
		Enabled:    aws.Bool(processingConfiguration["enabled"].(bool)),
		Processors: extractProcessors(processingConfiguration["processors"].([]interface{})),
	}
}

func extractProcessors(processingConfigurationProcessors []interface{}) []*firehose.Processor {
	processors := []*firehose.Processor{}

	for _, processor := range processingConfigurationProcessors {
		processors = append(processors, extractProcessor(processor.(map[string]interface{})))
	}

	return processors
}

func extractProcessor(processingConfigurationProcessor map[string]interface{}) *firehose.Processor {
	return &firehose.Processor{
		Type:       aws.String(processingConfigurationProcessor["type"].(string)),
		Parameters: extractProcessorParameters(processingConfigurationProcessor["parameters"].([]interface{})),
	}
}

func extractProcessorParameters(processorParameters []interface{}) []*firehose.ProcessorParameter {
	parameters := []*firehose.ProcessorParameter{}

	for _, attr := range processorParameters {
		parameters = append(parameters, extractProcessorParameter(attr.(map[string]interface{})))
	}

	return parameters
}

func extractProcessorParameter(processorParameter map[string]interface{}) *firehose.ProcessorParameter {
	parameter := &firehose.ProcessorParameter{
		ParameterName:  aws.String(processorParameter["parameter_name"].(string)),
		ParameterValue: aws.String(processorParameter["parameter_value"].(string)),
	}

	return parameter
}

func extractEncryptionConfiguration(s3 map[string]interface{}) *firehose.EncryptionConfiguration {
	if key, ok := s3["kms_key_arn"]; ok && len(key.(string)) > 0 {
		return &firehose.EncryptionConfiguration{
			KMSEncryptionConfig: &firehose.KMSEncryptionConfig{
				AWSKMSKeyARN: aws.String(key.(string)),
			},
		}
	}

	return &firehose.EncryptionConfiguration{
		NoEncryptionConfig: aws.String("NoEncryption"),
	}
}

func extractCloudWatchLoggingConfiguration(s3 map[string]interface{}) *firehose.CloudWatchLoggingOptions {
	config := s3["cloudwatch_logging_options"].(*schema.Set).List()
	if len(config) == 0 {
		return nil
	}

	loggingConfig := config[0].(map[string]interface{})
	loggingOptions := &firehose.CloudWatchLoggingOptions{
		Enabled: aws.Bool(loggingConfig["enabled"].(bool)),
	}

	if v, ok := loggingConfig["log_group_name"]; ok {
		loggingOptions.LogGroupName = aws.String(v.(string))
	}

	if v, ok := loggingConfig["log_stream_name"]; ok {
		loggingOptions.LogStreamName = aws.String(v.(string))
	}

	return loggingOptions

}

func extractPrefixConfiguration(s3 map[string]interface{}) *string {
	if v, ok := s3["prefix"]; ok {
		return aws.String(v.(string))
	}

	return nil
}

func createRedshiftConfig(d *schema.ResourceData, s3Config *firehose.S3DestinationConfiguration) (*firehose.RedshiftDestinationConfiguration, error) {
	redshiftRaw, ok := d.GetOk("redshift_configuration")
	if !ok {
		return nil, fmt.Errorf("[ERR] Error loading Redshift Configuration for Kinesis Firehose: redshift_configuration not found")
	}
	rl := redshiftRaw.([]interface{})

	redshift := rl[0].(map[string]interface{})

	configuration := &firehose.RedshiftDestinationConfiguration{
		ClusterJDBCURL:  aws.String(redshift["cluster_jdbcurl"].(string)),
		RetryOptions:    extractRedshiftRetryOptions(redshift),
		Password:        aws.String(redshift["password"].(string)),
		Username:        aws.String(redshift["username"].(string)),
		RoleARN:         aws.String(redshift["role_arn"].(string)),
		CopyCommand:     extractCopyCommandConfiguration(redshift),
		S3Configuration: s3Config,
	}

	if _, ok := redshift["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(redshift)
	}
	if s3BackupMode, ok := redshift["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
		configuration.S3BackupConfiguration = expandS3BackupConfig(d.Get("redshift_configuration").([]interface{})[0].(map[string]interface{}))
	}

	return configuration, nil
}

func updateRedshiftConfig(d *schema.ResourceData, s3Update *firehose.S3DestinationUpdate) (*firehose.RedshiftDestinationUpdate, error) {
	redshiftRaw, ok := d.GetOk("redshift_configuration")
	if !ok {
		return nil, fmt.Errorf("[ERR] Error loading Redshift Configuration for Kinesis Firehose: redshift_configuration not found")
	}
	rl := redshiftRaw.([]interface{})

	redshift := rl[0].(map[string]interface{})

	configuration := &firehose.RedshiftDestinationUpdate{
		ClusterJDBCURL: aws.String(redshift["cluster_jdbcurl"].(string)),
		RetryOptions:   extractRedshiftRetryOptions(redshift),
		Password:       aws.String(redshift["password"].(string)),
		Username:       aws.String(redshift["username"].(string)),
		RoleARN:        aws.String(redshift["role_arn"].(string)),
		CopyCommand:    extractCopyCommandConfiguration(redshift),
		S3Update:       s3Update,
	}

	if _, ok := redshift["cloudwatch_logging_options"]; ok {
		configuration.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(redshift)
	}
	if s3BackupMode, ok := redshift["s3_backup_mode"]; ok {
		configuration.S3BackupMode = aws.String(s3BackupMode.(string))
		configuration.S3BackupUpdate = updateS3BackupConfig(d.Get("redshift_configuration").([]interface{})[0].(map[string]interface{}))
	}

	return configuration, nil
}

func createElasticsearchConfig(d *schema.ResourceData, s3Config *firehose.S3DestinationConfiguration) (*firehose.ElasticsearchDestinationConfiguration, error) {
	esConfig, ok := d.GetOk("elasticsearch_configuration")
	if !ok {
		return nil, fmt.Errorf("[ERR] Error loading Elasticsearch Configuration for Kinesis Firehose: elasticsearch_configuration not found")
	}
	esList := esConfig.([]interface{})

	es := esList[0].(map[string]interface{})

	config := &firehose.ElasticsearchDestinationConfiguration{
		BufferingHints:  extractBufferingHints(es),
		DomainARN:       aws.String(es["domain_arn"].(string)),
		IndexName:       aws.String(es["index_name"].(string)),
		RetryOptions:    extractElasticSearchRetryOptions(es),
		RoleARN:         aws.String(es["role_arn"].(string)),
		TypeName:        aws.String(es["type_name"].(string)),
		S3Configuration: s3Config,
	}

	if _, ok := es["cloudwatch_logging_options"]; ok {
		config.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(es)
	}

	if indexRotationPeriod, ok := es["index_rotation_period"]; ok {
		config.IndexRotationPeriod = aws.String(indexRotationPeriod.(string))
	}
	if s3BackupMode, ok := es["s3_backup_mode"]; ok {
		config.S3BackupMode = aws.String(s3BackupMode.(string))
	}

	return config, nil
}

func updateElasticsearchConfig(d *schema.ResourceData, s3Update *firehose.S3DestinationUpdate) (*firehose.ElasticsearchDestinationUpdate, error) {
	esConfig, ok := d.GetOk("elasticsearch_configuration")
	if !ok {
		return nil, fmt.Errorf("[ERR] Error loading Elasticsearch Configuration for Kinesis Firehose: elasticsearch_configuration not found")
	}
	esList := esConfig.([]interface{})

	es := esList[0].(map[string]interface{})

	update := &firehose.ElasticsearchDestinationUpdate{
		BufferingHints: extractBufferingHints(es),
		DomainARN:      aws.String(es["domain_arn"].(string)),
		IndexName:      aws.String(es["index_name"].(string)),
		RetryOptions:   extractElasticSearchRetryOptions(es),
		RoleARN:        aws.String(es["role_arn"].(string)),
		TypeName:       aws.String(es["type_name"].(string)),
		S3Update:       s3Update,
	}

	if _, ok := es["cloudwatch_logging_options"]; ok {
		update.CloudWatchLoggingOptions = extractCloudWatchLoggingConfiguration(es)
	}

	if indexRotationPeriod, ok := es["index_rotation_period"]; ok {
		update.IndexRotationPeriod = aws.String(indexRotationPeriod.(string))
	}

	return update, nil
}

func extractBufferingHints(es map[string]interface{}) *firehose.ElasticsearchBufferingHints {
	bufferingHints := &firehose.ElasticsearchBufferingHints{}

	if bufferingInterval, ok := es["buffering_interval"].(int); ok {
		bufferingHints.IntervalInSeconds = aws.Int64(int64(bufferingInterval))
	}
	if bufferingSize, ok := es["buffering_size"].(int); ok {
		bufferingHints.SizeInMBs = aws.Int64(int64(bufferingSize))
	}

	return bufferingHints
}

func extractElasticSearchRetryOptions(es map[string]interface{}) *firehose.ElasticsearchRetryOptions {
	retryOptions := &firehose.ElasticsearchRetryOptions{}

	if retryDuration, ok := es["retry_duration"].(int); ok {
		retryOptions.DurationInSeconds = aws.Int64(int64(retryDuration))
	}

	return retryOptions
}

func extractRedshiftRetryOptions(redshift map[string]interface{}) *firehose.RedshiftRetryOptions {
	retryOptions := &firehose.RedshiftRetryOptions{}

	if retryDuration, ok := redshift["retry_duration"].(int); ok {
		retryOptions.DurationInSeconds = aws.Int64(int64(retryDuration))
	}

	return retryOptions
}

func extractCopyCommandConfiguration(redshift map[string]interface{}) *firehose.CopyCommand {
	cmd := &firehose.CopyCommand{
		DataTableName: aws.String(redshift["data_table_name"].(string)),
	}
	if copyOptions, ok := redshift["copy_options"]; ok {
		cmd.CopyOptions = aws.String(copyOptions.(string))
	}
	if columns, ok := redshift["data_table_columns"]; ok {
		cmd.DataTableColumns = aws.String(columns.(string))
	}

	return cmd
}

func resourceAwsKinesisFirehoseDeliveryStreamCreate(d *schema.ResourceData, meta interface{}) error {
	validateError := validateAwsKinesisFirehoseSchema(d)

	if validateError != nil {
		return validateError
	}

	conn := meta.(*AWSClient).firehoseconn

	sn := d.Get("name").(string)

	createInput := &firehose.CreateDeliveryStreamInput{
		DeliveryStreamName: aws.String(sn),
	}

	if v, ok := d.GetOk("kinesis_source_configuration"); ok {
		sourceConfig := createSourceConfig(v.([]interface{})[0].(map[string]interface{}))
		createInput.KinesisStreamSourceConfiguration = sourceConfig
		createInput.DeliveryStreamType = aws.String(firehose.DeliveryStreamTypeKinesisStreamAsSource)
	} else {
		createInput.DeliveryStreamType = aws.String(firehose.DeliveryStreamTypeDirectPut)
	}

	if d.Get("destination").(string) == "extended_s3" {
		extendedS3Config := createExtendedS3Config(d)
		createInput.ExtendedS3DestinationConfiguration = extendedS3Config
	} else {
		s3Config := createS3Config(d)

		if d.Get("destination").(string) == "s3" {
			createInput.S3DestinationConfiguration = s3Config
		} else if d.Get("destination").(string) == "elasticsearch" {
			esConfig, err := createElasticsearchConfig(d, s3Config)
			if err != nil {
				return err
			}
			createInput.ElasticsearchDestinationConfiguration = esConfig
		} else {
			rc, err := createRedshiftConfig(d, s3Config)
			if err != nil {
				return err
			}
			createInput.RedshiftDestinationConfiguration = rc
		}
	}

	var lastError error
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := conn.CreateDeliveryStream(createInput)
		if err != nil {
			log.Printf("[DEBUG] Error creating Firehose Delivery Stream: %s", err)
			lastError = err

			if awsErr, ok := err.(awserr.Error); ok {
				// IAM roles can take ~10 seconds to propagate in AWS:
				// http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html#launch-instance-with-role-console
				if awsErr.Code() == "InvalidArgumentException" && strings.Contains(awsErr.Message(), "Firehose is unable to assume role") {
					log.Printf("[DEBUG] Firehose could not assume role referenced, retrying...")
					return resource.RetryableError(awsErr)
				}
			}
			// Not retryable
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		if awsErr, ok := lastError.(awserr.Error); ok {
			return fmt.Errorf("[WARN] Error creating Kinesis Firehose Delivery Stream: %s", awsErr.Error())
		}
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    firehoseStreamStateRefreshFunc(conn, sn),
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	firehoseStream, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Kinesis Stream (%s) to become active: %s",
			sn, err)
	}

	s := firehoseStream.(*firehose.DeliveryStreamDescription)
	d.SetId(*s.DeliveryStreamARN)
	d.Set("arn", s.DeliveryStreamARN)

	return resourceAwsKinesisFirehoseDeliveryStreamRead(d, meta)
}

func validateAwsKinesisFirehoseSchema(d *schema.ResourceData) error {

	_, s3Exists := d.GetOk("s3_configuration")
	_, extendedS3Exists := d.GetOk("extended_s3_configuration")

	if d.Get("destination").(string) == "extended_s3" {
		if !extendedS3Exists {
			return fmt.Errorf(
				"When destination is 'extended_s3', extended_s3_configuration is required",
			)
		} else if s3Exists {
			return fmt.Errorf(
				"When destination is 'extended_s3', s3_configuration must not be set",
			)
		}
	} else {
		if !s3Exists {
			return fmt.Errorf(
				"When destination is %s, s3_configuration is required",
				d.Get("destination").(string),
			)
		} else if extendedS3Exists {
			return fmt.Errorf(
				"extended_s3_configuration can only be used when destination is 'extended_s3'",
			)
		}
	}

	return nil
}

func resourceAwsKinesisFirehoseDeliveryStreamUpdate(d *schema.ResourceData, meta interface{}) error {
	validateError := validateAwsKinesisFirehoseSchema(d)

	if validateError != nil {
		return validateError
	}

	conn := meta.(*AWSClient).firehoseconn

	sn := d.Get("name").(string)

	updateInput := &firehose.UpdateDestinationInput{
		DeliveryStreamName:             aws.String(sn),
		CurrentDeliveryStreamVersionId: aws.String(d.Get("version_id").(string)),
		DestinationId:                  aws.String(d.Get("destination_id").(string)),
	}

	if d.Get("destination").(string) == "extended_s3" {
		extendedS3Config := updateExtendedS3Config(d)
		updateInput.ExtendedS3DestinationUpdate = extendedS3Config
	} else {
		s3Config := updateS3Config(d)

		if d.Get("destination").(string) == "s3" {
			updateInput.S3DestinationUpdate = s3Config
		} else if d.Get("destination").(string) == "elasticsearch" {
			esUpdate, err := updateElasticsearchConfig(d, s3Config)
			if err != nil {
				return err
			}
			updateInput.ElasticsearchDestinationUpdate = esUpdate
		} else {
			rc, err := updateRedshiftConfig(d, s3Config)
			if err != nil {
				return err
			}
			updateInput.RedshiftDestinationUpdate = rc
		}
	}

	_, err := conn.UpdateDestination(updateInput)
	if err != nil {
		return fmt.Errorf(
			"Error Updating Kinesis Firehose Delivery Stream: \"%s\"\n%s",
			sn, err)
	}

	return resourceAwsKinesisFirehoseDeliveryStreamRead(d, meta)
}

func resourceAwsKinesisFirehoseDeliveryStreamRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).firehoseconn

	resp, err := conn.DescribeDeliveryStream(&firehose.DescribeDeliveryStreamInput{
		DeliveryStreamName: aws.String(d.Get("name").(string)),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[WARN] Error reading Kinesis Firehose Delivery Stream: \"%s\", code: \"%s\"", awsErr.Message(), awsErr.Code())
		}
		return err
	}

	s := resp.DeliveryStreamDescription
	d.Set("version_id", s.VersionId)
	d.Set("arn", *s.DeliveryStreamARN)
	if len(s.Destinations) > 0 {
		destination := s.Destinations[0]
		d.Set("destination_id", *destination.DestinationId)
	}

	return nil
}

func resourceAwsKinesisFirehoseDeliveryStreamDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).firehoseconn

	sn := d.Get("name").(string)
	_, err := conn.DeleteDeliveryStream(&firehose.DeleteDeliveryStreamInput{
		DeliveryStreamName: aws.String(sn),
	})

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DESTROYED"},
		Refresh:    firehoseStreamStateRefreshFunc(conn, sn),
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Delivery Stream (%s) to be destroyed: %s",
			sn, err)
	}

	d.SetId("")
	return nil
}

func firehoseStreamStateRefreshFunc(conn *firehose.Firehose, sn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		describeOpts := &firehose.DescribeDeliveryStreamInput{
			DeliveryStreamName: aws.String(sn),
		}
		resp, err := conn.DescribeDeliveryStream(describeOpts)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == "ResourceNotFoundException" {
					return 42, "DESTROYED", nil
				}
				return nil, awsErr.Code(), err
			}
			return nil, "failed", err
		}

		return resp.DeliveryStreamDescription, *resp.DeliveryStreamDescription.DeliveryStreamStatus, nil
	}
}
