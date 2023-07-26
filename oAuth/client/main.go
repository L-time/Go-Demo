package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// 你在注册时得到的
const (
	clientID     = "你的ID"
	clientSecret = "你的Secret"
)

var httpClient = http.Client{}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func main() {
	file, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//fmt.Println(file)

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile(file+"/public", false)))
	//r.StaticFS("/", http.Dir(file+"\\public"))
	group := r.Group("/oauth")

	group.GET("/redirect", HandleOAuth)

	r.Run(":8080")

}

func HandleOAuth(c *gin.Context) {
	// 第一步：从此处拿取授权码
	// 只有拿到授权码才可以进行后续操作
	code := c.Query("code")
	// 第二步：生成重定向网址
	// 这里我们需要使用在申请OAuth应用时的clientID、clientSecret来生成
	// code便是我们的授权码
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		clientID, clientSecret, code)
	// 第三步，发送授权请求
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		log.Printf("could not create HTTP request: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	// 设置返回数据的格式为json
	req.Header.Set("accept", "application/json")

	// 这里发送出去
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("could not send HTTP request: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	// 养成关闭的好习惯
	defer res.Body.Close()

	// 从这读取拿到的令牌
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		log.Printf("could not parse JSON response: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	// 最后一步：重定向至我们授权成功后的页面
	c.Redirect(http.StatusFound, "/welcome.html?access_token="+t.AccessToken)
}
