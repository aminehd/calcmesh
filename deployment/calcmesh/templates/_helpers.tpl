{{- define "calcmesh.labels" -}}
app.kubernetes.io/name: {{ .name }}
app.kubernetes.io/part-of: calcmesh
{{- end -}}
