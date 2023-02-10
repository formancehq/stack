package v1beta1

type Status struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions Conditions `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

type ReplicationStatus struct {
	Status `json:",inline"`
	// +optional
	Replicas int32 `json:"replicas"`
	// +optional
	Selector string `json:"selector"`
}
