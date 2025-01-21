package api

var shapeNameAliases = map[string]map[string]string{
	"APIGateway": {
		"RequestValidator": "UpdateRequestValidatorOutput",
		"GatewayResponse":  "UpdateGatewayResponseOutput",
	},
	"Lambda": {
		"Concurrency": "PutFunctionConcurrencyOutput",
	},
	"Neptune": {
		"DBClusterParameterGroupNameMessage": "ResetDBClusterParameterGroupOutput",
		"DBParameterGroupNameMessage":        "ResetDBParameterGroupOutput",
	},
	"RDS": {
		"DBClusterBacktrack": "BacktrackDBClusterOutput",
	},
}
