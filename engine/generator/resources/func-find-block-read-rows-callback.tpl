{{ define "funcFindBlockReadRowsCallback" }}
        {{- /* 定义接收数据的参数*/}}
        {{- range $n, $v := $.Fields -}}
        {{ if $v.VarAlias }}
        {{- if or $v.Blob $v.JSON }}
        {{ if .Ptr }}{{ $.BeanVarName }}.{{.Name}} = &{{.FieldType}}{}{{ end }}
        {{- end }}
        var {{ $v.VarAlias }} {{ $v.VarType }}
        {{- end }}
        {{- end }}
        if err = {{ .SQLReturnVarName }}.Scan({{- range $n, $v := .Fields -}}
            {{ if ne $n 0 }}, {{ end }}&{{ if $v.VarAlias }}{{ $v.VarAlias }}{{ else }}{{ $.BeanVarName }}.{{.Name}}{{ end }}
            {{- end -}}); err != nil {
            return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}} err
        }
        {{/* 把数据赋给返回值 */}}
        {{- range $n, $v := $.Fields -}}
        {{- if $v.VarAlias -}}
        {{ if $v.Blob }}
        if err = {{ $.BeanVarName }}.{{.Name}}.Read({{ $v.VarAlias }}); err != nil {
            return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}} err
        }
        {{ else if $v.JSON }}
        err = {{ $.DBUtilPackage }}.ParseStruct({{$v.VarAlias}}, {{if not .Ptr}}&{{end}}{{ $.BeanVarName }}.{{.Name}})
        if err != nil {
            return {{ if $.Result.Bean }}{{ if $.Result.Bean.Object.Type.IsPtr }}nil{{else}}{{ $.BeanTypeName }}{}{{end}}, {{end}} err
        }
        {{ else if $v.Time }}
        {{ $.BeanVarName }}.{{.Name}} = {{if .Ptr}}&{{end}}{{$v.VarAlias}}.Time()
        {{- end }}
        {{- end }}
        {{- end }}
{{- end }}