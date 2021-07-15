{{ define "funcUpdate" }}
func ({{.Repo.Instance}} *{{.Repo.Name}}) {{.Name}}(
{{- range $n, $v := .Input.Args -}}
{{ if ne $n 0 }}, {{ end }}{{.Name}} {{.Type}}
{{- end -}}
) {{ if gt .Result.Count 1 }}({{ end }}{{ if .Result.AffectVar }}int64, {{ end }}{{ if .Result.Bean }}{{ if .Result.Bean.Ptr }}*{{end}}{{ .BeanTypeName }}, {{ end }}error{{ if gt .Result.Count 1 }}){{ end }} {
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
	sqlParams := []interface{}{ {{- range $n, $v := $.SQLParams -}}{{if ne $n 0}}, {{ end }}{{ if $v.VarAlias }}{{ $v.VarAlias }}{{ else }}{{ $v.VarName }}{{ end }}{{- end }}}
	// where 参数
	{{- range $n, $v := $.SQLWhereParams -}}
	{{ if $v.IsInOperator }}
	if len({{ $v.VarName }}) <= 0 {
		return {{ if $.Result.AffectVar }}0, {{ end }}{{ if $.Result.Bean }}{{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}{{ $.DBUtilPackage }}.NewErrParamIsEmpty("{{ $v.VarName }}")
	}
	sql = {{ $.DBUtilPackage }}.TakePlaceHolder(sql, "{{$v.InParamPlaceHolder}}", len({{ $v.VarName }}))
	for _, varP := range {{ $v.VarName }} {
		sqlParams = append(sqlParams, varP)
	}
	{{- else }}
	sqlParams = append(sqlParams, {{ if $v.VarAlias }}{{ $v.VarAlias }}{{ else }}{{ $v.VarName }}{{ end }})
	{{- end }}
	{{- end }}
	{{ if .Result.AffectVar }}var {{ .SQLReturnVarName }} {{ .SQLPackage }}.Result{{ end }}
	{{ if .Input.ContextArgName -}}
	{{ if .Result.AffectVar }}{{ .SQLReturnVarName }}{{ else }}_{{ end }}, err = {{.Repo.Instance}}.p.Executor().ExecContext({{ .Input.ContextArgName }}, sql, sqlParams...)
	{{ else -}}
	{{ if .Result.AffectVar }}{{ .SQLReturnVarName }}{{ else }}_{{ end }}, err = {{.Repo.Instance}}.p.Executor().Exec(sql, sqlParams...)
	{{ end -}}
	if err != nil {
		return {{ if $.Result.AffectVar }}0, {{ end }}{{ if $.Result.Bean }}{{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}err
	}
	{{ if .Result.AffectVar -}}
	{{ .Result.AffectVar }}, err := {{ .SQLReturnVarName }}.RowsAffected()
	if err != nil {
		return {{ .Result.AffectVar }}, {{ if $.Result.Bean }}{{ if $.Result.Bean.Ptr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}err
	}
	{{- end }}
	return {{ if $.Result.AffectVar }}{{ .Result.AffectVar }}, {{ end }}{{ if $.Result.Bean }}{{ if $.Result.Bean.Ptr }}{{ .BeanVarName }}{{else}}*{{ .BeanVarName }}{{end}}, {{end}}nil
}
{{ end }}