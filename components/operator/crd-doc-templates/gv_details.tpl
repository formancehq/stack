{{- define "rootType" }}
{{ template "type" (dict "Type" . "Recurse" false "Prefix" 4)}}

{{- range $k, $field := .Fields }}
{{ if $field.Type.IsBasic }}{{ continue }}{{ end }}
{{- if eq $field.Name "spec" }}
{{ template "type" (dict "Type" $field.Type "Recurse" true "Prefix" 5) }}
{{- end }}
{{- end }}

{{- range $k, $field := .Fields }}
{{ if $field.Type.IsBasic }}{{ continue }}{{ end }}
{{- if eq $field.Name "status" }}
{{ template "type" (dict "Type" $field.Type "Recurse" true "Prefix" 5) }}
{{- end }}
{{- end }}
{{- end }}

{{- define "gvDetails" -}}
{{- $gv := . -}}

## {{ $gv.GroupVersionString }}

{{ $gv.Doc }}

{{- if $gv.Kinds  }}
Modules :
{{- range $gv.SortedTypes }}
{{- $k8sMetadata := index .Markers "kubebuilder:metadata" }}
{{- if gt (len $k8sMetadata) 0 }}
{{- $labels := (index $k8sMetadata 0).Labels }}
{{- if gt (len $labels) 0 }}
{{- $isModule := has "formance.com/kind=module" $labels }}
{{- if $isModule }}
- {{ $gv.TypeForKind .Name | markdownRenderTypeLink }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

Other resources :
{{- range $gv.SortedTypes }}
{{- if eq .Name "Stack"}}{{ continue }}{{ end }}
{{- if eq .Name "Settings"}}{{ continue }}{{ end }}
{{- $k8sMetadata := index .Markers "kubebuilder:metadata" }}
{{- $isModule := and (gt (len $k8sMetadata) 0) (has "formance.com/kind=module" (index $k8sMetadata 0).Labels) }}
{{- $kubeBuilderIsRoot := index .Markers "kubebuilder:object:root" }}
{{- $isRoot := and (gt (len $kubeBuilderIsRoot) 0) (index $kubeBuilderIsRoot 0) }}
{{- if and (not $isModule) $isRoot }}
- {{ $gv.TypeForKind .Name | markdownRenderTypeLink }}
{{- end }}
{{- end }}

### Main resources

{{- /** Display details about Stack resource */}}
{{ template "rootType" (index $gv.Types "Stack") }}

{{- /** Display details about Setting resource */}}
{{ template "rootType" (index $gv.Types "Settings") }}

### Modules

{{- range $gv.SortedTypes }}
{{- $k8sMetadata := index .Markers "kubebuilder:metadata" }}
{{- if gt (len $k8sMetadata) 0 }}
{{- $labels := (index $k8sMetadata 0).Labels }}
{{- if gt (len $labels) 0 }}
{{- $isModule := has "formance.com/kind=module" $labels }}
{{- if $isModule }}
{{ template "rootType" . }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

### Other resources

{{- range $gv.SortedTypes }}
{{- if eq .Name "Stack" }}{{ continue }}{{ end }}
{{- if eq .Name "Settings" }}{{ continue }}{{ end }}
{{- $k8sMetadata := index .Markers "kubebuilder:metadata" }}
{{- $isModule := and (gt (len $k8sMetadata) 0) (has "formance.com/kind=module" (index $k8sMetadata 0).Labels) }}
{{- $kubeBuilderIsRoot := index .Markers "kubebuilder:object:root" }}
{{- $isRoot := and (gt (len $kubeBuilderIsRoot) 0) (index $kubeBuilderIsRoot 0) }}

{{- if and (not $isModule) $isRoot }}
{{ template "rootType" . }}
{{- end }}
{{- end }}


{{- end }}
{{- end }}

