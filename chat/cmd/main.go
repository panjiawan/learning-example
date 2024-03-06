package main

import (
	"fmt"
	"github.com/panjiawan/note/chat/boot"
	"github.com/spf13/cobra"
)

func main() {
	rootCommand := &cobra.Command{}
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "输出版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("1.0.0")
		},
	}
	rootCommand.AddCommand(versionCommand)

	httpCommand := &cobra.Command{
		Use:   "http",
		Short: "启动http服务器",
		Long:  "启动http服务器\n$./cmd http etc配置路径 log日志路径",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Println("http启动参数有误: http etc配置路径 log日志路径")
				return
			}
			etcPath := args[0]
			logPath := args[1]
			boot.Start(etcPath, logPath)
		},
	}
	rootCommand.AddCommand(httpCommand)

	rootCommand.Execute()
}
