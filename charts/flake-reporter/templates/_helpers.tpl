{{- define "flakeReporter.labels" -}}
app.kubernetes.io/instance: "{{.Release.Name}}"
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/name: flake-reporter
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
helm.sh/chart: "{{.Chart.Name}}-{{.Chart.Version}}"
{{- end }}
