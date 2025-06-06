/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	flagLevel            = "log.level"
	flagDisableCaller    = "log.disable-caller"
	flagFormat           = "log.format"
	flagName             = "log.name"
	flagEnableStacktrace = "log.enable-stacktrace"
	flagEnableColor      = "log.enable-color"
	flagOutputPaths      = "log.output-paths"
	flagErrorOutput      = "log.error-output"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// LogOption contains configuration items related to log.
type LogOption struct {
	Level            string   `json:"level"`
	Format           string   `json:"format"`
	DisableCaller    bool     `json:"disable-caller"`
	LogName          string   `json:"name"`
	EnableStacktrace bool     `json:"disableStacktrace"`
	EnableColor      bool     `json:"enableColor"`
	OutputPaths      []string `json:"outputPaths"`
	ErrorOutput      []string `json:"errorOutput"`
}

// Name flag set group name.
func (o *LogOption) Name() string {
	return "logs"
}

// NewLogOptions creates an LogOption object with default parameters.
func NewLogOptions() *LogOption {
	return &LogOption{
		Level:  zapcore.InfoLevel.String(),
		Format: consoleFormat,
		//EnableColor: true,
	}
}

func (o *LogOption) ApplyTo() *LogOption {
	return o
}

// Validate the options fields.
func (o *LogOption) Validate() []error {
	// 用于属性校验
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, fmt.Errorf("%s option: %v", o.Name(), err))
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("%s option: not a valid log format: %q", o.Name(), o.Format))
	}
	// 两个参数只能有一个有效
	if o.LogName != "" && len(o.OutputPaths) > 0 {
		errs = append(errs, fmt.Errorf("not a valid format, only one: %s,%v", o.Name(), o.OutputPaths))
	}

	return errs
}

// Complete set default LogOption.
func (o *LogOption) Complete() error {

	//o.LogName = viper.GetString(flagName)
	//o.Level = viper.GetString(flagLevel)
	//o.Format = viper.GetString(flagFormat)
	//o.EnableColor = viper.GetBool(flagEnableColor)
	//o.EnableStacktrace = viper.GetBool(flagDisableStacktrace)
	//o.DisableCaller = viper.GetBool(flagDisableCaller)
	//
	//o.OutputPaths = viper.GetStringSlice(flagOutputPaths)
	//o.ErrorOutput = viper.GetStringSlice(flagErrorOutput)
	return nil
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *LogOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.StringVar(&o.LogName, flagName, o.LogName, "The name of file for the logger.")

	fs.BoolVar(&o.EnableStacktrace, flagEnableStacktrace,
		o.EnableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	// 定义多个日志输出文件
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutput, flagErrorOutput, o.ErrorOutput, "Error log output paths of log.")
}

//func (o *LogOption) String() string {
//	data, _ := json.Marshal(o)
//
//	return string(data)
//}

// Build constructs a global zap logger from the Config and LogOption.
func (o *LogOption) Build() *zap.SugaredLogger {

	// 定义日志级别
	logger.SetOptions(
		logger.WithLevel(o.Level),
		logger.WithDisableCaller(o.DisableCaller),
		logger.WithFormat(o.Format),
		logger.WithDisableStacktrace(!o.EnableStacktrace),
		logger.WithEnableColor(o.EnableColor),
		logger.WithAddCallerSkip(1),
	)

	// 定义单独一个日志输出文件
	if len(o.LogName) > 0 {
		outputPaths := []string{o.LogName}
		logger.SetOptions(logger.WithOutputPaths(outputPaths))
	}

	// 根据定义的路径输出日志
	if len(o.OutputPaths) > 0 {
		logger.SetOptions(logger.WithOutputPaths(o.OutputPaths))
	}

	if len(o.ErrorOutput) > 0 {
		logger.SetOptions(logger.WithErrorOutputPaths(o.ErrorOutput))
	}

	// 获取一个logger 和 sugar logger
	return logger.Log.Sugar()
}
