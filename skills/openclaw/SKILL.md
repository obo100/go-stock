# OpenClaw Skill: go-stock CLI

## 说明

该技能通过调用本地 `go-stock-cli` 执行业务逻辑，适合 OpenClaw 以命令行方式集成。

## 前置条件

- 已构建 `go-stock-cli`
- 已完成 `go-stock-cli init` 或已导入配置

## 示例命令

### AI 分析

```bash
go-stock-cli ai chat --question "分析一下立讯精密" --stock "立讯精密" --code "sz002475" --ai-config-id 1 --format json
```

### AI 市场总结（启用工具生成推荐记录）

```bash
go-stock-cli ai summary --question "总结今天市场并给出建议" --ai-config-id 1 --enable-tools --format json
```

### Agent 对话

```bash
go-stock-cli agent chat --question "最近半导体行情" --ai-config-id 1 --format json
```

### 查看推荐股票列表

```bash
go-stock-cli recommend list --page 1 --page-size 10 --format json
```

## 建议

- 在 OpenClaw 中将输出格式固定为 `json`，便于结构化解析
- 可结合 `--stream` 获取流式输出
- 若需要生成“推荐股票列表”，请使用 `ai summary` 且开启 `--enable-tools`
