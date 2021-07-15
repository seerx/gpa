{{ define "funcCount" }}
func ({{.Repo.Instance}} *{{.Repo.Name}}) {{.Name}}(
{{- range $n, $v := .Input.Args -}}
{{ if ne $n 0 }}, {{ end }}{{.Name}} {{.Type}}
{{- end -}}
) (int64, {{ if .Result.Bean }}{{ if .Result.Bean.Ptr }}*{{end}}{{ .BeanTypeName }}, {{ end }}error) {
	{{ if .BeanVarNeedCreate -}}
	{{ .BeanVarName }} := &{{ .BeanTypeName }}{}
	{{ end -}}
	{{- range .BeanFieldSetValues -}}
	{{ $.BeanVarName }}.{{.VarName}} = {{.ValueName}}
	{{ end -}}
	var err error
	{{- range $n, $v := $.SQLWhereParams -}}
	{{- if $v.VarAlias }}
	{{ if $v.JSON -}}
	var {{ $v.VarAlias }} string
	{{ $v.VarAlias }}, err = {{ $.DBUtilPackage }}.Struct2String({{ $v.VarName }})
	if err != nil {
		return {{ if $.Result.AffectVar }}0, {{ end }}{{ if $.Result.Bean }}{{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}err
	}
	{{ else if $v.Time }}
	var {{ $v.VarAlias }}tp *{{ $.DBUtilPackage }}.TimeProp
	{{ $v.VarAlias }}tp, err = {{ $.DBUtilPackage }}.NewTimeProp("{{ $v.TimeProp.TypeName }}", {{ $v.TimeProp.Nullable }}, "{{ $v.TimeProp.TimeZone }}")
	if err != nil {
		return {{ if $.Result.AffectVar }}0, {{ end }}{{ if $.Result.Bean }}{{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}err
	}
	{{ $v.VarAlias }} := {{ $.DBUtilPackage }}.FormatColumnTime({{$.Repo.Instance}}.p.GetTimeStampzFormat(),
		{{ $.Repo.Instance }}.p.GetTimezone(),
		{{ $v.VarAlias }}tp,
		{{ $v.VarName }})
	{{- end }}
	{{- end }}
	{{- end }}
	sql := "{{ .SQL }}"
	{{ if .Input.ContextArgName -}}
	{{ .SQLReturnVarName }} := {{.Repo.Instance}}.p.Executor().QueryContextRow({{ .Input.ContextArgName }}, sql
		{{- range $n, $v := $.SQLWhereParams -}}
		{{ if $v.VarAlias }}, {{ $v.VarAlias }}{{ else }}, {{ $v.VarName }}{{ end }}
		{{- end -}}
	)
	{{ else -}}
	{{ .SQLReturnVarName }} := {{.Repo.Instance}}.p.Executor().QueryRow(sql
		{{- range $n, $v := $.SQLWhereParams -}}
		{{ if $v.VarAlias }}, {{ $v.VarAlias }}{{ else }}, {{ $v.VarName }}{{ end }}
		{{- end -}}
	)
	{{ end -}}
	var {{.Result.CountVar}} int64
	err = {{ .SQLReturnVarName }}.Scan(&{{.Result.CountVar}})
	if err != nil {
		return 0, {{ if $.Result.Bean }}{{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}err
	}
	return {{.Result.CountVar}}, {{ if $.Result.Bean }}{{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}nil
}
{{ end }}