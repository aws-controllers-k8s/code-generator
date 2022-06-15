{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"context"
	"errors"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/{{ .ServicePackageName }}"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/{{.ServicePackageName }}-controller/apis/{{ .APIVersion }}"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.{{ .APIInterfaceTypeName}}{}
	_ = &svcapitypes.{{ .CRD.Names.Camel }}{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
)

// sdkFind returns SDK-specific information about a supplied resource
{{ if .CRD.Ops.ReadOne }}
	{{- template "sdk_find_read_one" . }}
{{- else if .CRD.Ops.GetAttributes }}
	{{- template "sdk_find_get_attributes" . }}
{{- else if .CRD.Ops.ReadMany }}
	{{- template "sdk_find_read_many" . }}
{{- else }}
	{{- template "sdk_find_not_implemented" . }}
{{- end }}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()

{{- if $hookCode := Hook .CRD "sdk_create_pre_build_request" }}
{{ $hookCode }}
{{- end }}
{{- if $customMethod := .CRD.GetCustomImplementation .CRD.Ops.Create }}
	created, err = rm.{{ $customMethod }}(ctx, desired)
	if created != nil || err != nil {
		return created, err
	}
{{- end }}
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}
{{- if $hookCode := Hook .CRD "sdk_create_post_build_request" }}
{{ $hookCode }} 
{{- end }}

	var resp {{ .CRD.GetOutputShapeGoType .CRD.Ops.Create }}; _ = resp;
	resp, err = rm.sdkapi.{{ .CRD.Ops.Create.ExportedName }}WithContext(ctx, input)
{{- if $hookCode := Hook .CRD "sdk_create_post_request" }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("CREATE", "{{ .CRD.Ops.Create.ExportedName }}", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()
{{- if $hookCode := Hook .CRD "sdk_create_pre_set_output" }}
{{ $hookCode }}
{{- end }}
{{ GoCodeSetCreateOutput .CRD "resp" "ko" 1 }}
	rm.setStatusDefaults(ko)
{{- if $setOutputCustomMethodName := .CRD.SetOutputCustomMethodName .CRD.Ops.Create }}
	// custom set output from response
	ko, err = rm.{{ $setOutputCustomMethodName }}(ctx, desired, resp, ko)
	if err != nil {
		return nil, err
	}
{{- end }}
{{- if $hookCode := Hook .CRD "sdk_create_post_set_output" }}
{{ $hookCode }}
{{- end }}
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.{{ .CRD.Ops.Create.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ .CRD.Ops.Create.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetCreateInput .CRD "r.ko" "res" 1 }}
	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
{{ if .CRD.CustomUpdateMethodName }}
	{{- template "sdk_update_custom" . }}
{{- else if .CRD.Ops.Update }}
	{{- template "sdk_update" . }}
{{- else if .CRD.Ops.SetAttributes }}
	{{- template "sdk_update_set_attributes" . }}
{{- else }}
	{{- template "sdk_update_not_implemented" . }}
{{- end }}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()

{{- if .CRD.Ops.Delete }}
{{- if $hookCode := Hook .CRD "sdk_delete_pre_build_request" }}
{{ $hookCode }}
{{- end }}
{{- if $customMethod := .CRD.GetCustomImplementation .CRD.Ops.Delete }}
	if err = rm.{{ $customMethod }}(ctx, r); err != nil {
		return nil, err
	}
{{- end }}
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
{{- if $hookCode := Hook .CRD "sdk_delete_post_build_request" }}
{{ $hookCode }}
{{- end }}
	var resp {{ .CRD.GetOutputShapeGoType .CRD.Ops.Delete }}; _ = resp;
	resp, err = rm.sdkapi.{{ .CRD.Ops.Delete.Name }}WithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "{{ .CRD.Ops.Delete.Name }}", err)
{{- if $hookCode := Hook .CRD "sdk_delete_post_request" }}
{{ $hookCode }}
{{- end }}
	return nil, err
{{- else }}
	// TODO(jaypipes): Figure this out...
	return nil, nil
{{ end }}
}

{{ if .CRD.Ops.Delete -}}
// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.{{ .CRD.Ops.Delete.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ .CRD.Ops.Delete.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetDeleteInput .CRD "r.ko" "res" 1 }}
	return res, nil
}
{{- end }}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults (
	ko *svcapitypes.{{ .CRD.Names.Camel }},
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions (
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err ==  ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type:   ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type:   ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}

{{- if $reconcileRequeuOnSuccessSeconds := .CRD.ReconcileRequeuOnSuccessSeconds }}
	if syncCondition == nil && onSuccess {
		syncCondition = &ackv1alpha1.Condition{
			Type:   ackv1alpha1.ConditionTypeResourceSynced,
			Status: corev1.ConditionTrue,
		}
		ko.Status.Conditions = append(ko.Status.Conditions, syncCondition)
	}
{{- else }}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
{{- end }}

{{- if $updateConditionsCustomMethodName := .CRD.UpdateConditionsCustomMethodName }}
	// custom update conditions
	customUpdate := rm.{{ $updateConditionsCustomMethodName }}(ko, r, err)
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil || customUpdate {
		return &resource{ko}, true // updated
	}
{{- else }}
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
{{- end }}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
{{- if .CRD.TerminalExceptionCodes }}
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case {{ range $x, $terminalCode := .CRD.TerminalExceptionCodes -}}{{ if ne ($x) (0) }},
		{{ end }} "{{ $terminalCode }}"{{ end }}:
		return true
	default:
		return false
	}
{{- else }}
	// No terminal_errors specified for this resource in generator config
	return false
{{- end }}
}

{{- if .CRD.HasImmutableFieldChanges }}
// getImmutableFieldChanges returns list of immutable fields from the
func (rm *resourceManager) getImmutableFieldChanges(
	delta *ackcompare.Delta,
) []string {
	var fields []string;

	{{- range $immutableField := .CRD.GetImmutableFieldPaths }}
		if delta.DifferentAt("{{$immutableField}}") {
			fields = append(fields,"{{$immutableField}}")
		}
	{{- end }}

	return fields
}

// handleImmutableFieldsChangedCondition validates the immutable fields and set appropriate condition
func (rm *resourceManager) handleImmutableFieldsChangedCondition(
	r *resource,
	delta *ackcompare.Delta,
) *resource {

	fields := rm.getImmutableFieldChanges(delta)
	ko := r.ko.DeepCopy()
	var advisoryCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeAdvisory {
			advisoryCondition = condition
			break
		}
	}

	// Remove the advisory condition if issue is no longer present
	if len(fields) == 0 && advisoryCondition != nil{
		var newConditions []*ackv1alpha1.Condition
		for _, condition := range ko.Status.Conditions {
			if condition.Type != ackv1alpha1.ConditionTypeAdvisory {
				newConditions = append(newConditions,condition)
			}
		}
		ko.Status.Conditions = newConditions
	}

	if len(fields) > 0 {
		if advisoryCondition == nil {
			advisoryCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeAdvisory,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, advisoryCondition)
		}

		advisoryCondition.Status = corev1.ConditionTrue
		message := "Immutable Spec fields have been modified : " + strings.Join(fields, ",")
		advisoryCondition.Message = &message
	}

	return &resource{ko}
}
{{- end }}
{{- if $hookCode := Hook .CRD "sdk_file_end" }}
{{ $hookCode }}
{{- end }}
