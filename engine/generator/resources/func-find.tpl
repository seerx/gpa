{{ define "funcFind" }}
func ({{.Repo.Instance}} *{{.Repo.Name}}) {{.Name}}(
{{- range $n, $v := .Input.Args -}}
{{ if ne $n 0 }}, {{ end }}{{.Name}} {{.Type}}
{{- end -}}
) {{ if gt .Result.Count 1 }}({{ end }}{{ if .Result.Bean }}{{ .Result.ReturnTypeName }}, {{ end }}error{{ if gt .Result.Count 1 }}){{ end }}{
    var err error
	{{- range $n, $v := $.SQLWhereParams -}}
	{{- if $v.VarAlias }}
	{{ if $v.JSON -}}
	var {{ $v.VarAlias }} string
	{{ $v.VarAlias }}, err = {{ $.DBUtilPackage }}.Struct2String({{ $v.VarName }})
	if err != nil {
		return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}err
	}
	{{ else if $v.Time }}
	var {{ $v.VarAlias }}tp *{{ $.DBUtilPackage }}.TimeProp
	{{ $v.VarAlias }}tp, err = {{ $.DBUtilPackage }}.NewTimeProp("{{ $v.TimeProp.TypeName }}", {{ $v.TimeProp.Nullable }}, "{{ $v.TimeProp.TimeZone }}")
	if err != nil {
		return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}err
	}
	{{ $v.VarAlias }} := {{ $.DBUtilPackage }}.FormatColumnTime({{$.Repo.Instance}}.p.GetTimeStampzFormat(),
		{{ $.Repo.Instance }}.p.GetTimezone(),
		{{ $v.VarAlias }}tp,
		{{ $v.VarName }})
	{{- end }}
	{{- end }}
	{{- end }}
	sql := "{{ .SQL }}"
	var sqlParams []interface{}
	// where 参数
	{{- range $n, $v := $.SQLWhereParams -}}
	{{ if $v.IsInOperator }}
	if len({{ $v.VarName }}) <= 0 {
		return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}}{{ $.DBUtilPackage }}.NewErrParamIsEmpty("{{ $v.VarName }}")
	}
	sql = {{ $.DBUtilPackage }}.TakePlaceHolder(sql, "{{$v.InParamPlaceHolder}}", len({{ $v.VarName }}))
	for _, varP := range {{ $v.VarName }} {
		sqlParams = append(sqlParams, varP)
	}
	{{- else }}
	sqlParams = append(sqlParams, {{ if $v.VarAlias }}{{ $v.VarAlias }}{{ else }}{{ $v.VarName }}{{ end }})
	{{- end }}
	{{- end }}
    {{ if .Result.FindOne }}
    {{- /* 只返回一条记录 */}}
    {{ if .Input.ContextArgName -}}
	{{ .SQLReturnVarName }} := {{.Repo.Instance}}.p.Executor().QueryContextRow({{ .Input.ContextArgName }}, sql, sqlParams...)
	{{ else -}}
	{{ .SQLReturnVarName }} := {{.Repo.Instance}}.p.Executor().QueryRow(sql, sqlParams...)
	{{- end -}}
	{{- template "funcFindBlockReadRow" . }}
	return {{ .BeanVarName }}, nil
    {{- else }} {{- /* 只返回一条记录 -- 结束 */}}
    {{- /*返回多条记录*/}}
    var {{ .SQLReturnVarName }} *{{ .SQLPackage }}.Rows
	{{ if .Input.ContextArgName -}}
	{{ .SQLReturnVarName }}, err = {{.Repo.Instance}}.p.Executor().QueryContextRows({{ .Input.ContextArgName }}, sql, sqlParams...)
	{{ else -}}
	{{ .SQLReturnVarName }}, err = {{.Repo.Instance}}.p.Executor().QueryRows(sql, sqlParams...)
	{{- end }}
	if err != nil {
		return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}} err
	}
	{{ if .Input.Callback -}}
	{{- /*使用回调函数返回数据*/}}
	{{ .BeanVarName }} := &{{ .BeanTypeName }}{}
	for {{ .SQLReturnVarName }}.Next() {
	{{- template "funcFindBlockReadRowsCallback" . }}
		if err = {{.Input.Callback.Name}}({{if not .Input.CallbackArgIsPtr}}*{{end}}{{.BeanVarName}}); err != nil {
			return err
		}
	}
	return nil
	{{- else -}} {{- /*使用回调函数返回数据 -- 结束*/}}
	{{- /*使用数据或 map 返回多条数据*/}}
	{{ .SQLReturnVarName }}Results := {{.Result.ReturnTypeName}}{}
	for {{ .SQLReturnVarName }}.Next() {
	{{- template "funcFindBlockReadRows" . }}
		{{- if .Result.Bean.IsMap }}
		{{ .SQLReturnVarName }}Results[{{.Input.KeyGenerator.Name}}({{if not .Input.KeyGeneratorArgIsPtr}}*{{end}}{{.BeanVarName}})] = {{ if not .Result.Bean.Object.Type.IsPtr }}*{{end}}{{.BeanVarName}}
		{{- else if .Result.Bean.IsSlice }}
		{{ .SQLReturnVarName }}Results = append({{ .SQLReturnVarName }}Results, {{ if not .Result.Bean.Object.Type.IsPtr }}*{{end}}{{.BeanVarName}})
		{{- end }}
	}
	return {{ .SQLReturnVarName }}Results, nil
	{{- end -}} {{- /*使用数据或 map 返回多条数据 -- 结束*/}}
    {{- end }} {{- /*返回多条记录 -- 结束*/}}
}
{{ end }}