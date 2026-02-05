# OpenClaw Skill: go-stock CLI

## 说明

该技能通过调用本地 `go-stock-cli` 执行业务逻辑，适合 OpenClaw 以命令行方式集成。

## 前置条件

- 已构建 `go-stock-cli`
- 已完成 `go-stock-cli init` 或已导入配置

## 示例命令

### AI 分析

```bash
go-stock-cli ai chat --question "分析一下立讯精密" --stock "立讯精密" --code "sz002475" --ai-config-id 1 --quiet
```

### AI 市场总结（启用工具生成推荐记录）

```bash
go-stock-cli ai summary --question "总结今天市场并给出建议" --ai-config-id 1 --enable-tools --quiet
```

### Agent 对话

```bash
go-stock-cli agent chat --question "最近半导体行情" --ai-config-id 1 --quiet
```

### 查看推荐股票列表

```bash
go-stock-cli recommend list --page 1 --page-size 10 --quiet
```

### 设置东财用户标识

```bash
go-stock-cli config set-qgqp --qgqp-b-id "YOUR_QGQP_B_ID" --quiet
```

## 持仓管理

### 新增持仓

```bash
go-stock-cli position add --code "000592.SZ" --price 11.83 --volume 1000 --take-profit 13.50 --stop-loss 10.50 --quiet
```

### 查看持仓列表

```bash
go-stock-cli position list --quiet
```

### 分析持仓全部股票操作建议

```bash
go-stock-cli position analyze --question "分析持仓全部股票操作建议" --ai-config-id 1 --enable-tools --quiet
```

### 卖出/移除持仓

```bash
go-stock-cli position sell --code "000592.SZ" --quiet
```

## 建议

- 如需输出 JSON 结构数据，请使用 `--format json`。在 OpenClaw 中不建议开启，部分模型容易编码混乱。
- 可结合 `--stream` 获取流式输出。
- 若需生成“推荐股票列表”，请使用 `ai summary` 并开启 `--enable-tools`。
- 为避免日志和 CLI 输出混在一起，请使用 `--quiet` 参数。
