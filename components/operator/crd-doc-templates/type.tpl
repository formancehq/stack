{{- define "type" -}}
{{- $type := .Type -}}
{{- $recurse := .Recurse -}}
{{- $prefix := .Prefix -}}

{{ repeat ($prefix | int) "#" }} {{ $type.Name }}

{{ if $type.IsAlias }}_Underlying type:_ _{{ markdownRenderTypeLink $type.UnderlyingType  }}_{{ end }}

{{ $type.Doc }}

{{/*{{ if $type.Validation -}}*/}}
{{/*_Validation:_*/}}
{{/*{{- range $type.Validation }}*/}}
{{/*- {{ . }}*/}}
{{/*{{- end }}*/}}
{{/*{{- end }}*/}}

{{/*{{ if $type.References -}}*/}}
{{/*_Appears in:_*/}}
{{/*{{- range $type.SortedReferences }}*/}}
{{/*- {{ markdownRenderTypeLink . }}*/}}
{{/*{{- end }}*/}}
{{/*{{- end }}*/}}

{{ if $type.Members -}}
| Field | Description | Default | Validation |
| --- | --- | --- | --- |
{{ if $type.GVK -}}
| `apiVersion` _string_ | `{{ $type.GVK.Group }}/{{ $type.GVK.Version }}` | | |
| `kind` _string_ | `{{ $type.GVK.Kind }}` | | |
{{ end -}}

{{ range $type.Members -}}
{{- if or (has "Type: string" .Type.Validation) (eq (markdownRenderType .Type) "[Duration](#duration)") -}}
| `{{ .Name  }}` _string_ | {{ template "type_members" . }} | {{ markdownRenderDefault .Default }} | {{ range .Validation -}} {{ . }} <br />{{ end }} |
{{- else -}}
| `{{ .Name  }}` _{{ markdownRenderType .Type }}_ | {{ template "type_members" . }} | {{ markdownRenderDefault .Default }} | {{ range .Validation -}} {{ . }} <br />{{ end }} |
{{- end }}
{{ end -}}

{{ end -}}

{{- if $recurse }}
{{- if eq $type.Kind 4}}
{{- $dummy := set . "Fields" $type.UnderlyingType.Fields }}
{{- else }}
{{- $dummy := set . "Fields" $type.Fields }}
{{- end }}
{{- range $k, $field := .Fields }}
{{- if hasPrefix "github.com/formancehq/operator/api" $field.Type.Package }}
{{- if has "Type: string" $field.Type.Validation }}{{ continue }}{{ end }}
{{ template "type" (dict "Type" $field.Type "Recurse" true "Prefix" (min (add $prefix 1) 6)) }}
{{- end }}
{{- end }}
{{- end }}

{{- end -}}
