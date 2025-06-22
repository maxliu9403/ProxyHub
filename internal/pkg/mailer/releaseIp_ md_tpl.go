package mailer

import (
	"bytes"
	"text/template"

	"github.com/maxliu9403/ProxyHub/models"
	"github.com/yuin/goldmark"
)

const releaseMarkdownTpl = `
# ClashProxyHub 清理报告

共处理分组：{{ len . }}

{{ range . }}
## 分组：{{ .GroupName }}

- 最大在线数：{{ .MaxOnline }}

### 已释放代理 IP

| IP 地址     | 释放次数 |
|------------|----------|
{{- range .ReleaseIPDetail }}
| {{ .IP }} | {{ .Count }} |
{{- end }}

### 解绑模拟器列表

| 浏览器ID | UUID |
|----------|------|
{{- range .UnbindEmulator }}
| {{ .BrowserID }} | {{ .UUID }} |
{{- end }}

---
{{ end }}
`

func MarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RenderReleaseMarkdown(data []*models.GroupReleaseResult) (string, error) {
	tpl, err := template.New("release_md").Parse(releaseMarkdownTpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
