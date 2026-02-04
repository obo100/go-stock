# go-stock CLI

## 安装与构建

```bash
# 在仓库根目录

go build -o go-stock-cli ./cmd/go-stock-cli
```

## 全局参数

- `--db` 数据库路径，默认 `data/stock.db`
- `--format` 输出格式：`markdown` 或 `json`
- `--timeout` 超时秒数
- `--quiet` 关闭控制台日志（仅输出结果）

## 初始化

```bash
go-stock-cli init

go-stock-cli init --no-stock-data

go-stock-cli init --update-basic-info
```

## AI 聊天

```bash
go-stock-cli ai chat --question "分析一下立讯精密" --stock "立讯精密" --code "sz002475" --ai-config-id 1

go-stock-cli ai chat --question "分析一下立讯精密" --ai-config-id 1 --enable-tools --stream
```

## AI 新闻总结

```bash
go-stock-cli ai summary --question "总结今天市场" --ai-config-id 1
go-stock-cli ai summary --question "总结今天市场并给出建议" --ai-config-id 1 --enable-tools
```

## Agent 对话

```bash
go-stock-cli agent chat --question "最近半导体行情" --ai-config-id 1
```

## 市场快讯

```bash
go-stock-cli market news --limit 20
```

## 股票行情

```bash
go-stock-cli stock quote --code "sh600000"
```

## 配置管理

```bash
go-stock-cli config export --file config.json

go-stock-cli config import --file config.json

go-stock-cli config add-ai --name "deepseek" --base-url "https://api.deepseek.com" --api-key "KEY" --model "deepseek-chat" --ai-timeout 60

go-stock-cli config list-ai

go-stock-cli recommend list --page 1 --page-size 10
```

## HTTP 服务

```bash
go-stock-cli serve --listen 127.0.0.1:17688
```
