package mailer

import (
	"bytes"
	"html/template"

	"github.com/maxliu9403/ProxyHub/models"
)

const mailTpl = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; font-size: 14px; color: #333; }
    h2 { color: #2c3e50; }
    table { width: 100%; border-collapse: collapse; margin: 10px 0; }
    th, td { border: 1px solid #ccc; padding: 8px 12px; text-align: left; }
    th { background-color: #f4f4f4; }
    hr { border: none; border-top: 1px solid #ccc; margin: 20px 0; }
    .section { margin-bottom: 30px; }
  </style>
</head>
<body>
  <h2>ClashProxyHub 定时清理报告</h2>
  <p>本次共处理分组数：<strong>{{len .}}</strong></p>

  {{range .}}
  <div class="section">
    <hr>
    <h3>分组：{{.GroupName}}</h3>
    <p><strong>最大在线数：</strong> {{.MaxOnline}}</p>

    <h4>解绑模拟器列表</h4>
    <table>
      <tr>
        <th>BrowserID</th>
        <th>UUID</th>
      </tr>
      {{range .UnbindEmulator}}
      <tr>
        <td>{{.BrowserID}}</td>
        <td>{{.UUID}}</td>
      </tr>
      {{end}}
    </table>

    <h4>已释放代理 IP 列表</h4>
    <table>
      <tr>
        <th>IP 地址</th>
        <th>释放次数</th>
      </tr>
      {{range .ReleaseIPDetail}}
      <tr>
        <td>{{.IP}}</td>
        <td>{{.Count}}</td>
      </tr>
      {{end}}
    </table>
  </div>
  {{end}}
</body>
</html>
`

func RenderReleaseHTML(data []*models.GroupReleaseResult) (string, error) {
	var buf bytes.Buffer
	tpl := template.Must(template.New("release").Parse(mailTpl))
	err := tpl.Execute(&buf, data)
	return buf.String(), err
}
