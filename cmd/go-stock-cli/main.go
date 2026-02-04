package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go-stock/backend/agent"
	"go-stock/backend/aitools"
	"go-stock/backend/bootstrap"
	"go-stock/backend/data"
	"go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/backend/util"
	"go-stock/internal/assets"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
)

type OutputFormat string

const (
	FormatMarkdown OutputFormat = "markdown"
	FormatJSON     OutputFormat = "json"
)

type GlobalFlags struct {
	DbPath  string
	Format  OutputFormat
	Timeout time.Duration
	Quiet   bool
}

type AIRequest struct {
	Question    string `json:"question"`
	Stock       string `json:"stock"`
	StockCode   string `json:"stock_code"`
	AiConfigID  int    `json:"ai_config_id"`
	SysPromptID int    `json:"sys_prompt_id"`
	EnableTools bool   `json:"enable_tools"`
	Thinking    bool   `json:"thinking"`
	Stream      bool   `json:"stream"`
	Format      string `json:"format"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "init":
		os.Exit(runInit(os.Args[2:]))
	case "ai":
		os.Exit(runAI(os.Args[2:]))
	case "agent":
		os.Exit(runAgent(os.Args[2:]))
	case "market":
		os.Exit(runMarket(os.Args[2:]))
	case "stock":
		os.Exit(runStock(os.Args[2:]))
	case "config":
		os.Exit(runConfig(os.Args[2:]))
	case "recommend":
		os.Exit(runRecommend(os.Args[2:]))
	case "serve":
		os.Exit(runServe(os.Args[2:]))
	case "-h", "--help", "help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("go-stock-cli 用法:")
	fmt.Println("  go-stock-cli init [--no-stock-data] [--update-basic-info]")
	fmt.Println("  go-stock-cli ai chat --question \"...\" [--stock \"\"] [--code \"\"] [--ai-config-id 1]")
	fmt.Println("  go-stock-cli ai summary [--question \"...\"] [--ai-config-id 1]")
	fmt.Println("  go-stock-cli agent chat --question \"...\" [--ai-config-id 1]")
	fmt.Println("  go-stock-cli market news [--source \"\"] [--limit 20]")
	fmt.Println("  go-stock-cli stock quote --code \"sh600000\"")
	fmt.Println("  go-stock-cli config export --file config.json")
	fmt.Println("  go-stock-cli config import --file config.json")
	fmt.Println("  go-stock-cli config add-ai --name \"\" --base-url \"\" --api-key \"\" --model \"\"")
	fmt.Println("  go-stock-cli config list-ai")
	fmt.Println("  go-stock-cli recommend list [--page 1] [--page-size 10]")
	fmt.Println("  go-stock-cli serve --listen 127.0.0.1:17688 [--token \"\"]")
	fmt.Println("")
	fmt.Println("全局参数:")
	fmt.Println("  --db       数据库路径 (默认 data/stock.db)")
	fmt.Println("  --format   输出格式 markdown|json (默认 markdown)")
	fmt.Println("  --timeout  超时秒数")
}

func parseGlobalFlags(fs *flag.FlagSet) (*GlobalFlags, *string, *int, *bool) {
	g := &GlobalFlags{}
	format := fs.String("format", "markdown", "输出格式 markdown|json")
	timeout := fs.Int("timeout", 0, "超时秒数")
	quiet := fs.Bool("quiet", false, "关闭控制台日志")
	fs.StringVar(&g.DbPath, "db", "data/stock.db", "数据库路径")
	return g, format, timeout, quiet
}

func finalizeGlobalFlags(g *GlobalFlags, format *string, timeout *int, quiet *bool) {
	g.Format = normalizeFormat(*format)
	if *timeout > 0 {
		g.Timeout = time.Duration(*timeout) * time.Second
	}
	if quiet != nil && *quiet {
		g.Quiet = true
		_ = os.Setenv("GO_STOCK_LOG_STDOUT", "0")
		logger.InitLogger()
	}
}

func normalizeFormat(v string) OutputFormat {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "json":
		return FormatJSON
	case "md":
		return FormatMarkdown
	default:
		return FormatMarkdown
	}
}

func initApp(g *GlobalFlags, opts bootstrap.Options) error {
	return bootstrap.Init(g.DbPath, assets.StockData(), opts)
}

func runInit(args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	noStock := fs.Bool("no-stock-data", false, "跳过基础股票数据导入")
	updateBasic := fs.Bool("update-basic-info", false, "启动后更新基础信息")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	err := initApp(g, bootstrap.Options{
		NoStockData:      *noStock,
		UpdateBasicInfo:  *updateBasic,
		AutoMigrateAsync: false,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	fmt.Println("初始化完成")
	return 0
}

func runAI(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "缺少子命令: chat 或 summary")
		return 1
	}
	switch args[0] {
	case "chat":
		return runAIChat(args[1:])
	case "summary":
		return runAISummary(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "未知 ai 子命令")
		return 1
	}
}

func runAIChat(args []string) int {
	fs := flag.NewFlagSet("ai chat", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	question := fs.String("question", "", "问题")
	stock := fs.String("stock", "", "股票名称")
	code := fs.String("code", "", "股票代码")
	aiConfigID := fs.Int("ai-config-id", 0, "AI配置ID")
	sysPromptID := fs.Int("sys-prompt-id", 0, "系统提示词ID")
	enableTools := fs.Bool("enable-tools", false, "启用工具调用")
	thinking := fs.Bool("think", false, "启用思考模式")
	stream := fs.Bool("stream", false, "流式输出")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if strings.TrimSpace(*question) == "" {
		fmt.Fprintln(os.Stderr, "question 不能为空")
		return 1
	}
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	configID := resolveAiConfigID(*aiConfigID)
	if configID == 0 {
		fmt.Fprintln(os.Stderr, "未找到AI配置，请先使用 config add-ai 或 config import")
		return 1
	}
	ctx := context.Background()
	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}
	var sysPromptPtr *int
	if *sysPromptID > 0 {
		sysPromptPtr = sysPromptID
	}
	ai := data.NewDeepSeekOpenAi(ctx, configID)
	tools := []data.Tool{}
	if *enableTools {
		tools = aitools.DefaultTools()
	}
	ch := ai.NewChatStream(*stock, *code, *question, sysPromptPtr, tools, *thinking)
	if *stream {
		_, _ = streamMapChannel(os.Stdout, ch, g.Format)
		return 0
	}
	_, result := collectMapChannel(ch)
	return writeOutput(os.Stdout, g.Format, result)
}

func runAISummary(args []string) int {
	fs := flag.NewFlagSet("ai summary", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	question := fs.String("question", "", "问题")
	aiConfigID := fs.Int("ai-config-id", 0, "AI配置ID")
	sysPromptID := fs.Int("sys-prompt-id", 0, "系统提示词ID")
	enableTools := fs.Bool("enable-tools", false, "启用工具调用")
	thinking := fs.Bool("think", false, "启用思考模式")
	stream := fs.Bool("stream", false, "流式输出")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	configID := resolveAiConfigID(*aiConfigID)
	if configID == 0 {
		fmt.Fprintln(os.Stderr, "未找到AI配置，请先使用 config add-ai 或 config import")
		return 1
	}
	ctx := context.Background()
	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}
	var sysPromptPtr *int
	if *sysPromptID > 0 {
		sysPromptPtr = sysPromptID
	}
	ai := data.NewDeepSeekOpenAi(ctx, configID)
	if *enableTools {
		ch := ai.NewSummaryStockNewsStreamWithTools(*question, sysPromptPtr, aitools.DefaultTools(), *thinking)
		if *stream {
			_, _ = streamMapChannel(os.Stdout, ch, g.Format)
			return 0
		}
		_, result := collectMapChannel(ch)
		return writeOutput(os.Stdout, g.Format, result)
	}
	ch := ai.NewSummaryStockNewsStream(*question, sysPromptPtr, *thinking)
	if *stream {
		_, _ = streamMapChannel(os.Stdout, ch, g.Format)
		return 0
	}
	_, result := collectMapChannel(ch)
	return writeOutput(os.Stdout, g.Format, result)
}

func runAgent(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "缺少子命令: chat")
		return 1
	}
	switch args[0] {
	case "chat":
		return runAgentChat(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "未知 agent 子命令")
		return 1
	}
}

func runAgentChat(args []string) int {
	fs := flag.NewFlagSet("agent chat", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	question := fs.String("question", "", "问题")
	aiConfigID := fs.Int("ai-config-id", 0, "AI配置ID")
	sysPromptID := fs.Int("sys-prompt-id", 0, "系统提示词ID")
	stream := fs.Bool("stream", false, "流式输出")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if strings.TrimSpace(*question) == "" {
		fmt.Fprintln(os.Stderr, "question 不能为空")
		return 1
	}
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	configID := resolveAiConfigID(*aiConfigID)
	if configID == 0 {
		fmt.Fprintln(os.Stderr, "未找到AI配置，请先使用 config add-ai 或 config import")
		return 1
	}
	var sysPromptPtr *int
	if *sysPromptID > 0 {
		sysPromptPtr = sysPromptID
	}
	if g.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
		defer cancel()
		_ = ctx
	}
	ch := agent.NewStockAiAgentApi().Chat(*question, configID, sysPromptPtr)
	if *stream {
		_, result := streamAgentChannel(os.Stdout, ch, g.Format)
		_ = result
		return 0
	}
	_, result := collectAgentChannel(ch)
	return writeOutput(os.Stdout, g.Format, result)
}

func runMarket(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "缺少子命令: news")
		return 1
	}
	switch args[0] {
	case "news":
		return runMarketNews(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "未知 market 子命令")
		return 1
	}
}

func runMarketNews(args []string) int {
	fs := flag.NewFlagSet("market news", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	source := fs.String("source", "", "来源")
	limit := fs.Int("limit", 20, "数量")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	_ = data.NewMarketNewsApi().TelegraphList(30)
	news := data.NewMarketNewsApi().GetNewsList(*source, *limit)
	if g.Format == FormatMarkdown {
		content := util.MarkdownTableWithTitle("市场快讯", news)
		fmt.Println(content)
		return 0
	}
	return writeOutput(os.Stdout, g.Format, map[string]any{"data": news})
}

func runStock(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "缺少子命令: quote")
		return 1
	}
	switch args[0] {
	case "quote":
		return runStockQuote(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "未知 stock 子命令")
		return 1
	}
}

func runStockQuote(args []string) int {
	fs := flag.NewFlagSet("stock quote", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	code := fs.String("code", "", "股票代码")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if strings.TrimSpace(*code) == "" {
		fmt.Fprintln(os.Stderr, "code 不能为空")
		return 1
	}
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	info, err := data.NewStockDataApi().GetStockCodeRealTimeData(*code)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	if g.Format == FormatMarkdown {
		content := util.MarkdownTableWithTitle("股票行情", info)
		fmt.Println(content)
		return 0
	}
	return writeOutput(os.Stdout, g.Format, map[string]any{"data": info})
}

func runConfig(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "缺少子命令: export/import/add-ai/list-ai")
		return 1
	}
	switch args[0] {
	case "export":
		return runConfigExport(args[1:])
	case "import":
		return runConfigImport(args[1:])
	case "add-ai":
		return runConfigAddAI(args[1:])
	case "list-ai":
		return runConfigListAI(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "未知 config 子命令")
		return 1
	}
}

func runRecommend(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "缺少子命令: list")
		return 1
	}
	switch args[0] {
	case "list":
		return runRecommendList(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "未知 recommend 子命令")
		return 1
	}
}

func runRecommendList(args []string) int {
	fs := flag.NewFlagSet("recommend list", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	page := fs.Int("page", 1, "页码")
	pageSize := fs.Int("page-size", 10, "每页数量")
	stockCode := fs.String("stock-code", "", "股票代码")
	stockName := fs.String("stock-name", "", "股票名称")
	bkCode := fs.String("bk-code", "", "板块代码")
	bkName := fs.String("bk-name", "", "板块名称")
	startDate := fs.String("start-date", "", "开始日期")
	endDate := fs.String("end-date", "", "结束日期")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	query := &models.AiRecommendStocksQuery{
		StockCode: *stockCode,
		StockName: *stockName,
		BkCode:    *bkCode,
		BkName:    *bkName,
		StartDate: *startDate,
		EndDate:   *endDate,
		Page:      *page,
		PageSize:  *pageSize,
	}
	pageData, err := data.NewAiRecommendStocksService().GetAiRecommendStocksList(query)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	if g.Format == FormatMarkdown {
		content := util.MarkdownTableWithTitle("AI推荐股票列表", pageData.List)
		fmt.Println(content)
		return 0
	}
	return writeOutput(os.Stdout, g.Format, pageData)
}

func runConfigExport(args []string) int {
	fs := flag.NewFlagSet("config export", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	file := fs.String("file", "config.json", "输出文件")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	content := data.NewSettingsApi().Export()
	if err := os.WriteFile(*file, []byte(content), os.ModePerm); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	fmt.Println("导出完成:", *file)
	return 0
}

func runConfigImport(args []string) int {
	fs := flag.NewFlagSet("config import", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	file := fs.String("file", "", "导入文件")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if strings.TrimSpace(*file) == "" {
		fmt.Fprintln(os.Stderr, "file 不能为空")
		return 1
	}
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	raw, err := os.ReadFile(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	cfg := &data.SettingConfig{Settings: &data.Settings{}}
	if err := json.Unmarshal(raw, cfg); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	msg := data.UpdateConfig(cfg)
	fmt.Println(msg)
	return 0
}

func runConfigAddAI(args []string) int {
	fs := flag.NewFlagSet("config add-ai", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	name := fs.String("name", "", "配置名称")
	baseURL := fs.String("base-url", "", "Base URL")
	apiKey := fs.String("api-key", "", "API Key")
	model := fs.String("model", "", "模型名称")
	temperature := fs.Float64("temperature", 0.1, "温度")
	maxTokens := fs.Int("max-tokens", 4096, "最大tokens")
	timeOut := fs.Int("ai-timeout", 60, "请求超时秒数")
	httpProxy := fs.String("http-proxy", "", "HTTP代理")
	httpProxyEnabled := fs.Bool("http-proxy-enabled", false, "启用HTTP代理")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if *baseURL == "" || *apiKey == "" || *model == "" {
		fmt.Fprintln(os.Stderr, "base-url、api-key、model 不能为空")
		return 1
	}
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	cfg := data.GetSettingConfig()
	if cfg.Settings == nil {
		cfg.Settings = &data.Settings{}
	}
	cfg.OpenAiEnable = true
	cfg.AiConfigs = append(cfg.AiConfigs, &data.AIConfig{
		Name:             *name,
		BaseUrl:          *baseURL,
		ApiKey:           *apiKey,
		ModelName:        *model,
		Temperature:      *temperature,
		MaxTokens:        *maxTokens,
		TimeOut:          *timeOut,
		HttpProxy:        *httpProxy,
		HttpProxyEnabled: *httpProxyEnabled,
	})
	msg := data.UpdateConfig(cfg)
	fmt.Println(msg)
	return 0
}

func runConfigListAI(args []string) int {
	fs := flag.NewFlagSet("config list-ai", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	cfg := data.GetSettingConfig()
	if g.Format == FormatMarkdown {
		content := util.MarkdownTableWithTitle("AI配置列表", cfg.AiConfigs)
		fmt.Println(content)
		return 0
	}
	return writeOutput(os.Stdout, g.Format, map[string]any{"data": cfg.AiConfigs})
}

func runServe(args []string) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	g, format, timeout, quiet := parseGlobalFlags(fs)
	listen := fs.String("listen", "127.0.0.1:17688", "监听地址")
	token := fs.String("token", "", "鉴权token")
	if hasHelp(args) {
		fs.Usage()
		return 0
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}
	finalizeGlobalFlags(g, format, timeout, quiet)
	if err := initApp(g, bootstrap.Options{NoStockData: true, AutoMigrateAsync: false}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	fmt.Printf("HTTP服务启动: %s\n", *listen)
	return startHTTPServer(*listen, *token, g)
}

func startHTTPServer(addr, token string, g *GlobalFlags) int {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if !checkToken(r, token) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		writeJSON(w, map[string]any{"ok": true})
	})
	mux.HandleFunc("/ai/chat", func(w http.ResponseWriter, r *http.Request) {
		if !checkToken(r, token) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		handleAIChatHTTP(w, r, g)
	})
	mux.HandleFunc("/ai/summary", func(w http.ResponseWriter, r *http.Request) {
		if !checkToken(r, token) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		handleAISummaryHTTP(w, r, g)
	})
	mux.HandleFunc("/agent/chat", func(w http.ResponseWriter, r *http.Request) {
		if !checkToken(r, token) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		handleAgentChatHTTP(w, r, g)
	})
	mux.HandleFunc("/market/news", func(w http.ResponseWriter, r *http.Request) {
		if !checkToken(r, token) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		handleMarketNewsHTTP(w, r, g)
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	return 0
}

func handleAIChatHTTP(w http.ResponseWriter, r *http.Request, g *GlobalFlags) {
	req, err := decodeAIRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	configID := resolveAiConfigID(req.AiConfigID)
	if configID == 0 {
		http.Error(w, "AI配置不存在", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}
	var sysPromptPtr *int
	if req.SysPromptID > 0 {
		sysPromptPtr = &req.SysPromptID
	}
	ai := data.NewDeepSeekOpenAi(ctx, configID)
	tools := []data.Tool{}
	if req.EnableTools {
		tools = aitools.DefaultTools()
	}
	ch := ai.NewChatStream(req.Stock, req.StockCode, req.Question, sysPromptPtr, tools, req.Thinking)
	if req.Stream {
		streamSSE(w, ch)
		return
	}
	_, result := collectMapChannel(ch)
	writeJSON(w, result)
}

func handleAISummaryHTTP(w http.ResponseWriter, r *http.Request, g *GlobalFlags) {
	req, err := decodeAIRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	configID := resolveAiConfigID(req.AiConfigID)
	if configID == 0 {
		http.Error(w, "AI配置不存在", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}
	var sysPromptPtr *int
	if req.SysPromptID > 0 {
		sysPromptPtr = &req.SysPromptID
	}
	ai := data.NewDeepSeekOpenAi(ctx, configID)
	var ch <-chan map[string]any
	if req.EnableTools {
		ch = ai.NewSummaryStockNewsStreamWithTools(req.Question, sysPromptPtr, aitools.DefaultTools(), req.Thinking)
	} else {
		ch = ai.NewSummaryStockNewsStream(req.Question, sysPromptPtr, req.Thinking)
	}
	if req.Stream {
		streamSSE(w, ch)
		return
	}
	_, result := collectMapChannel(ch)
	writeJSON(w, result)
}

func handleAgentChatHTTP(w http.ResponseWriter, r *http.Request, g *GlobalFlags) {
	req, err := decodeAIRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Question) == "" {
		http.Error(w, "question 不能为空", http.StatusBadRequest)
		return
	}
	configID := resolveAiConfigID(req.AiConfigID)
	if configID == 0 {
		http.Error(w, "AI配置不存在", http.StatusBadRequest)
		return
	}
	var sysPromptPtr *int
	if req.SysPromptID > 0 {
		sysPromptPtr = &req.SysPromptID
	}
	ch := agent.NewStockAiAgentApi().Chat(req.Question, configID, sysPromptPtr)
	if req.Stream {
		streamAgentSSE(w, ch)
		return
	}
	_, result := collectAgentChannel(ch)
	writeJSON(w, result)
}

func handleMarketNewsHTTP(w http.ResponseWriter, r *http.Request, g *GlobalFlags) {
	source := r.URL.Query().Get("source")
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}
	_ = data.NewMarketNewsApi().TelegraphList(30)
	news := data.NewMarketNewsApi().GetNewsList(source, limit)
	writeJSON(w, map[string]any{"data": news})
}

func decodeAIRequest(r *http.Request) (*AIRequest, error) {
	if r.Method != http.MethodPost {
		return nil, errors.New("只支持POST")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	req := &AIRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		return nil, err
	}
	req.Format = strings.ToLower(strings.TrimSpace(req.Format))
	return req, nil
}

func streamSSE(w http.ResponseWriter, ch <-chan map[string]any) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, _ := w.(http.Flusher)
	for msg := range ch {
		writeSSE(w, "delta", msg)
		if flusher != nil {
			flusher.Flush()
		}
	}
	writeSSE(w, "done", map[string]any{"ok": true})
	if flusher != nil {
		flusher.Flush()
	}
}

func streamAgentSSE(w http.ResponseWriter, ch <-chan *schema.Message) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, _ := w.(http.Flusher)
	for msg := range ch {
		writeSSE(w, "delta", msg)
		if flusher != nil {
			flusher.Flush()
		}
	}
	writeSSE(w, "done", map[string]any{"ok": true})
	if flusher != nil {
		flusher.Flush()
	}
}

func writeSSE(w http.ResponseWriter, event string, data any) {
	payload, _ := json.Marshal(data)
	_, _ = fmt.Fprintf(w, "event: %s\n", event)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", payload)
}

func checkToken(r *http.Request, token string) bool {
	if token == "" {
		return true
	}
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		auth = strings.TrimPrefix(auth, "Bearer ")
	}
	if auth == token {
		return true
	}
	if r.Header.Get("X-Token") == token {
		return true
	}
	return false
}

func collectMapChannel(ch <-chan map[string]any) (string, map[string]any) {
	var b strings.Builder
	meta := map[string]any{}
	for msg := range ch {
		if v, ok := msg["content"].(string); ok {
			b.WriteString(v)
		}
		if v, ok := msg["model"]; ok {
			meta["model"] = v
		}
		if v, ok := msg["chatId"]; ok {
			meta["chat_id"] = v
		}
		if v, ok := msg["question"]; ok {
			meta["question"] = v
		}
	}
	meta["content"] = b.String()
	return b.String(), meta
}

func streamMapChannel(w io.Writer, ch <-chan map[string]any, format OutputFormat) (string, map[string]any) {
	var b strings.Builder
	meta := map[string]any{}
	for msg := range ch {
		if v, ok := msg["content"].(string); ok {
			b.WriteString(v)
			if format == FormatMarkdown {
				_, _ = fmt.Fprint(w, v)
			}
		}
		if format == FormatJSON {
			_ = writeJSONLine(w, map[string]any{"event": "delta", "data": msg})
		}
		if v, ok := msg["model"]; ok {
			meta["model"] = v
		}
		if v, ok := msg["chatId"]; ok {
			meta["chat_id"] = v
		}
		if v, ok := msg["question"]; ok {
			meta["question"] = v
		}
	}
	meta["content"] = b.String()
	if format == FormatJSON {
		_ = writeJSONLine(w, map[string]any{"event": "done"})
	} else {
		_, _ = fmt.Fprintln(w)
	}
	return b.String(), meta
}

func collectAgentChannel(ch <-chan *schema.Message) (string, map[string]any) {
	var b strings.Builder
	meta := map[string]any{}
	for msg := range ch {
		if msg == nil {
			continue
		}
		if msg.ReasoningContent != "" {
			b.WriteString(msg.ReasoningContent)
		}
		if msg.Content != "" {
			b.WriteString(msg.Content)
		}
		if msg.ResponseMeta != nil {
			meta["finish_reason"] = msg.ResponseMeta.FinishReason
		}
	}
	meta["content"] = b.String()
	return b.String(), meta
}

func streamAgentChannel(w io.Writer, ch <-chan *schema.Message, format OutputFormat) (string, map[string]any) {
	var b strings.Builder
	meta := map[string]any{}
	for msg := range ch {
		if msg == nil {
			continue
		}
		if msg.ReasoningContent != "" {
			b.WriteString(msg.ReasoningContent)
			if format == FormatMarkdown {
				_, _ = fmt.Fprint(w, msg.ReasoningContent)
			}
		}
		if msg.Content != "" {
			b.WriteString(msg.Content)
			if format == FormatMarkdown {
				_, _ = fmt.Fprint(w, msg.Content)
			}
		}
		if format == FormatJSON {
			_ = writeJSONLine(w, map[string]any{"event": "delta", "data": msg})
		}
		if msg.ResponseMeta != nil {
			meta["finish_reason"] = msg.ResponseMeta.FinishReason
		}
	}
	meta["content"] = b.String()
	if format == FormatJSON {
		_ = writeJSONLine(w, map[string]any{"event": "done"})
	} else {
		_, _ = fmt.Fprintln(w)
	}
	return b.String(), meta
}

func resolveAiConfigID(input int) int {
	cfg := data.GetSettingConfig()
	if input > 0 {
		for _, item := range cfg.AiConfigs {
			if int(item.ID) == input {
				return input
			}
		}
		return 0
	}
	if len(cfg.AiConfigs) == 0 {
		return 0
	}
	return int(cfg.AiConfigs[0].ID)
}

func hasHelp(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

func writeOutput(w io.Writer, format OutputFormat, data any) int {
	switch format {
	case FormatJSON:
		if err := writeJSON(w, data); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
	default:
		if m, ok := data.(map[string]any); ok {
			if content, ok := m["content"].(string); ok {
				fmt.Fprintln(w, content)
				return 0
			}
		}
		fmt.Fprintln(w, data)
	}
	return 0
}

func writeJSON(w io.Writer, data any) error {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.Header().Set("Content-Type", "application/json")
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}

func writeJSONLine(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}

// readStdin is reserved for future interactive usage.
