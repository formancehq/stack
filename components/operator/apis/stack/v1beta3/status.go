package v1beta3

import (
	"github.com/formancehq/operator/internal/collectionutils"
)

type Status struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions Conditions `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

func (in *Status) GetConditions() []Condition {
	return in.Conditions
}

func (in *Status) GetCondition(conditionType string) *Condition {
	if in != nil {
		for _, condition := range in.Conditions {
			if condition.Type == conditionType {
				return &condition
			}
		}
	}
	return nil
}

func (in *Status) SetCondition(condition Condition) {
	for i, c := range in.Conditions {
		if c.Type == condition.Type {
			in.Conditions[i] = condition
			return
		}
	}
	in.Conditions = append(in.Conditions, condition)
}

func (in *Status) RemoveCondition(v string) {
	in.Conditions = collectionutils.Filter(in.Conditions, func(stack Condition) bool {
		return stack.Type != v
	})
}
