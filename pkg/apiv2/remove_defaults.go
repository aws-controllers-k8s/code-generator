package apiv2

// BadDefaultsAssignment stores the Service Models that have members with defaultValues (considered non pointers)
// but still have pointers. This list was retrieved from aws-sdk-go-v2, link below
// https://github.com/aws/aws-sdk-go-v2/blob/4ad9d5996fd752f0756be2dbbdd4f8a4841fe362/codegen/smithy-aws-go-codegen/src/main/java/software/amazon/smithy/aws/go/codegen/customization/RemoveDefaults.java#L19-L39
var BadDefaultsAssignment = map[string]map[string]bool{
	"AWSS3ControlServiceV20180820": {
		"BlockPublicAcls":       true,
		"IgnorePublicAcls":      true,
		"BlockPublicPolicy":     true,
		"RestrictPublicBuckets": true,
	},
	"Evidently": {
		"ResultsPeriod": true,
	},
	"AmplifyUIBuilder": {
		"MaxResults":                  true,
		"PlaceIndexSearchResultLimit": true,
	},
	"PaymentCryptographyDataPlane": {
		"IntegerRangeBetween4And12": true,
	},
	"AwsToledoWebService": {
		// WorkerCount targets WorkerCounts, which contains bad default
		//https://github.com/aws/aws-sdk-go-v2/blob/dfcf25b6d15674848f71bd7e5ee8ed22b115c6a2/codegen/sdk-codegen/aws-models/emr-serverless.json#L2111
		"WorkerCount": true,
	},
	"imagebuilder": {
		"setDefaultVersion": true,
	},
	"AmazonBedrockAgentBuildTimeLambda": {
		"StorageDays": true,
	},
}

func hasBadDefualtAssignment(serviceName, shapeName string) bool {
	service, found := BadDefaultsAssignment[serviceName]
	if !found {
		return false
	}

	return service[shapeName]
}
