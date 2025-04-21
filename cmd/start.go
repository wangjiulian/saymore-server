package cmd

import (
	"bytes"
	"com.say.more.server/config"
	"com.say.more.server/internal/app"
	"com.say.more.server/internal/app/constant"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	localConfig string

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start SayMore Server",
		Run: func(cmd *cobra.Command, args []string) {
			online()
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&localConfig, "CONFIG", "c", "", "config path: /opt/local.toml")
}

func online() {
	data, err := ioutil.ReadFile(localConfig)
	if err != nil {
		fmt.Println("failed to get config data error: ", err)
		return
	}
	cfg, err := readConfig(data)
	setDefaultConfig(cfg)
	if err != nil {
		fmt.Println("failed to read config error: ", err)
		return
	}
	run(cfg)
}

func setDefaultConfig(cfg *config.Config) {
	if cfg.Course == nil {
		// Set default values for course cancellation
		cfg.Course = &config.Course{
			CancelInterval: constant.CancelCourseDefaultInterval,
			CancelRefund:   constant.CancelCourseDefaultRefund,
			CancelRule:     constant.CancelCourseDefaultRule,
			CourseUnit:     constant.CourseUnitDefault,
		}
	}
	if cfg.Course.CancelInterval <= 0 {
		cfg.Course.CancelInterval = constant.CancelCourseDefaultInterval
	}
	if cfg.Course.CancelRefund <= 0 {
		cfg.Course.CancelRefund = constant.CancelCourseDefaultRefund
	}
	if cfg.Course.CancelRule == "" {
		cfg.Course.CancelRule = constant.CancelCourseDefaultRule
	}
	if cfg.Course.CourseUnit <= 0 {
		cfg.Course.CourseUnit = constant.CourseUnitDefault
	}
}

func run(cfg *config.Config) {
	app.Run(cfg)
}

func readConfig(data []byte) (*config.Config, error) {
	v := viper.New()
	v.SetConfigType("toml")
	reader := bytes.NewReader(data)
	err := v.ReadConfig(reader)
	if err != nil {
		return nil, err
	}
	cfg := config.NewConfig()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
