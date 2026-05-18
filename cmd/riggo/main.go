// cmd/claw/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ailoc/riggo/internal/engine"
	"github.com/Ailoc/riggo/internal/feishu"
	"github.com/Ailoc/riggo/internal/provider"
	"github.com/Ailoc/riggo/internal/tools"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	// 确保设置了 ZHIPU_API_KEY
	if os.Getenv("ZHIPU_API_KEY") == "" {
		log.Fatal("请先导出 ZHIPU_API_KEY 环境变量")
	}

	// 1. 获取工作区物理边界
	workDir, _ := os.Getwd()

	// 2. 初始化真实的大脑 (指向智谱 GLM-4.5，使用上一讲的 OpenAI 适配器)
	llmProvider := provider.NewZhipuOpenAIProvider("glm-4.5-air")

	// 3. 初始化真实的 Tool Registry
	registry := tools.NewRegistry()

	// 4. 将真实的 ReadFile 工具挂载到注册表中
	readFileTool := tools.NewReadFileTool(workDir)
	writeFileTool := tools.NewWriteFileTool(workDir)
	bashTool := tools.NewBashTool(workDir)
	registry.Register(tools.NewEditFileTool(workDir))

	registry.Register(readFileTool)
	registry.Register(writeFileTool)
	registry.Register(bashTool)

	// 5. 实例化核心引擎，由于任务简单，我们关闭思考阶段 (EnableThinking = false) 以加快速度
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 飞书模式：有环境变量时后台启动
	if os.Getenv("FEISHU_APP_ID") != "" && os.Getenv("FEISHU_APP_SECRET") != "" {
		bot := feishu.NewFeishuBot(eng)
		go func() {
			log.Println("🚀 飞书 WebSocket 长连接模式启动...")
			if err := bot.StartWebSocket(ctx); err != nil {
				log.Printf("❌ WebSocket 连接失败: %v\n", err)
			}
		}()
	}
}
