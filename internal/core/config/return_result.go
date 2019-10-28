package config

import (
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var RR = &ReturnResult{}

type Result struct {
	Code    int    `json:"code" mapstructure:"code"`
	Message string `json:"message" mapstructure:"message"`
}

type ReturnResult struct {
	InvalidParam    Result `mapstructure:"invalid_param"`
	InvalidRequest  Result `mapstructure:"invalid_request"`
	InvalidFormat   Result `mapstructure:"invalid_format"`
	InvalidDownload Result `mapstructure:"invalid_download"`
	InvalidUpload   Result `mapstructure:"invalid_upload"`
	FileNotFound    Result `mapstructure:"file_not_found"`
	Internal        struct {
		Success           Result `mapstructure:"success"`
		General           Result `mapstructure:"general"`
		BadRequest        Result `mapstructure:"bad_request"`
		ConnectionError   Result `mapstructure:"connection_error"`
		DBSessionNotFound Result `mapstructure:"db_session_not_found"`
		Unauthorized      Result `mapstructure:"unauthorized"`
	} `mapstructure:"internal"`
}

func (ec Result) Error() string {
	return ec.Message
}

func (ec Result) ErrorCode() int {
	return ec.Code
}

func (ec Result) HTTPStatusCode() int {
	switch {
	case ec.Code == 0000: // success
		return http.StatusOK
	case ec.Code == 400: // bad request
		return http.StatusBadRequest
	case ec.Code == 404: // connection_error
		return http.StatusNotFound
	case ec.Code == 401: // unauthorized
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}

func InitReturnResult(configPath string) error {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("return_result")

	if err := v.ReadInConfig(); err != nil {
		logrus.Error("read config file error:", err)
		return err
	}

	if err := bindingReturnResult(v, RR); err != nil {
		logrus.Error("binding config error:", err)
		return err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logrus.Info("config file changed:", e.Name)
		if err := bindingReturnResult(v, RR); err != nil {
			logrus.Error("binding error:", err)
		}
		logrus.Infof("Initial 'Return Result'. %+v", RR)
	})
	return nil
}

func bindingReturnResult(vp *viper.Viper, rr *ReturnResult) error {
	if err := vp.Unmarshal(&rr); err != nil {
		logrus.Error("unmarshal config error:", err)
		return err
	}
	return nil
}

// NewResult for new result
func NewResult(code int, message string) *Result {
	return &Result{
		Code:    code,
		Message: message,
	}
}
