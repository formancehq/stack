{{/*
  .Values: Values to search within
  .Key: Key to find in the <service>.config.<key> then in <global>.<key>
  .Default: default an object where each key is a string
*/}}
{{- define "resolveGlobalOrServiceValue" -}}
  {{- $values := .Values -}}
  {{- $key := .Key -}}
  {{- $default := .Default -}}

  {{- $keys := splitList "." $key -}}
  

  {{- $configkeys := splitList "." (print "config." $key) -}}
  {{- $subchartValue := $values -}}
  {{- $found := true -}}
  {{- range $configkeys -}}
    {{- if hasKey $subchartValue . -}}
      {{- $subchartValue = index $subchartValue . -}}
    {{- else -}}
      {{- $found = false -}}
      {{- break -}}
    {{- end -}}
  {{- end -}}

  {{- if not $found -}}
    {{- $subchartValue = $values.global -}}
    {{- $found = true -}}
    {{- range $keys -}}
      {{- if hasKey $subchartValue . -}}
        {{- $subchartValue = index $subchartValue . -}}
      {{- else -}}
        {{- $subchartValue = $default -}}
        {{- $found = false -}}
        {{- break -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}

  {{- if not $found -}}
    {{- $subchartValue = $default -}}
  {{- end -}}

  {{- $subchartValue -}}
{{- end -}}
