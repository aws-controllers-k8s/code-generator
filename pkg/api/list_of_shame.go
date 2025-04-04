package api

type persistAPIType struct {
	output bool
	input  bool
}

type persistAPITypes map[string]map[string]persistAPIType

func (ts persistAPITypes) Lookup(serviceName, opName string) persistAPIType {
	service, ok := shamelist[serviceName]
	if !ok {
		return persistAPIType{}
	}

	return service[opName]
}

func (ts persistAPITypes) Input(serviceName, opName string) bool {
	return ts.Lookup(serviceName, opName).input
}

func (ts persistAPITypes) Output(serviceName, opName string) bool {
	return ts.Lookup(serviceName, opName).output
}

// shamelist is used to not rename certain operation's input and output shapes.
// We need to maintain backwards compatibility with pre-existing services. Since
// not generating unique input/output shapes is not desired, we will generate
// unique input/output shapes for new operations.
var shamelist = persistAPITypes{
	"AutoScaling": {
		"ResumeProcesses": {
			input: true,
		},
		"SuspendProcesses": {
			input: true,
		},
	},
	"CognitoIdentity": {
		"CreateIdentityPool": {
			output: true,
		},
		"DescribeIdentity": {
			output: true,
		},
		"DescribeIdentityPool": {
			output: true,
		},
		"UpdateIdentityPool": {
			input:  true,
			output: true,
		},
	},
	"DirectConnect": {
		"AllocateConnectionOnInterconnect": {
			output: true,
		},
		"AllocateHostedConnection": {
			output: true,
		},
		"AllocatePrivateVirtualInterface": {
			output: true,
		},
		"AllocatePublicVirtualInterface": {
			output: true,
		},
		"AssociateConnectionWithLag": {
			output: true,
		},
		"AssociateHostedConnection": {
			output: true,
		},
		"AssociateVirtualInterface": {
			output: true,
		},
		"CreateConnection": {
			output: true,
		},
		"CreateInterconnect": {
			output: true,
		},
		"CreateLag": {
			output: true,
		},
		"CreatePrivateVirtualInterface": {
			output: true,
		},
		"CreatePublicVirtualInterface": {
			output: true,
		},
		"DeleteConnection": {
			output: true,
		},
		"DeleteLag": {
			output: true,
		},
		"DescribeConnections": {
			output: true,
		},
		"DescribeConnectionsOnInterconnect": {
			output: true,
		},
		"DescribeHostedConnections": {
			output: true,
		},
		"DescribeLoa": {
			output: true,
		},
		"DisassociateConnectionFromLag": {
			output: true,
		},
		"UpdateLag": {
			output: true,
		},
	},
	"EC2": {
		"AttachVolume": {
			output: true,
		},
		"CreateSnapshot": {
			output: true,
		},
		"CreateVolume": {
			output: true,
		},
		"DetachVolume": {
			output: true,
		},
	},
	"ElasticBeanstalk": {
		"ComposeEnvironments": {
			output: true,
		},
		"CreateApplication": {
			output: true,
		},
		"CreateApplicationVersion": {
			output: true,
		},
		"CreateConfigurationTemplate": {
			output: true,
		},
		"CreateEnvironment": {
			output: true,
		},
		"DescribeEnvironments": {
			output: true,
		},
		"TerminateEnvironment": {
			output: true,
		},
		"UpdateApplication": {
			output: true,
		},
		"UpdateApplicationVersion": {
			output: true,
		},
		"UpdateConfigurationTemplate": {
			output: true,
		},
		"UpdateEnvironment": {
			output: true,
		},
	},
	"ElasticTranscoder": {
		"CreateJob": {
			output: true,
		},
	},
	"Glacier": {
		"DescribeJob": {
			output: true,
		},
		"UploadArchive": {
			output: true,
		},
		"CompleteMultipartUpload": {
			output: true,
		},
	},
	"IAM": {
		"GetContextKeysForCustomPolicy": {
			output: true,
		},
		"GetContextKeysForPrincipalPolicy": {
			output: true,
		},
		"SimulateCustomPolicy": {
			output: true,
		},
		"SimulatePrincipalPolicy": {
			output: true,
		},
	},
	"Kinesis": {
		"DisableEnhancedMonitoring": {
			output: true,
		},
		"EnableEnhancedMonitoring": {
			output: true,
		},
	},
	"Lambda": {
		"UpdateFunctionCode": {
			output: true,
		},
		"UpdateFunctionConfiguration": {
			output: true,
		},
	},
	"MQ": {
		"CreateConfiguration": {
			input:  true,
			output: true,
		},
		"CreateUser": {
			input: true,
		},
		"DescribeUser": {
			output: true,
		},
		"DescribeConfigurationRevision": {
			output: true,
		},
		"ListBrokers": {
			output: true,
		},
		"ListConfigurations": {
			output: true,
		},
		"ListConfigurationRevisions": {
			output: true,
		},
		"ListUsers": {
			output: true,
		},
		"UpdateConfiguration": {
			input:  true,
			output: true,
		},
		"UpdateUser": {
			input: true,
		},
	},
	"RDS": {
		"ModifyDBClusterParameterGroup": {
			output: true,
		},
		"ModifyDBParameterGroup": {
			output: true,
		},
		"ResetDBClusterParameterGroup": {
			output: true,
		},
		"ResetDBParameterGroup": {
			output: true,
		},
	},
	"Redshift": {
		"DescribeLoggingStatus": {
			output: true,
		},
		"DisableLogging": {
			output: true,
		},
		"EnableLogging": {
			output: true,
		},
		"ModifyClusterParameterGroup": {
			output: true,
		},
		"ResetClusterParameterGroup": {
			output: true,
		},
	},
	"S3": {
		"GetBucketNotification": {
			input:  true,
			output: true,
		},
		"GetBucketNotificationConfiguration": {
			input:  true,
			output: true,
		},
	},
	"ServerlessApplicationRepository": {
		"CreateApplication": {
			input: true,
		},
		"CreateApplicationVersion": {
			input: true,
		},
		"CreateCloudFormationChangeSet": {
			input: true,
		},
		"UpdateApplication": {
			input: true,
		},
	},
	"SWF": {
		"CountClosedWorkflowExecutions": {
			output: true,
		},
		"CountOpenWorkflowExecutions": {
			output: true,
		},
		"CountPendingActivityTasks": {
			output: true,
		},
		"CountPendingDecisionTasks": {
			output: true,
		},
		"ListClosedWorkflowExecutions": {
			output: true,
		},
		"ListOpenWorkflowExecutions": {
			output: true,
		},
	},
}
