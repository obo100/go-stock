package aitools

import (
	"go-stock/backend/data"
)

func AddTools(tools []data.Tool) []data.Tool {
	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name: "SearchStockByIndicators",
			Description: "根据自然语言筛选股票，返回自然语言选股条件要求的股票所有相关数据。输入股票名称可以获取当前股票最新的股价交易数据和基础财务指标信息，多个股票名称使用,分隔。" +
				"例如:查看涨停板：涨停板，按涨幅从高到低排序。" +
				"例如:查看跌停板：跌停板，按跌幅从高到低排序。" +
				"例如:查看龙虎榜：龙虎榜，按涨幅从高到低排序。" +
				"例如:查看昨日龙虎榜：昨日龙虎榜。" +
				"例如:查看板块龙头行情：板块/概念龙头，按涨幅从高到低排序。" +
				"例如:查看板块龙头行情：龙头股，按成交量从高到低排序。" +
				"例如:查看技术指标：上海贝岭,macd,rsi,kdj,boll,5日均线,14日均线,30日均线,60日均线,成交量,OBV,EMA。" +
				"例如:查看近期趋势：量比连续2天>1，主力连续2日净流入且递增，主力净额>3000万元，行业，股价在20日线上。按成交量从高到低排序。" +
				"例如:当日成交量 ≥ 近 5 日平均成交量 ×1.5，收盘价 ≥ 20 日均线，20 日均线 ≥ 60 日均线，当日涨幅 3%-7%， 3日主力资金净流入累计≥5000 万元，当日换手率 5%-15%，筹码集中度（90% 筹码峰）≤15%，非创业板非科创板非ST非北交所，行业。按成交量从高到低排序。" +
				"例如:查看有潜力的成交量爆发股：最近7日成交量量比大于3，出现过一次，非ST。按成交量从高到低排序。" +
				"例如:超短线策略：当日成交量大于前一日成交量的1.8倍;当日最高价创60日新高当日收盘价大于5日均线;当日为阳线;股价小于200;" +
				"例1：创新药,半导体;PE<30;净利润增长率>50%。 按成交量从高到低排序。" +
				"例2：上证指数,科创50。 " +
				"例3：长电科技,上海贝岭。" +
				"例4：长电科技,上海贝岭;KDJ,MACD,RSI,BOLL,主力资金。" +
				"例5：换手率大于3%小于25%.量比1以上. 10日内有过涨停.股价处于峰值的二分之一以下.流通股本<100亿.当日和连续四日净流入;股价在20日均线以上.分时图股价在均线之上.热门板块下涨幅领先的A股. 当日量能20000手以上.沪深个股.近一年市盈率波动小于150%.MACD金叉;不要ST股及不要退市股，非北交所，每股收益>0。按成交量从高到低排序。" +
				"例6：沪深主板.流通市值小于100亿.市值大于10亿.60分钟dif大于dea.60分钟skdj指标k值大于d值.skdj指标k值小于90.换手率大于3%.成交额大于1亿元.量比大于2.涨幅大于2%小于7%.股价大于5小于50.创业板.10日均线大于20日均线;不要ST股及不要退市股;不要北交所;不要科创板;不要创业板。按成交量从高到低排序。" +
				"例7：股价在20日线上，一月之内涨停次数>=1，量比大于1，换手率大于3%。按成交量从高到低排序。" +
				"例8：基本条件：前期有爆量，回调到 10 日线，当日是缩量阴线，均线趋势向上。;优选条件：一月之内涨停次数>=1。按成交量从高到低排序。",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"words": map[string]any{
						"type": "string",
						"description": "选股自然语言。" +
							"例如:查看涨停板：涨停板，按涨幅从高到低排序。" +
							"例如:查看跌停板：跌停板，按跌幅从高到低排序。" +
							"例如:查看龙虎榜：龙虎榜，按涨幅从高到低排序。" +
							"例如:查看昨日龙虎榜：昨日龙虎榜。" +
							"例如:查看板块龙头行情：板块/概念龙头，按涨幅从高到低排序。" +
							"例如:查看板块龙头行情：龙头股，按成交量从高到低排序。" +
							"例如:查看技术指标：上海贝岭,macd,rsi,kdj,boll,5日均线,14日均线,30日均线,60日均线,成交量,OBV,EMA。" +
							"例如:查看近期趋势：量比连续2天>1，主力连续2日净流入且递增，主力净额>3000万元，行业，股价在20日线上。按成交量从高到低排序。" +
							"例如:当日成交量 ≥ 近 5 日平均成交量 ×1.5，收盘价 ≥ 20 日均线，20 日均线 ≥ 60 日均线，当日涨幅 3%-7%， 3日主力资金净流入累计≥5000 万元，当日换手率 5%-15%，筹码集中度（90% 筹码峰）≤15%，非创业板非科创板非ST非北交所，行业。按成交量从高到低排序。" +
							"例如:查看有潜力的成交量爆发股：最近7日成交量量比大于3，出现过一次，非ST。按成交量从高到低排序。" +
							"例如:超短线策略：当日成交量大于前一日成交量的1.8倍;当日最高价创60日新高当日收盘价大于5日均线;当日为阳线;股价小于200;" +
							"例1：创新药,半导体;PE<30;净利润增长率>50%。 按成交量从高到低排序。" +
							"例2：上证指数,科创50。 " +
							"例3：长电科技,上海贝岭。" +
							"例4：长电科技,上海贝岭;KDJ,MACD,RSI,BOLL,主力资金。" +
							"例5：换手率大于3%小于25%.量比1以上. 10日内有过涨停.股价处于峰值的二分之一以下.流通股本<100亿.当日和连续四日净流入;股价在20日均线以上.分时图股价在均线之上.热门板块下涨幅领先的A股. 当日量能20000手以上.沪深个股.近一年市盈率波动小于150%.MACD金叉;不要ST股及不要退市股，非北交所，每股收益>0。按成交量从高到低排序。" +
							"例6：沪深主板.流通市值小于100亿.市值大于10亿.60分钟dif大于dea.60分钟skdj指标k值大于d值.skdj指标k值小于90.换手率大于3%.成交额大于1亿元.量比大于2.涨幅大于2%小于7%.股价大于5小于50.创业板.10日均线大于20日均线;不要ST股及不要退市股;不要北交所;不要科创板;不要创业板。按成交量从高到低排序。" +
							"例7：股价在20日线上，一月之内涨停次数>=1，量比大于1，换手率大于3%。按成交量从高到低排序。" +
							"例8：基本条件：前期有爆量，回调到 10 日线，当日是缩量阴线，均线趋势向上。;优选条件：一月之内涨停次数>=1。按成交量从高到低排序。",
					},
				},
				Required: []string{"words"},
			},
		},
	})

	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name: "SearchBk",
			Description: "根据自然语言查询板块/概念/指数整体数据。" +
				"例如:存储芯片，成分股" +
				"例如:查看指数：上证指数。" +
				"例如:查看存储芯片板块：存储芯片。" +
				"例如:查看概念板块排名：今日涨幅前5的概念板块。" +
				"例如:查看概念板块排名：今日净流入前5的概念板块。" +
				"例如:查看板块/概念排名数据：今日主力净流出前15的概念板块。" +
				"例如:查看板块板块/概念：今日成交量前15的概念板块。" +
				"例如:通过市盈率查询板块：当前市盈率介于30-50的板块/概念。",

			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"words": map[string]any{
						"type": "string",
						"description": "板块/概念数据查询自然语言。" +
							"例如:存储芯片，成分股" +
							"例如:查看指数：上证指数。" +
							"例如:查看存储芯片板块：存储芯片。" +
							"例如:查看概念板块排名：今日涨幅前5的概念板块。" +
							"例如:查看概念板块排名：今日净流入前5的概念板块。" +
							"例如:查看板块/概念排名数据：今日主力净流出前15的概念板块。" +
							"例如:查看板块板块/概念：今日成交量前15的概念板块。" +
							"例如:通过市盈率查询板块：当前市盈率介于30-50的板块/概念。",
					},
				},
				Required: []string{"words"},
			},
		},
	})

	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "GetStockKLine",
			Description: "获取股票日K线数据。",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"days": map[string]any{
						"type":        "string",
						"description": "日K数据条数",
					},
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码（A股：sh,sz开头;港股hk开头,美股：us开头）",
					},
				},
				Required: []string{"days", "stockCode"},
			},
		},
	})

	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "InteractiveAnswer",
			Description: "获取投资者与上市公司互动问答的数据,反映当前投资者关注的热点问题",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"page": map[string]any{
						"type":        "string",
						"description": "分页号",
					},
					"pageSize": map[string]any{
						"type":        "string",
						"description": "分页大小",
					},
					"keyWord": map[string]any{
						"type":        "string",
						"description": "搜索关键词（可输入股票名称或者当前热门板块/行业/概念/标的/事件等）",
					},
				},
				Required: []string{"page", "pageSize"},
			},
		},
	})

	//tools = append(tools, data.Tool{
	//	Type: "function",
	//	Function: data.ToolFunction{
	//		Name:        "QueryBKDictInfo",
	//		Description: "获取所有板块/行业名称或者代码(bkCode,bkName)",
	//	},
	//})

	//tools = append(tools, data.Tool{
	//	Type: "function",
	//	Function: data.ToolFunction{
	//		Name:        "GetIndustryResearchReport",
	//		Description: "获取行业/板块研究报告,请先使用QueryBKDictInfo工具获取行业代码，然后输入行业代码调用",
	//		Parameters: data.FunctionParameters{
	//			Type: "object",
	//			Properties: map[string]any{
	//				"bkCode": map[string]any{
	//					"type":        "string",
	//					"description": "板块/行业代码",
	//				},
	//			},
	//			Required: []string{"bkCode"},
	//		},
	//	},
	//})

	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "GetStockResearchReport",
			Description: "获取市场分析师的股票研究报告",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码",
					},
				},
				Required: []string{"stockCode"},
			},
		},
	})

	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "HotStrategyTable",
			Description: "获取当前热门选股策略",
		},
	})

	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "HotStockTable",
			Description: "当前热门股票排名",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"pageSize": map[string]any{
						"type":        "string",
						"description": "分页大小",
					},
				},
				Required: []string{"pageSize"},
			},
		},
	})

	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "GetStockMoneyData",
			Description: "今日股票资金流入排名",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"pageSize": map[string]any{
						"type":        "string",
						"description": "分页大小",
					},
				},
				Required: []string{"pageSize"},
			},
		},
	})
	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "GetStockConceptInfo",
			Description: "获取股票所属概念详细信息",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"code": map[string]any{
						"type":        "string",
						"description": "股票代码,如：601138.SH。注意 上海证券交易所股票以.SH结尾，深圳证券交易所股票以.SZ结尾，港股股票以.HK结尾，北交所股票以.BJ结尾，",
					},
				},
				Required: []string{"code"},
			},
		},
	})

	//CreateAiRecommendStocks
	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "CreateAiRecommendStocks",
			Description: "创建/保存AI推荐股票记录",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"modelName": map[string]any{
						"type":        "string",
						"description": "模型名称",
					},
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码,如：601138.SH。注意 上海证券交易所股票以.SH结尾，深圳证券交易所股票以.SZ结尾，港股股票以.HK结尾，北交所股票以.BJ结尾，",
					},
					"stockName": map[string]any{
						"type":        "string",
						"description": "股票名称",
					},
					"bkCode": map[string]any{
						"type":        "string",
						"description": "板块/行业代码",
					},
					"bkName": map[string]any{
						"type":        "string",
						"description": "板块/概念/行业名称",
					},
					"stockPrice": map[string]any{
						"type":        "string",
						"description": "推荐时股票价格",
					},
					"stockPrePrice": map[string]any{
						"type":        "string",
						"description": "前一交易日股票价格",
					},
					"stockClosePrice": map[string]any{
						"type":        "string",
						"description": "推荐时股票收盘价格",
					},
					"recommendReason": map[string]any{
						"type":        "string",
						"description": "推荐理由/驱动因素/逻辑",
					},
					"recommendBuyPrice": map[string]any{
						"type":        "string",
						"description": "ai建议买入价区间最低价和最高价之间用`-`分隔",
					},
					"recommendStopProfitPrice": map[string]any{
						"type":        "string",
						"description": "ai建议止盈价区间最低价和最高价之间用`-`分隔",
					},
					"recommendStopLossPrice": map[string]any{
						"type":        "string",
						"description": "ai建议止损价",
					},
					"riskRemarks": map[string]any{
						"type":        "string",
						"description": "风险提示",
					},
					"remarks": map[string]any{
						"type":        "string",
						"description": "操作总结/备注",
					},
				},
				Required: []string{"stockCode", "stockName", "bkName"},
			},
		},
	})

	//BatchCreateAiRecommendStocks
	tools = append(tools, data.Tool{
		Type: "function",
		Function: data.ToolFunction{
			Name:        "BatchCreateAiRecommendStocks",
			Description: "批量创建/保存AI推荐股票记录，建议每次批量保存5条记录",
			Parameters: &data.FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stocks": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"modelName": map[string]any{
									"type":        "string",
									"description": "模型名称",
								},
								"stockCode": map[string]any{
									"type":        "string",
									"description": "股票代码,如：601138.SH。注意 上海证券交易所股票以.SH结尾，深圳证券交易所股票以.SZ结尾，港股股票以.HK结尾，北交所股票以.BJ结尾，",
								},
								"stockName": map[string]any{
									"type":        "string",
									"description": "股票名称",
								},
								"bkCode": map[string]any{
									"type":        "string",
									"description": "板块/行业代码",
								},
								"bkName": map[string]any{
									"type":        "string",
									"description": "板块/概念/行业名称",
								},
								"stockPrice": map[string]any{
									"type":        "string",
									"description": "推荐时股票价格",
								},
								"stockPrePrice": map[string]any{
									"type":        "string",
									"description": "前一交易日股票价格",
								},
								"stockClosePrice": map[string]any{
									"type":        "string",
									"description": "推荐时股票收盘价格",
								},
								"recommendReason": map[string]any{
									"type":        "string",
									"description": "推荐理由/驱动因素/逻辑",
								},
								"recommendBuyPrice": map[string]any{
									"type":        "string",
									"description": "ai建议买入价区间最低价和最高价之间用`-`分隔",
								},
								"recommendStopProfitPrice": map[string]any{
									"type":        "string",
									"description": "ai建议止盈价区间最低价和最高价之间用`-`分隔",
								},
								"recommendStopLossPrice": map[string]any{
									"type":        "string",
									"description": "ai建议止损价",
								},
								"riskRemarks": map[string]any{
									"type":        "string",
									"description": "风险提示",
								},
								"remarks": map[string]any{
									"type":        "string",
									"description": "操作总结/备注",
								},
							},
						},
					},
				},

				Required: []string{"stockCode", "stockName", "bkName"},
			},
		},
	})

	return tools
}

// DefaultTools returns the default tool list.
func DefaultTools() []data.Tool {
	return AddTools(nil)
}
