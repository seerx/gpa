{{ define "repo" }}
type {{.Name}} struct {
	p *{{ .DBPackage }}.Provider
}
{{ range .Funcs -}}
{{ if eq .Template "insert" -}}
	{{- template "funcInsert" . -}}
{{- else if eq .Template "update" -}}
	{{- template "funcUpdate" . }}
{{- else if eq .Template "delete" -}}
	{{- template "funcDelete" . -}}
{{- else if eq .Template "find" -}}
	{{- template "funcFind" . -}}
{{- else if eq .Template "count" -}}
	{{- template "funcCount" . -}}
{{- end }}
{{- end -}}
{{- end -}}