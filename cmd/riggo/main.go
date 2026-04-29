// cmd/claw/main.go
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/Ailoc/riggo/internal/engine"
	"github.com/Ailoc/riggo/internal/provider"
	"github.com/Ailoc/riggo/internal/tools"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")
	// 确保设置了 ZHIPU_API_KEY
	if os.Getenv("ZHIPU_API_KEY") == "" {
		log.Fatal("请先导出 ZHIPU_API_KEY 环境变量")
	}

	// 1. 获取工作区物理边界
	workDir, _ := os.Getwd()
	workDir = filepath.Dir(workDir)
	workDir = filepath.Dir(workDir) // 当前项目的根目录

	// 2. 初始化真实的大脑 (指向智谱 GLM-4.5，使用上一讲的 OpenAI 适配器)
	llmProvider := provider.NewZhipuOpenAIProvider("glm-4.5-air")

	// 3. 初始化真实的 Tool Registry
	registry := tools.NewRegistry()

	// 4. 将真实的 ReadFile 工具挂载到注册表中
	readFileTool := tools.NewReadFileTool(workDir)
	writeFileTool := tools.NewWriteFileTool(workDir)
	bashTool := tools.NewBashTool(workDir)
	registry.Register(readFileTool)
	registry.Register(writeFileTool)
	registry.Register(bashTool)

	// 5. 实例化核心引擎，由于任务简单，我们关闭思考阶段 (EnableThinking = false) 以加快速度
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	// 6. 下发一个必须通过真实工具才能完成的任务
	prompt := ` 请帮我执行以下操作：
	1. 用 bash 查看一下我当前电脑的 Go 版本。
	2. 帮我写一个简单的 helloworld.go 文件，输出 "Hello, go-tiny-claw!"。
	3. 用 bash 编译并运行这个 go 文件，确认它能正常工作。 `

	err := eng.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
