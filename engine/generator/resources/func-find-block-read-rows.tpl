{{ define "funcFindBlockReadRows" }}
        {{/* 定义接收数据的参数*/}}
        {{- .BeanVarName }} := &{{ .BeanTypeName }}{}
        {{- range $n, $v := $.Fields -}}
        {{ if $v.VarAlias }}
        var {{ $v.VarAlias }} {{ $v.VarType }}
        {{- end }}
        {{- end }}
        err = {{ .SQLReturnVarName }}.Scan({{- range $n, $v := .Fields -}}
            {{ if ne $n 0 }}, {{ end }}&{{ if $v.VarAlias }}{{ $v.VarAlias }}{{ else }}{{ $.BeanVarName }}.{{.Name}}{{ end }}
            {{- end -}})
        if err != nil {
            return {{ if $.Result.Bean }}nil, {{end}} err
        }
        {{/* 把数据赋给返回值 */}}
        {{- range $n, $v := $.Fields -}}
        {{- if $v.VarAlias -}}
        {{ if $v.JSON }}
        err = {{ $.DBUtilPackage }}.ParseStruct({{$v.VarAlias}}, {{if not .Ptr}}&{{end}}{{ $.BeanVarName }}.{{.Name}})
        if err != nil {
            return {{ if $.Result.Bean }}nil, {{end}} err
        }
        {{ else if $v.Time }}
        {{ $.BeanVarName }}.{{.Name}} = {{if .Ptr}}&{{end}}{{$v.VarAlias}}.Time()
        {{- end }}
        {{- end }}
        {{- end }}
{{- end }}