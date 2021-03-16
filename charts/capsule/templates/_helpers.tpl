{{/*
Expand the name of the chart.
*/}}
{{- define "capsule.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "capsule.fullname" -}}
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
{{- define "capsule.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "capsule.labels" -}}
helm.sh/chart: {{ include "capsule.chart" . }}
{{ include "capsule.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "capsule.selectorLabels" -}}
app.kubernetes.io/name: {{ include "capsule.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "capsule.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "capsule.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the manager fully-qualified Docker image to use
*/}}
{{- define "capsule.managerFullyQualifiedDockerImage" -}}
{{- printf "%s:%s" .Values.manager.image.repository ( .Values.manager.image.tag | default (printf "v%s" .Chart.AppVersion) ) -}}
{{- end }}

{{/*
Create the kube-rbac-proxy fully-qualified Docker image to use
*/}}
{{- define "capsule.sidecarFullyQualifiedDockerImage" -}}
{{- printf "%s:%s" .Values.sidecar.image.repository .Values.sidecar.image.tag -}}
{{- end }}

{{/*
Create the jobs fully-qualified Docker image to use
*/}}
{{- define "capsule.jobsFullyQualifiedDockerImage" -}}
{{- printf "%s:%s" .Values.jobs.image.repository .Values.jobs.image.tag -}}
{{- end }}

{{/*
Create the Capsule Deployment name to use
*/}}
{{- define "capsule.deploymentName" -}}
{{- printf "%s-controller-manager" (include "capsule.fullname" .) -}}
{{- end }}

{{/*
Create the Capsule CA Secret name to use
*/}}
{{- define "capsule.secretCaName" -}}
{{- printf "%s-ca" (include "capsule.fullname" .) -}}
{{- end }}

{{/*
Create the Capsule TLS Secret name to use
*/}}
{{- define "capsule.secretTlsName" -}}
{{- printf "%s-tls" (include "capsule.fullname" .) -}}
{{- end }}
