{{ define "funcInsert" }}
func ({{.Repo.Instance}} *{{.Repo.Name}}) {{.Name}}(
{{- range $n, $v := .Input.Args -}}
{{ if ne $n 0 }}, {{ end }}{{.Name}} {{.Type}}
{{- end -}}
) {{ if gt .Result.Count 1 }}({{ end }}{{ if .Result.Bean }}{{ if .Result.Bean.Object.Type.IsPtr }}*{{end}}{{ .BeanTypeName }}, {{ end }}error{{ if gt .Result.Count 1 }}){{ end }} {
	{{ if .BeanVarNeedCreate -}}
	{{ .BeanVarName }} := &{{ .BeanTypeName }}{}
	{{ end -}}
	{{- range .BeanFieldSetValues -}}
	{{ $.BeanVarName }}.{{.VarName}} = {{.ValueName}}
	{{ end -}}
	var err error
	{{- range $n, $v := $.SQLParams -}}
	{{- if $v.VarAlias }}
	{{ if $v.JSON -}}
	var {{ $v.VarAlias }} string
	{{ $v.VarAlias }}, err = {{ $.DBUtilPackage }}.Struct2String({{ $v.VarName }})
	if err != nil {
		return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{ end }}err
	}
	{{ else if $v.Time }}
	var {{ $v.VarAlias }}tp *{{ $.DBUtilPackage }}.TimeProp
	{{ $v.VarAlias }}tp, err = {{ $.DBUtilPackage }}.NewTimeProp("{{ $v.TimeProp.TypeName }}", {{ $v.TimeProp.Nullable }}, "{{ $v.TimeProp.TimeZone }}")
	if err != nil {
		return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{ end }}err
	}
	{{ $v.VarAlias }} := {{ $.DBUtilPackage }}.FormatColumnTime({{$.Repo.Instance}}.p.GetTimeStampzFormat(),
		{{$.Repo.Instance}}.p.GetTimezone(),
		{{ $v.VarAlias }}tp,
		{{ $v.VarName }})
	{{- end }}
	{{- end }}
	{{- end }}
	{{- if .AutoinrPrimaryKeyField }}
	var {{ .SQLReturnVarName }} {{ .SQLPackage }}.Result
	{{- end }}
	sql := "{{ .SQL }}"
	{{ if .Input.ContextArgName -}}
	{{- if .AutoinrPrimaryKeyField }}{{ .SQLReturnVarName }}{{ else }}_{{ end }}, err = {{.Repo.Instance}}.p.Executor().ExecContext({{ .Input.ContextArgName }}, sql
		{{- range $n, $v := $.SQLParams -}}
		{{ if $v.VarAlias }}, {{ $v.VarAlias }}{{ else }}, {{ $v.VarName }}{{ end }}
		{{- end -}}
	)
	{{ else -}}
	{{- if .AutoinrPrimaryKeyField }}{{ .SQLReturnVarName }}{{ else }}_{{ end }}, err = {{.Repo.Instance}}.p.Executor().Exec(sql
		{{- range $n, $v := $.SQLParams -}}
		{{ if $v.VarAlias }}, {{ $v.VarAlias }}{{ else }}, {{ $v.VarName }}{{ end }}
		{{- end -}}
	)
	{{ end -}}
	if err != nil {
		return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{ end }}err
	}
	{{ if .Result.Bean -}}
	{{ if .AutoinrPrimaryKeyField }}
	{{ .SQLReturnVarName }}InsertID, err := {{ .SQLReturnVarName }}.LastInsertId()
	if err != nil {
		return {{ if .Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ .BeanTypeName }}{}{{end}}, err
	}
	{{ if eq .AutoinrPrimaryKeyFieldType "int64" }}
	{{- .BeanVarName }}.{{ .AutoinrPrimaryKeyVarName }} = {{ .SQLReturnVarName }}InsertID
	{{ else }}
	{{- .BeanVarName }}.{{ .AutoinrPrimaryKeyVarName }} = {{.AutoinrPrimaryKeyFieldType}}({{ .SQLReturnVarName }}InsertID)
	{{ end }}
	{{- end }}
	{{- end }}
	return {{ if $.Result.Bean }}{{ if not $.Result.Bean.Object.Type.IsPtr }}*{{end}}{{ .BeanVarName }}, {{ end }}err 
}
{{ end }}