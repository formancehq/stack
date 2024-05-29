{{/*
Selector labels
*/}}
{{- define "agent.selectorLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "agent.monitoring.logs" -}}
{{- if eq .Values.global.monitoring.logs.format "json" }}
- name: JSON_FORMATTING_LOGGER
  value: "true"
{{- end }}
{{- end }}