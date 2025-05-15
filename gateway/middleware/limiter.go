package middleware

import (
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

var limiter *rate.Limiter

func InitLimiter() {
	rps := viper.GetInt("limiter.rps")
	burst := viper.GetInt("limiter.burst")

	limiter = rate.NewLimiter(rate.Limit(rps), burst)
}

func Allow() bool {
	return limiter.Allow()
}
