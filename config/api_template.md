# {{.Name}}文档

> 文档描述： {{.Description}}

## 请求

### 请求地址

```text
POST /api/{{.Path}}
```

### 请求方式

```text
Content-Type: application/json;charset=utf8
```

### 请求参数

| 参数名称 | 参数位置 | 参数类型 | 是否必须 | 默认值 | 参数说明 |
| -------- | -------- | -------- | -------- | -------- | -------- |
| page | query | int | 否 | 1 | 分页页数 |
| size | query | int | 否 | 10 | 分页每页显示的条数  |
{{- range .RequestParams }}
| {{.Name}} | {{ pl .ParamLocation}} | {{ pt .ParamType }} | {{ if .Required }} 是 {{ else }} 否 {{ end }} | {{.DefaultValue}} | {{.Description}} |
{{- end }}

## 响应

### 响应码

| 响应码 | 描述 |
| -------- | -------- |
| 200 | 服务器正常响应 |
| 400 | 请求参数不正确 |
| 401 | 未认证 |
| 403 | 权限不足 |
| 404 | 请求地址不存在 |
| 405 | 请求方式不正确 |
| 500 | 服务器错误 |

### 响应体

响应格式均为 JSON 格式：

```json
{
  "success": true,
  "message": "",
  "data": []
}
```

其中：

| 参数名称 | 参数类型 | 参数说明 |
| -------- | -------- | -------- |
| success | boolean | 正常为true，错误时为false |
| message | string | 错误描述 |
| data | array(object) | 数据 |

其中 data 数据参数如下：

| 参数名称 | 参数类型 | 参数说明 |
| -------- | -------- | -------- |
{{- range .ResponseParams }}
| {{ if eq .ConvertType 1 }} {{.ConvertValue}} {{ else }} {{.Name}} {{ end }} | {{ pt .ParamType }} | {{.Description}} |
{{- end}}
