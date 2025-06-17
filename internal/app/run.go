/*
@Date: 2021/1/12 下午2:23
@Author: max.liu
@File : run
@Desc:
*/

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/maxliu9403/ProxyHub/internal/common"
	"github.com/maxliu9403/ProxyHub/internal/config"
	"github.com/maxliu9403/ProxyHub/internal/handler"
	"github.com/maxliu9403/ProxyHub/models"
	"github.com/maxliu9403/common/apiserver"
	"github.com/maxliu9403/common/apiserver/conf"
	"github.com/maxliu9403/common/logger"
	"github.com/maxliu9403/common/version"
	"github.com/spf13/cobra"
)

const projectName = "ProxyHub"

var (
	configFile string
	rootCmd    = &cobra.Command{
		Short: projectName,
		RunE: func(*cobra.Command, []string) error {
			return run()
		},
	}

	versionCommand = version.NewVerCommand(projectName)
	envCommand     = apiserver.NewConfigEnvCommand(config.G)
	initDB         = models.NewCreateDatabaseCommand(&configFile)
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "configs/dev.yaml", "configuration file path")
	rootCmd.AddCommand(versionCommand, envCommand, initDB)
}

func run() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = conf.LoadConfig(configFile, config.G)
	if err != nil {
		return fmt.Errorf("config file init failed: %s", err.Error())
	}

	// 数据表迁移，新增表时修改 AllTables
	m := apiserver.Migration(models.AllTables)
	// 初始化项目所需的 API，定时从配置中心获取最新 API 地址
	server := apiserver.CreateNewServer(ctx, config.G.APIConfig, m)
	defer server.Stop()

	logger.Debugf("%+v", config.G)

	group := server.AddGinGroup("")
	tra := server.GetTracer()
	handler.RegisterRouter(tra, group)

	// 以下是初始化 kafka 生产者和消费者，按需取消注释，或者删除代码
	// producer.NewProducer(config.G.Kafka)
	//
	// c, err := kafka.Default().NewConsumer()
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// 	return
	// }
	// if err = consumer.RunConsume(c); err != nil {
	// 	logger.Fatal(err.Error())
	// 	return
	// }

	// 初始化 validator 翻译器
	if err = common.InitTrans("zh"); err != nil {
		return fmt.Errorf("init trans failed, err: %v", err)
	}

	server.Start()
	return err
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
