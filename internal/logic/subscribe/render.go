package subscribe

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/maxliu9403/ProxyHub/models"
)

type ClashTemplateData struct {
	ISPProtocol string
	ISPServer   string
	ISPPort     int
	ISPUsername string
	ISPPassword string
}

func (s *Svc) renderClashConfig(proxy *models.Proxy) (string, error) {
	// 准备数据
	data := ClashTemplateData{
		ISPProtocol: proxy.ProxyType, // 示例字段
		ISPServer:   proxy.IP,
		ISPPort:     int(proxy.Port),
		ISPUsername: proxy.Username,
		ISPPassword: proxy.Password,
	}

	// 读取模板
	tmplContent, err := os.ReadFile("/Users/liuxiang/go/src/github.com/ProxyHub/configs/base_proxy.yaml")
	if err != nil {
		return "", fmt.Errorf("读取模板失败: %w", err)
	}

	// 创建模板
	tmpl, err := template.New("clash").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("模板解析失败: %w", err)
	}

	// 渲染模板
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", fmt.Errorf("模板渲染失败: %w", err)
	}

	return rendered.String(), nil
}
