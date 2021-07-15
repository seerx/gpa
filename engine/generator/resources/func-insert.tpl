{{ define "funcInsert" }}
func ({{.Repo.Instance}} *{{.Repo.Name}}) {{.Name}}(
{{- range $n, $v := .Input.Args -}}
{{ if ne $n 0 }}, {{ end }}{{.Name}} {{.Type}}
{{- end -}}
) {{ if gt .Result.Count 1 }}({{ end }}{{ if .Result.Bean }}{{ if .Result.Bean.Ptr }}*{{end}}{{ .BeanTypeName }}, {{ end }}error{{ if gt .Result.Count 1 }}){{ end }} {
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
		{{ if $.Result.Bean -}}
		return {{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, err
		{{- else }}return err
		{{- end }}
	}
	{{ else if $v.Time }}
	var {{ $v.VarAlias }}tp *{{ $.DBUtilPackage }}.TimeProp
	{{ $v.VarAlias }}tp, err = {{ $.DBUtilPackage }}.NewTimeProp("{{ $v.TimeProp.TypeName }}", {{ $v.TimeProp.Nullable }}, "{{ $v.TimeProp.TimeZone }}")
	if err != nil {
		{{ if $.Result.Bean -}}
		return {{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, err
		{{- else }}return err
		{{- end }}
	}
	{{ $v.VarAlias }} := {{ $.DBUtilPackage }}.FormatColumnTime({{$.Repo.Instance}}.p.GetTimeStampzFormat(),
		{{$.Repo.Instance}}.p.GetTimezone(),
		{{ $v.VarAlias }}tp,
		{{ $v.VarName }})
	{{- end }}
	{{- end }}
	{{- end }}
	sql := "{{ .SQL }}"
	{{ if .Input.ContextArgName -}}
	_, err = {{.Repo.Instance}}.p.Executor().ExecContext({{ .Input.ContextArgName }}, sql
		{{- range $n, $v := $.SQLParams -}}
		{{ if $v.VarAlias }}, {{ $v.VarAlias }}{{ else }}, {{ $v.VarName }}{{ end }}
		{{- end -}}
	)
	{{ else -}}
	_, err = {{.Repo.Instance}}.p.Executor().Exec(sql
		{{- range $n, $v := $.SQLParams -}}
		{{ if $v.VarAlias }}, {{ $v.VarAlias }}{{ else }}, {{ $v.VarName }}{{ end }}
		{{- end -}}
	)
	{{ end -}}
	if err != nil {
		{{ if .Result.Bean -}}
		return {{ if .Result.Bean.Ptr }}nil{{else}}{{ .BeanTypeName }}{}{{end}}, err
		{{- else }}return err
		{{- end }}
	}
	{{ if .Result.Bean -}}
	return {{ if not .Result.Bean.Ptr }}*{{end}}{{ .BeanVarName }}, nil
	{{- else }}return nil
	{{- end }}
}
{{ end }}