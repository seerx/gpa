//+mro-ignore
// DO NOT EDIT THIS FILE
// Generated by mro at {{.Time}}
package {{.package}}

import (
	{{- range $k, $v := .imports }}
	{{ $k }} "{{ $v }}"
	{{- end }}
)
{{ range .Repos -}}
{{ template "repo" . }}
{{- end }}