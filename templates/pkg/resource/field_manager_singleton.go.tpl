{{- define "field_manager_singleton" -}}
// key returns the primary key for a singleton managed field. There is at
// most one managed field per parent, so the key is a constant; this causes
// computeDelta to route value changes through toUpdate instead of
// toCreate+toDelete (which would wipe the item).
func key(r *resource) string {
	_ = r
	return "_singleton"
}

// sync reconciles the singleton managed field. Since there can only be one
// item per parent, computeDelta yields at most one entry in exactly one of
// {toCreate, toUpdate, toDelete}; no iteration or batching is needed.
func (rm *resourceManager) sync(ctx context.Context, d delta) error {
	var err error
	if len(d.toCreate) > 0 {
		_, err = rm.sdkCreate(ctx, d.toCreate[0])
		return err
	}
	if len(d.toUpdate) > 0 {
		_, err = rm.sdkCreate(ctx, d.toUpdate[0])
		return err
	}
	if len(d.toDelete) > 0 {
		_, err = rm.sdkDelete(ctx, d.toDelete[0])
		return err
	}
	return nil
}

// Get reads the current state of the managed field from AWS and writes the
// result back onto the parent's injected spec field.
func (rm *resourceManager) Get(
	ctx context.Context,
	parent *svcapitypes.{{ ParentKind .CRD }},
) error {
	// Build a minimal seed carrying only the primary key — do NOT use
	// convertFromParent here, because that copies the full desired spec
	// into the seed, which would make latest look identical to desired
	// even when AWS hasn't applied the write yet.
	seed := &resource{ko: &svcapitypes.{{ .CRD.Kind }}{}}
{{ ParentPrimaryKeyAssign .CRD (ParentKind .CRD) "seed.ko" }}
	found, err := rm.sdkFind(ctx, seed)
	if err != nil {
		if err == ackerr.NotFound {
			parent.{{ ManagedFieldPath .CRD }} = nil
			return nil
		}
		return err
	}
	if found == nil {
		parent.{{ ManagedFieldPath .CRD }} = nil
		return nil
	}
	parent.{{ ManagedFieldPath .CRD }} = &found.ko.Spec
	return nil
}

// convertFromParent builds a single managed field instance from the parent.
// The parent's injected field is copied wholesale into the managed field Spec,
// then the parent's primary key is written into the managed field's primary
// key field so the SDK call can identify the resource.
func convertFromParent(parent *svcapitypes.{{ ParentKind .CRD }}) []resource {
	ko := &svcapitypes.{{ .CRD.Kind }}{}
	if parent.{{ ManagedFieldPath .CRD }} != nil {
		ko.Spec = *parent.{{ ManagedFieldPath .CRD }}
	}
{{ ParentPrimaryKeyAssign .CRD (ParentKind .CRD) "ko" }}
	return []resource{ {ko: ko} }
}
{{- end }}
