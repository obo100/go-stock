# go-stock HTTP API

## 启动服务

```bash
go-stock-cli serve --listen 127.0.0.1:17688
```

## 认证

- 如设置 `--token`，需在请求头携带 `Authorization: Bearer <token>` 或 `X-Token: <token>`

## 接口列表

### GET /health

返回服务健康状态。

### POST /ai/chat

请求体示例：

```json
{
  "question": "分析一下立讯精密",
  "stock": "立讯精密",
  "stock_code": "sz002475",
  "ai_config_id": 1,
  "sys_prompt_id": 0,
  "enable_tools": true,
  "thinking": false,
  "stream": false,
  "format": "json"
}
```

### POST /ai/summary

请求体示例：

```json
{
  "question": "总结今天市场",
  "ai_config_id": 1,
  "enable_tools": true,
  "stream": true
}
```

### POST /agent/chat

请求体示例：

```json
{
  "question": "最近半导体行情",
  "ai_config_id": 1,
  "stream": true
}
```

### GET /market/news

查询参数：

- `source` 新闻来源，可选
- `limit` 数量，默认 20

## 流式响应

当 `stream=true` 时，使用 SSE：

- `event: delta`
- `event: done`
