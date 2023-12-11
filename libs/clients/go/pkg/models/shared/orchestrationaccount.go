// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type OrchestrationAccount struct {
	Address          string                         `json:"address"`
	EffectiveVolumes map[string]OrchestrationVolume `json:"effectiveVolumes,omitempty"`
	Metadata         map[string]string              `json:"metadata"`
	Volumes          map[string]OrchestrationVolume `json:"volumes,omitempty"`
}

func (o *OrchestrationAccount) GetAddress() string {
	if o == nil {
		return ""
	}
	return o.Address
}

func (o *OrchestrationAccount) GetEffectiveVolumes() map[string]OrchestrationVolume {
	if o == nil {
		return nil
	}
	return o.EffectiveVolumes
}

func (o *OrchestrationAccount) GetMetadata() map[string]string {
	if o == nil {
		return map[string]string{}
	}
	return o.Metadata
}

func (o *OrchestrationAccount) GetVolumes() map[string]OrchestrationVolume {
	if o == nil {
		return nil
	}
	return o.Volumes
}
