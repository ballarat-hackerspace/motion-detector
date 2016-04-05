package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func chk(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/motion-detector/")
	viper.AddConfigPath("$HOME/.motion-detector/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	chk(err)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/enable", func(c *gin.Context) {
		apiUrl := viper.GetString("spark_endpoint")
		data := url.Values{}
		data.Set("access_token", viper.GetString("spark_api"))
		data.Add("args", "pir-enable")

		u, _ := url.ParseRequestURI(apiUrl)
		u.Path = "/v1/devices/" + viper.GetString("spark_device") + "/action"
		urlStr := fmt.Sprintf("%v", u)

		client := &http.Client{}
		req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		resp, _ := client.Do(req)
		c.JSON(200, gin.H{
			"message": resp.Status,
			"urlstr":  urlStr,
		})
	})
	r.GET("/disable", func(c *gin.Context) {
		disable_length, err := strconv.ParseInt(c.Query("for"), 10, 64)
		if err != nil {
			disable_length = 0
		}

		// backwards compatibility with old hour based usage
		if disable_length < 6 {
			disable_length = disable_length * 3600
		}
		apiUrl := viper.GetString("spark_endpoint")
		data := url.Values{}
		data.Set("access_token", viper.GetString("spark_api"))
		if disable_length != 0 {
			data.Add("args", "pir-disable," + fmt.Sprintf("%d", disable_length))
		} else {
			data.Add("args", "pir-disable")
		}

		u, _ := url.ParseRequestURI(apiUrl)
		u.Path = "/v1/devices/" + viper.GetString("spark_device") + "/action"
		urlStr := fmt.Sprintf("%v", u)

		client := &http.Client{}
		req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
		resp, _ := client.Do(req)
		c.JSON(200, gin.H{
			"message": resp.Status,
			"urlstr":  urlStr,
		})
	})
	r.Run("127.0.0.1:8081")
}
