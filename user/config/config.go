package config

import "github.com/spf13/viper"

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置失败：" + err.Error())
	} else {
		println("配置读取成功！使用配置文件:", viper.ConfigFileUsed())
	}
}
