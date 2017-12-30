package aws

import "github.com/jmcgill/formation/core"

func Importers() map[string]core.Importer {
    return map[string]core.Importer {
        // App Autoscaling Resources
        // "aws_appautoscaling_policy": &AwsAppautoscalingPolicyImporter{},
        // "aws_appautoscaling_scheduled_action": &AwsAppautoscalingScheduledActionImporter{},
        // "aws_appautoscaling_target": &AwsAppautoscalingTargetImporter{},

        // Athena Resources
        // "aws_athena_database": &AwsAthenaDatabaseImporter{},
        // "aws_athena_named_query": &AwsAthenaNamedQueryImporter{},

        // Batch Resources
        // "aws_batch_compute_environment": &AwsBatchComputeEnvironmentImporter{},
        // "aws_batch_job_definition": &AwsBatchJobDefinitionImporter{},
        // "aws_batch_job_queue": &AwsBatchJobQueueImporter{},

        // CloudFormation Resources
        // "aws_cloudformation_stack": &AwsCloudformationStackImporter{},

        // CloudFront Resources
        // "aws_cloudfront_distribution": &AwsCloudfrontDistributionImporter{},
        // "aws_cloudfront_origin_access_identity": &AwsCloudfrontOriginAccessIdentityImporter{},

        // CloudTrail Resources
        // "aws_cloudtrail": &AwsCloudtrailImporter{},

        // CloudWatch Resources
        // "aws_cloudwatch_dashboard": &AwsCloudwatchDashboardImporter{},
        // "aws_cloudwatch_event_rule": &AwsCloudwatchEventRuleImporter{},
        // "aws_cloudwatch_event_target": &AwsCloudwatchEventTargetImporter{},
        // "aws_cloudwatch_log_destination": &AwsCloudwatchLogDestinationImporter{},
        // "aws_cloudwatch_log_destination_policy": &AwsCloudwatchLogDestinationPolicyImporter{},
        // "aws_cloudwatch_log_group": &AwsCloudwatchLogGroupImporter{},
        // "aws_cloudwatch_log_metric_filter": &AwsCloudwatchLogMetricFilterImporter{},
        // "aws_cloudwatch_log_stream": &AwsCloudwatchLogStreamImporter{},
        // "aws_cloudwatch_log_subscription_filter": &AwsCloudwatchLogSubscriptionFilterImporter{},
        // "aws_cloudwatch_metric_alarm": &AwsCloudwatchMetricAlarmImporter{},

        // Config Resources
        // "aws_config_config_rule": &AwsConfigConfigRuleImporter{},
        // "aws_config_configuration_recorder": &AwsConfigConfigurationRecorderImporter{},
        // "aws_config_configuration_recorder_status": &AwsConfigConfigurationRecorderStatusImporter{},
        // "aws_config_delivery_channel": &AwsConfigDeliveryChannelImporter{},

        // Database Migration Service
        // "aws_dms_certificate": &AwsDmsCertificateImporter{},
        // "aws_dms_endpoint": &AwsDmsEndpointImporter{},
        // "aws_dms_replication_instance": &AwsDmsReplicationInstanceImporter{},
        // "aws_dms_replication_subnet_group": &AwsDmsReplicationSubnetGroupImporter{},
        // "aws_dms_replication_task": &AwsDmsReplicationTaskImporter{},

        // Device Farm Resources
        // "aws_devicefarm_project": &AwsDevicefarmProjectImporter{},

        // Directory Service Resources
        // "aws_directory_service_directory": &AwsDirectoryServiceDirectoryImporter{},

        // Direct Connect Resources
        // "aws_dx_connection": &AwsDxConnectionImporter{},
        // "aws_dx_connection_association": &AwsDxConnectionAssociationImporter{},
        // "aws_dx_lag": &AwsDxLagImporter{},

        // DynamoDB Resources
        // "aws_dynamodb_table": &AwsDynamodbTableImporter{},

        // EC2 Resources
        // "aws_ami": &AwsAmiImporter{},
        // "aws_ami_copy": &AwsAmiCopyImporter{},
        // "aws_ami_from_instance": &AwsAmiFromInstanceImporter{},
        // "aws_ami_launch_permission": &AwsAmiLaunchPermissionImporter{},
        // "aws_app_cookie_stickiness_policy": &AwsAppCookieStickinessPolicyImporter{},
        // "aws_autoscaling_attachment": &AwsAutoscalingAttachmentImporter{},
        // "aws_autoscaling_group": &AwsAutoscalingGroupImporter{},
        // "aws_autoscaling_lifecycle_hook": &AwsAutoscalingLifecycleHookImporter{},
        // "aws_autoscaling_notification": &AwsAutoscalingNotificationImporter{},
        // "aws_autoscaling_policy": &AwsAutoscalingPolicyImporter{},
        // "aws_autoscaling_schedule": &AwsAutoscalingScheduleImporter{},
        // "aws_snapshot_create_volume_permission": &AwsSnapshotCreateVolumePermissionImporter{},
        // "aws_ebs_snapshot": &AwsEbsSnapshotImporter{},
        "aws_ebs_volume": &AwsEbsVolumeImporter{},
        // "aws_eip": &AwsEipImporter{},
        // "aws_eip_association": &AwsEipAssociationImporter{},
        // "aws_elb": &AwsElbImporter{},
        // "aws_elb_attachment": &AwsElbAttachmentImporter{},
        "aws_instance": &AwsInstanceImporter{},
        // "aws_key_pair": &AwsKeyPairImporter{},
        // "aws_launch_configuration": &AwsLaunchConfigurationImporter{},
        // "aws_lb_cookie_stickiness_policy": &AwsLbCookieStickinessPolicyImporter{},
        // "aws_lb_ssl_negotiation_policy": &AwsLbSslNegotiationPolicyImporter{},
        // "aws_load_balancer_backend_server_policy": &AwsLoadBalancerBackendServerPolicyImporter{},
        // "aws_load_balancer_listener_policy": &AwsLoadBalancerListenerPolicyImporter{},
        // "aws_load_balancer_policy": &AwsLoadBalancerPolicyImporter{},
        // "aws_placement_group": &AwsPlacementGroupImporter{},
        // "aws_proxy_protocol_policy": &AwsProxyProtocolPolicyImporter{},
        // "aws_spot_datafeed_subscription": &AwsSpotDatafeedSubscriptionImporter{},
        // "aws_spot_fleet_request": &AwsSpotFleetRequestImporter{},
        // "aws_spot_instance_request": &AwsSpotInstanceRequestImporter{},
        "aws_volume_attachment": &AwsVolumeAttachmentImporter{},

        // Load Balancing Resources
        // "aws_lb": &AwsLbImporter{},
        // "aws_lb_listener": &AwsLbListenerImporter{},
        // "aws_lb_listener_rule": &AwsLbListenerRuleImporter{},
        // "aws_lb_target_group": &AwsLbTargetGroupImporter{},
        // "aws_lb_target_group_attachment": &AwsLbTargetGroupAttachmentImporter{},

        // ECS Resources
        // "aws_ecr_lifecycle_policy": &AwsEcrLifecyclePolicyImporter{},
        // "aws_ecr_repository": &AwsEcrRepositoryImporter{},
        // "aws_ecr_repository_policy": &AwsEcrRepositoryPolicyImporter{},
        // "aws_ecs_cluster": &AwsEcsClusterImporter{},
        "aws_ecs_service": &AwsEcsServiceImporter{},
        // "aws_ecs_task_definition": &AwsEcsTaskDefinitionImporter{},

        // EFS Resources
        // "aws_efs_file_system": &AwsEfsFileSystemImporter{},
        // "aws_efs_mount_target": &AwsEfsMountTargetImporter{},

        // ElastiCache Resources
        // "aws_elasticache_cluster": &AwsElasticacheClusterImporter{},
        // "aws_elasticache_parameter_group": &AwsElasticacheParameterGroupImporter{},
        // "aws_elasticache_replication_group": &AwsElasticacheReplicationGroupImporter{},
        // "aws_elasticache_security_group": &AwsElasticacheSecurityGroupImporter{},
        // "aws_elasticache_subnet_group": &AwsElasticacheSubnetGroupImporter{},

        // IAM Resources
        // "aws_iam_access_key": Secret - cannot be imported,
        "aws_iam_account_alias": &AwsIamAccountAliasImporter{},
        "aws_iam_account_password_policy": &AwsIamAccountPasswordPolicyImporter{},
        "aws_iam_group": &AwsIamGroupImporter{},
        "aws_iam_group_membership": &AwsIamGroupMembershipImporter{},
        // "aws_iam_group_policy": &AwsIamGroupPolicyImporter{},
        // "aws_iam_group_policy_attachment": &AwsIamGroupPolicyAttachmentImporter{},
        // "aws_iam_instance_profile": &AwsIamInstanceProfileImporter{},
        // "aws_iam_openid_connect_provider": &AwsIamOpenidConnectProviderImporter{},
        // "aws_iam_policy": &AwsIamPolicyImporter{},
        // "aws_iam_policy_attachment": &AwsIamPolicyAttachmentImporter{},
        "aws_iam_role": &AwsIamRoleImporter{},
        // "aws_iam_role_policy": &AwsIamRolePolicyImporter{},
        // "aws_iam_role_policy_attachment": &AwsIamRolePolicyAttachmentImporter{},
        // "aws_iam_saml_provider": &AwsIamSamlProviderImporter{},
        // "aws_iam_server_certificate": &AwsIamServerCertificateImporter{},
        "aws_iam_user": &AwsIamUserImporter{},
        // "aws_iam_user_login_profile": &AwsIamUserLoginProfileImporter{},
        // "aws_iam_user_policy": &AwsIamUserPolicyImporter{},
        // "aws_iam_user_policy_attachment": &AwsIamUserPolicyAttachmentImporter{},
        // "aws_iam_user_ssh_key": &AwsIamUserSshKeyImporter{},

        // Kinesis Resources
        // "aws_kinesis_stream": &AwsKinesisStreamImporter{},

        // Kinesis Firehose Resources
        // "aws_kinesis_firehose_delivery_stream": &AwsKinesisFirehoseDeliveryStreamImporter{},

        // KMS Resources
        // "aws_kms_alias": &AwsKmsAliasImporter{},
        // "aws_kms_key": &AwsKmsKeyImporter{},

        // Lambda Resources
        // "aws_lambda_alias": &AwsLambdaAliasImporter{},
        // "aws_lambda_event_source_mapping": &AwsLambdaEventSourceMappingImporter{},
        // "aws_lambda_function": &AwsLambdaFunctionImporter{},
        // "aws_lambda_permission": &AwsLambdaPermissionImporter{},

        // RDS Resources
        // "aws_db_event_subscription": &AwsDbEventSubscriptionImporter{},
        // "aws_db_instance": &AwsDbInstanceImporter{},
        // "aws_db_option_group": &AwsDbOptionGroupImporter{},
        // "aws_db_parameter_group": &AwsDbParameterGroupImporter{},
        // "aws_db_security_group": &AwsDbSecurityGroupImporter{},
        // "aws_db_snapshot": &AwsDbSnapshotImporter{},
        // "aws_db_subnet_group": &AwsDbSubnetGroupImporter{},
        // "aws_rds_cluster": &AwsRdsClusterImporter{},
        // "aws_rds_cluster_instance": &AwsRdsClusterInstanceImporter{},
        // "aws_rds_cluster_parameter_group": &AwsRdsClusterParameterGroupImporter{},

        // Redshift Resources
        // "aws_redshift_cluster": &AwsRedshiftClusterImporter{},
        // "aws_redshift_parameter_group": &AwsRedshiftParameterGroupImporter{},
        // "aws_redshift_security_group": &AwsRedshiftSecurityGroupImporter{},
        // "aws_redshift_subnet_group": &AwsRedshiftSubnetGroupImporter{},

        // Route53 Resources
        // "aws_route53_delegation_set": &AwsRoute53DelegationSetImporter{},
        // "aws_route53_health_check": &AwsRoute53HealthCheckImporter{},
        // "aws_route53_record": &AwsRoute53RecordImporter{},
        "aws_route53_zone": &AwsRoute53ZoneImporter{},
        // "aws_route53_zone_association": &AwsRoute53ZoneAssociationImporter{},

        // S3 Resources
        "aws_s3_bucket": &AwsS3BucketImporter{},
        // "aws_s3_bucket_notification": &AwsS3BucketNotificationImporter{},
        // "aws_s3_bucket_object": &AwsS3BucketObjectImporter{},
        // "aws_s3_bucket_policy": &AwsS3BucketPolicyImporter{},

        // SES Resources
        // "aws_ses_active_receipt_rule_set": &AwsSesActiveReceiptRuleSetImporter{},
        // "aws_ses_domain_identity": &AwsSesDomainIdentityImporter{},
        // "aws_ses_domain_dkim": &AwsSesDomainDkimImporter{},
        // "aws_ses_receipt_filter": &AwsSesReceiptFilterImporter{},
        // "aws_ses_receipt_rule": &AwsSesReceiptRuleImporter{},
        // "aws_ses_receipt_rule_set": &AwsSesReceiptRuleSetImporter{},
        // "aws_ses_configuration_set": &AwsSesConfigurationSetImporter{},
        // "aws_ses_event_destination": &AwsSesEventDestinationImporter{},
        // "aws_ses_template": &AwsSesTemplateImporter{},

        // SNS Resources
        // "aws_sns_topic": &AwsSnsTopicImporter{},
        // "aws_sns_topic_policy": &AwsSnsTopicPolicyImporter{},
        // "aws_sns_topic_subscription": &AwsSnsTopicSubscriptionImporter{},

        // SQS Resources
        "aws_sqs_queue": &AwsSqsQueueImporter{},
        // "aws_sqs_queue_policy": &AwsSqsQueuePolicyImporter{},

        // VPC Resources
        // "aws_customer_gateway": &AwsCustomerGatewayImporter{},
        // "aws_default_network_acl": &AwsDefaultNetworkAclImporter{},
        // "aws_default_route_table": &AwsDefaultRouteTableImporter{},
        // "aws_default_security_group": &AwsDefaultSecurityGroupImporter{},
        // "aws_default_subnet": &AwsDefaultSubnetImporter{},
        // "aws_default_vpc": &AwsDefaultVpcImporter{},
        // "aws_default_vpc_dhcp_options": &AwsDefaultVpcDhcpOptionsImporter{},
        // "aws_egress_only_internet_gateway": &AwsEgressOnlyInternetGatewayImporter{},
        // "aws_flow_log": &AwsFlowLogImporter{},
        // "aws_internet_gateway": &AwsInternetGatewayImporter{},
        // "aws_main_route_table_association": &AwsMainRouteTableAssociationImporter{},
        // "aws_nat_gateway": &AwsNatGatewayImporter{},
        // "aws_network_acl": &AwsNetworkAclImporter{},
        // "aws_network_acl_rule": &AwsNetworkAclRuleImporter{},
        // "aws_network_interface": &AwsNetworkInterfaceImporter{},
        // "aws_network_interface_attachment": &AwsNetworkInterfaceAttachmentImporter{},
        // "aws_route": &AwsRouteImporter{},
        // "aws_route_table": &AwsRouteTableImporter{},
        // "aws_route_table_association": &AwsRouteTableAssociationImporter{},
        // "aws_security_group": &AwsSecurityGroupImporter{},
        // "aws_network_interface_sg_attachment": &AwsNetworkInterfaceSgAttachmentImporter{},
        // "aws_security_group_rule": &AwsSecurityGroupRuleImporter{},
        // "aws_subnet": &AwsSubnetImporter{},
        // "aws_vpc": &AwsVpcImporter{},
        // "aws_vpc_dhcp_options": &AwsVpcDhcpOptionsImporter{},
        // "aws_vpc_dhcp_options_association": &AwsVpcDhcpOptionsAssociationImporter{},
        // "aws_vpc_endpoint": &AwsVpcEndpointImporter{},
        // "aws_vpc_endpoint_route_table_association": &AwsVpcEndpointRouteTableAssociationImporter{},
        // "aws_vpc_peering_connection": &AwsVpcPeeringConnectionImporter{},
        // "aws_vpc_peering_connection_accepter": &AwsVpcPeeringConnectionAccepterImporter{},
        // "aws_vpn_connection": &AwsVpnConnectionImporter{},
        // "aws_vpn_connection_route": &AwsVpnConnectionRouteImporter{},
        // "aws_vpn_gateway": &AwsVpnGatewayImporter{},
        // "aws_vpn_gateway_attachment": &AwsVpnGatewayAttachmentImporter{},
        // "aws_vpn_gateway_route_propagation": &AwsVpnGatewayRoutePropagationImporter{},

        // ***** TIER 2 RESOURCES ****

        // CodeBuild Resources
        // "aws_codebuild_project": &AwsCodebuildProjectImporter{},

        // CodeCommit Resources
        // "aws_codecommit_repository": &AwsCodecommitRepositoryImporter{},
        // "aws_codecommit_trigger": &AwsCodecommitTriggerImporter{},

        // CodeDeploy Resources
        // "aws_codedeploy_app": &AwsCodedeployAppImporter{},
        // "aws_codedeploy_deployment_config": &AwsCodedeployDeploymentConfigImporter{},
        // "aws_codedeploy_deployment_group": &AwsCodedeployDeploymentGroupImporter{},

        // CodePipeline Resources
        // "aws_codepipeline": &AwsCodepipelineImporter{},

        // Cognito Resources
        // "aws_cognito_identity_pool": &AwsCognitoIdentityPoolImporter{},
        // "aws_cognito_identity_pool_roles_attachment": &AwsCognitoIdentityPoolRolesAttachmentImporter{},
        // "aws_cognito_user_pool": &AwsCognitoUserPoolImporter{},

        // WAF Resources
        // "aws_waf_byte_match_set": &AwsWafByteMatchSetImporter{},
        // "aws_waf_ipset": &AwsWafIpsetImporter{},
        // "aws_waf_rule": &AwsWafRuleImporter{},
        // "aws_waf_rate_based_rule": &AwsWafRateBasedRuleImporter{},
        // "aws_waf_size_constraint_set": &AwsWafSizeConstraintSetImporter{},
        // "aws_waf_sql_injection_match_set": &AwsWafSqlInjectionMatchSetImporter{},
        // "aws_waf_web_acl": &AwsWafWebAclImporter{},
        // "aws_waf_xss_match_set": &AwsWafXssMatchSetImporter{},

        // WAF Regional Resources
        // "aws_wafregional_byte_match_set": &AwsWafregionalByteMatchSetImporter{},
        // "aws_wafregional_ipset": &AwsWafregionalIpsetImporter{},


        // SSM Resources
        // "aws_ssm_activation": &AwsSsmActivationImporter{},
        // "aws_ssm_association": &AwsSsmAssociationImporter{},
        // "aws_ssm_document": &AwsSsmDocumentImporter{},
        // "aws_ssm_maintenance_window": &AwsSsmMaintenanceWindowImporter{},
        // "aws_ssm_maintenance_window_target": &AwsSsmMaintenanceWindowTargetImporter{},
        // "aws_ssm_maintenance_window_task": &AwsSsmMaintenanceWindowTaskImporter{},
        // "aws_ssm_patch_baseline": &AwsSsmPatchBaselineImporter{},
        // "aws_ssm_patch_group": &AwsSsmPatchGroupImporter{},
        // "aws_ssm_parameter": &AwsSsmParameterImporter{},

        // API Gateway Resources
        // "aws_api_gateway_account": &AwsApiGatewayAccountImporter{},
        // "aws_api_gateway_api_key": &AwsApiGatewayApiKeyImporter{},
        // "aws_api_gateway_authorizer": &AwsApiGatewayAuthorizerImporter{},
        // "aws_api_gateway_base_path_mapping": &AwsApiGatewayBasePathMappingImporter{},
        // "aws_api_gateway_client_certificate": &AwsApiGatewayClientCertificateImporter{},
        // "aws_api_gateway_deployment": &AwsApiGatewayDeploymentImporter{},
        // "aws_api_gateway_domain_name": &AwsApiGatewayDomainNameImporter{},
        // "aws_api_gateway_gateway_response": &AwsApiGatewayGatewayResponseImporter{},
        // "aws_api_gateway_integration": &AwsApiGatewayIntegrationImporter{},
        // "aws_api_gateway_integration_response": &AwsApiGatewayIntegrationResponseImporter{},
        // "aws_api_gateway_method": &AwsApiGatewayMethodImporter{},
        // "aws_api_gateway_method_response": &AwsApiGatewayMethodResponseImporter{},
        // "aws_api_gateway_method_settings": &AwsApiGatewayMethodSettingsImporter{},
        // "aws_api_gateway_model": &AwsApiGatewayModelImporter{},
        // "aws_api_gateway_resource": &AwsApiGatewayResourceImporter{},
        // "aws_api_gateway_rest_api": &AwsApiGatewayRestApiImporter{},
        // "aws_api_gateway_stage": &AwsApiGatewayStageImporter{},
        // "aws_api_gateway_usage_plan": &AwsApiGatewayUsagePlanImporter{},
        // "aws_api_gateway_usage_plan_key": &AwsApiGatewayUsagePlanKeyImporter{},

        // Lightsail Resources
        // "aws_lightsail_domain": &AwsLightsailDomainImporter{},
        // "aws_lightsail_instance": &AwsLightsailInstanceImporter{},
        // "aws_lightsail_key_pair": &AwsLightsailKeyPairImporter{},
        // "aws_lightsail_static_ip": &AwsLightsailStaticIpImporter{},
        // "aws_lightsail_static_ip_attachment": &AwsLightsailStaticIpAttachmentImporter{},

        // MQ Resources
        // "aws_mq_broker": &AwsMqBrokerImporter{},
        // "aws_mq_configuration": &AwsMqConfigurationImporter{},

        // MediaStore Resources
        // "aws_media_store_container": &AwsMediaStoreContainerImporter{},

        // OpsWorks Resources
        // "aws_opsworks_application": &AwsOpsworksApplicationImporter{},
        // "aws_opsworks_custom_layer": &AwsOpsworksCustomLayerImporter{},
        // "aws_opsworks_ganglia_layer": &AwsOpsworksGangliaLayerImporter{},
        // "aws_opsworks_haproxy_layer": &AwsOpsworksHaproxyLayerImporter{},
        // "aws_opsworks_instance": &AwsOpsworksInstanceImporter{},
        // "aws_opsworks_java_app_layer": &AwsOpsworksJavaAppLayerImporter{},
        // "aws_opsworks_memcached_layer": &AwsOpsworksMemcachedLayerImporter{},
        // "aws_opsworks_mysql_layer": &AwsOpsworksMysqlLayerImporter{},
        // "aws_opsworks_nodejs_app_layer": &AwsOpsworksNodejsAppLayerImporter{},
        // "aws_opsworks_permission": &AwsOpsworksPermissionImporter{},
        // "aws_opsworks_php_app_layer": &AwsOpsworksPhpAppLayerImporter{},
        // "aws_opsworks_rails_app_layer": &AwsOpsworksRailsAppLayerImporter{},
        // "aws_opsworks_rds_db_instance": &AwsOpsworksRdsDbInstanceImporter{},
        // "aws_opsworks_stack": &AwsOpsworksStackImporter{},
        // "aws_opsworks_static_web_layer": &AwsOpsworksStaticWebLayerImporter{},
        // "aws_opsworks_user_profile": &AwsOpsworksUserProfileImporter{},

        // Service Catalog Resources
        // "aws_servicecatalog_portfolio": &AwsServicecatalogPortfolioImporter{},

        // Service Discovery Resources
        // "aws_service_discovery_private_dns_namespace": &AwsServiceDiscoveryPrivateDnsNamespaceImporter{},
        // "aws_service_discovery_public_dns_namespace": &AwsServiceDiscoveryPublicDnsNamespaceImporter{},

        // Step Function Resources
        // "aws_sfn_activity": &AwsSfnActivityImporter{},
        // "aws_sfn_state_machine": &AwsSfnStateMachineImporter{},

        // SimpleDB Resources
        // "aws_simpledb_domain": &AwsSimpledbDomainImporter{},

        // Elastic Beanstalk Resources
        // "aws_elastic_beanstalk_application": &AwsElasticBeanstalkApplicationImporter{},
        // "aws_elastic_beanstalk_application_version": &AwsElasticBeanstalkApplicationVersionImporter{},
        // "aws_elastic_beanstalk_configuration_template": &AwsElasticBeanstalkConfigurationTemplateImporter{},
        // "aws_elastic_beanstalk_environment": &AwsElasticBeanstalkEnvironmentImporter{},

        // Elastic Map Reduce Resources
        // "aws_emr_cluster": &AwsEmrClusterImporter{},
        // "aws_emr_instance_group": &AwsEmrInstanceGroupImporter{},
        // "aws_emr_security_configuration": &AwsEmrSecurityConfigurationImporter{},

        // ElasticSearch Resources
        // "aws_elasticsearch_domain": &AwsElasticsearchDomainImporter{},
        // "aws_elasticsearch_domain_policy": &AwsElasticsearchDomainPolicyImporter{},

        // Elastic Transcoder Resources

        // Glacier Resources
        // "aws_glacier_vault": &AwsGlacierVaultImporter{},

        // IoT Resources
        // "aws_iot_certificate": &AwsIotCertificateImporter{},
        // "aws_iot_policy": &AwsIotPolicyImporter{},

        // Inspector Resources
        // "aws_inspector_assessment_target": &AwsInspectorAssessmentTargetImporter{},
        // "aws_inspector_assessment_template": &AwsInspectorAssessmentTemplateImporter{},
    }
}