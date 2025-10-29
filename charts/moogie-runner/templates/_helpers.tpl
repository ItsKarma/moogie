{{/*
Expand the name of the chart.
*/}}
{{- define "moogie-runner.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "moogie-runner.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "moogie-runner.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "moogie-runner.labels" -}}
helm.sh/chart: {{ include "moogie-runner.chart" . }}
{{ include "moogie-runner.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "moogie-runner.selectorLabels" -}}
app.kubernetes.io/name: {{ include "moogie-runner.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "moogie-runner.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "moogie-runner.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get image repository for a check
*/}}
{{- define "moogie-runner.imageRepository" -}}
{{- if .check.image }}
{{- .check.image.repository | default .root.Values.global.image.repository }}
{{- else }}
{{- .root.Values.global.image.repository }}
{{- end }}
{{- end }}

{{/*
Get image tag for a check
*/}}
{{- define "moogie-runner.imageTag" -}}
{{- if .check.image }}
{{- .check.image.tag | default .root.Values.global.image.tag }}
{{- else }}
{{- .root.Values.global.image.tag }}
{{- end }}
{{- end }}

{{/*
Get image pull policy for a check
*/}}
{{- define "moogie-runner.imagePullPolicy" -}}
{{- if .check.image }}
{{- .check.image.pullPolicy | default .root.Values.global.image.pullPolicy }}
{{- else }}
{{- .root.Values.global.image.pullPolicy }}
{{- end }}
{{- end }}
