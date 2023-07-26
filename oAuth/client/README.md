# OAuth 客户端样例

本文同步发布在[本人博客](https://l-time.github.io/)，前往博客获得更佳的阅读体验。

OAuth 2.0是一个关于授权的开放网络标准，主要致力于简化客户端人员开发，同时为Web应用程序、桌面应用程序、移动电话和物联网设备提供特定的授权规范。它的官网在[这里](https://oauth.net/2/)。在[RFC6749](https://tools.ietf.org/html/rfc6749)中有明确协议规范。

## 理解OAuth 2.0

在开始编程前，我们先要掌握OAuth 2.0认证的流程是什么样的。

### 相关名词

在理解OAuth 2.0前，我们先需要理解下面这几个名词。

- `resource owner`： 资源所有者。
- `resource server`： 资源服务器，即服务提供商存放用户资源的服务器。
- `client`： 客户端，需要得到资源的应用程序。
- `authorization server`：认证服务器，即服务提供商专门用来处理认证的服务器。
- `user-agent`：用户代理。我们用来访问客户端的程序。

### 协议流程

下图来自[RFC6749](https://tools.ietf.org/html/rfc6749)。

```txt
     +--------+                               +---------------+
     |        |--(A)- Authorization Request ->|   Resource    |
     |        |                               |     Owner     |
     |        |<-(B)-- Authorization Grant ---|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(C)-- Authorization Grant -->| Authorization |
     | Client |                               |     Server    |
     |        |<-(D)----- Access Token -------|               |
     |        |                               +---------------+
     |        |
     |        |                               +---------------+
     |        |--(E)----- Access Token ------>|    Resource   |
     |        |                               |     Server    |
     |        |<-(F)--- Protected Resource ---|               |
     +--------+                               +---------------+

                     Figure 1: Abstract Protocol Flow
```

- (A) 客户端请求`resource owner`授权；

  这种授权可以是直接向`resource owner`请求，也可以通过`authorization server`间接请求。

- (B) 用户同意授权操作；

- (C) 客户端拿着上一步的授权，向`authorization server`申请令牌(access_token)；

- (D) `authorization server`确认授权无误后，发放令牌(access_token)；

- (E) 客户端拿着令牌（access_token）到`resource server`去获取资源;

- (F) `resource server`确认无误，同意向客户端下发受保护的资源。

下面详细说到其中四种授权模式。

### 授权模式

客户端要得到令牌(access_token), 必须需要得到用户的授权，在OAuth 2.0 中定义了四种授权模式：

- 授权码模式 (Authorization Code)
- 简化模式 (Implicit)
- 密码模式 (Resource Owner Password Credentials)
- 客户端模式 (Client Credentials)

每种模式的使用场景与流程都有一定的差别。

#### 授权码模式 (Authorization Code)

授权码模式的授权流程是基于重定向，流程图如下(来自RFC6749)：

```
    +----------+
     | Resource |
     |   Owner  |
     |          |
     +----------+
          ^
          |
         (B)
     +----|-----+          Client Identifier      +---------------+
     |         -+----(A)-- & Redirection URI ---->|               |
     |  User-   |                                 | Authorization |
     |  Agent  -+----(B)-- User authenticates --->|     Server    |
     |          |                                 |               |
     |         -+----(C)-- Authorization Code ---<|               |
     +-|----|---+                                 +---------------+
       |    |                                         ^      v
      (A)  (C)                                        |      |
       |    |                                         |      |
       ^    v                                         |      |
     +---------+                                      |      |
     |         |>---(D)-- Authorization Code ---------'      |
     |  Client |          & Redirection URI                  |
     |         |                                             |
     |         |<---(E)----- Access Token -------------------'
     +---------+       (w/ Optional Refresh Token)

   Note: The lines illustrating steps (A), (B), and (C) are broken into
   two parts as they pass through the user-agent.

                     Figure 3: Authorization Code Flow
```

- (A) 用户访问客户端，客户端将用户重定向到认证服务器；
- (B) 用户选择是否授权；
- (C) 如果用户同意授权，认证服务器重定向到客户端事先指定的地址，而且带上授权码(code)；
- (D) 客户端收到授权码，带着前面的重定向地址，向认证服务器申请访问令牌；
- (E) 认证服务器核对授权码与重定向地址，确认后向客户端发送访问令牌和更新令牌(可选)。

1. 在A中，客户端申请授权，重定向到认证服务器的URI中需要包含这些参数：

| 参数名称      | 参数含义                                                     | 是否必须 |
| ------------- | ------------------------------------------------------------ | -------- |
| response_type | 授权类型，此处的值为`code`                                   | 必须     |
| client_id     | 客户端ID，客户端到资源服务器注册的ID                         | 必须     |
| redirect_uri  | 重定向URI                                                    | 可选     |
| scope         | 申请的权限范围，多个逗号隔开                                 | 可选     |
| state         | 客户端的当前状态，可以指定任意值，认证服务器会原封不动的返回这个值 | 推荐     |

RFC6749中例子如下：

```bash
   GET /authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz
        &redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb HTTP/1.1
    Host: server.example.com
```

2. 在C中，认证服务器返回的URI中，需要包含下面这些参数：

| 参数名称 | 参数含义                                                     | 是否必须 |
| -------- | ------------------------------------------------------------ | -------- |
| code     | 授权码。认证服务器返回的授权码，生命周期不超过10分钟，而且要求只能使用一次，和A中的`client_id`,`redirect_uri`绑定。 | 必须     |
| state    | 如果A中请求包含这个参数，资源服务器原封不动的返回。          | 可选     |

如：

```bash
     HTTP/1.1 302 Found
     Location: https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA
               &state=xyz
```

3. 在D中客户端向认证服务器申请令牌(access_token)时，需要包含下面这些参数。

| 参名称       | 参数含义                               | 是否必须 |
| ------------ | -------------------------------------- | -------- |
| grant_type   | 授权模式，此处为`authorization_code`。 | 必须     |
| code         | 授权码，C中获取的`code`。              | 必须     |
| redirect_uri | 重定向URI，需要和A中一致。             | 必须     |
| client_id    | 客户端ID，与A中一致。                  | 必须     |

如：

```bash
     POST /token HTTP/1.1
     Host: server.example.com
     Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
     Content-Type: application/x-www-form-urlencoded

     grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA
     &redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb
```

4. 在E中，认证服务器返回的信息中，包含下面参数：

| 参数名称      | 参数含义                                                   | 是否必须 |
| ------------- | ---------------------------------------------------------- | -------- |
| access_token  | 访问令牌                                                   | 必须     |
| token_type    | 令牌类型，大小写不敏感。例如 Bearer，MAC。                 | 必须     |
| expires_in    | 过期时间(s)， 如果不设置也要通过其他方法设置一个。         | 推荐     |
| refresh_token | 更新令牌的token。当令牌过期的时候，可用通过该值刷新token。 | 可选     |
| scope         | 权限范围，如果与客户端申请范围一致，可省略。               | 可选     |

如：

```bash
     HTTP/1.1 200 OK
     Content-Type: application/json;charset=UTF-8
     Cache-Control: no-store
     Pragma: no-cache

     {
       "access_token":"2YotnFZFEjr1zCsicMWpAA",
       "token_type":"example",
       "expires_in":3600,
       "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
       "example_parameter":"example_value"
     }
```

5. 如果我们的令牌过期了，需要更新，这里就需要使用`refresh_token`获取一个新令牌了。此时发起HTTP请求需要的参数有：

| 参数名称      | 参数含义                        | 是否必须 |
| ------------- | ------------------------------- | -------- |
| grant_type    | 授权类型，此处是`refresh_token` | 必须     |
| refresh_token | 更新令牌的token。               | 必须     |
| scope         | 权限范围。                      | 可选     |

如：

```bash
     POST /token HTTP/1.1
     Host: server.example.com
     Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
     Content-Type: application/x-www-form-urlencoded

     grant_type=refresh_token&refresh_token=tGzv3JOkF0XG5Qx2TlKWIA
```

这就是授权码模式，应该也是我们平常见的较多的模式了。

#### 简化模式 (Implicit)

简化模式，相当于授权码模式中，C步骤不再通过客户端，直接在浏览器`(user-agent)`中向认证服务器申请令牌，认证服务器不再返回授权码，所有步骤都在浏览器中完成，最后资源服务器将令牌放在`Fragment`中，浏览器从中将令牌提取，发送给客户端。

所以这个令牌对访问者时可见的，而且客户端不需要认证。详细流程如下。

```txt
     +----------+
     | Resource |
     |  Owner   |
     |          |
     +----------+
          ^
          |
         (B)
     +----|-----+          Client Identifier     +---------------+
     |         -+----(A)-- & Redirection URI --->|               |
     |  User-   |                                | Authorization |
     |  Agent  -|----(B)-- User authenticates -->|     Server    |
     |          |                                |               |
     |          |<---(C)--- Redirection URI ----<|               |
     |          |          with Access Token     +---------------+
     |          |            in Fragment
     |          |                                +---------------+
     |          |----(D)--- Redirection URI ---->|   Web-Hosted  |
     |          |          without Fragment      |     Client    |
     |          |                                |    Resource   |
     |     (F)  |<---(E)------- Script ---------<|               |
     |          |                                +---------------+
     +-|--------+
       |    |
      (A)  (G) Access Token
       |    |
       ^    v
     +---------+
     |         |
     |  Client |
     |         |
     +---------+
```

- (A) 客户端将用户导向认证服务器， 携带客户端ID及重定向URI；
- (B) 用户授权；
- (C) 用户同意授权后，认证服务器重定向到A中指定的URI，并且在URI的`Fragment`中包含了访问令牌；
- (D) 浏览器向资源服务器发出请求，该请求中不包含C中的`Fragment`值；
- (E) 资源服务器返回一个网页，其中包含了可以提取C中`Fragment`里面访问令牌的脚本；
- (F) 浏览器执行E中获得的脚本，提取令牌；
- (G) 浏览器将令牌发送给客户端。

1. 在A步骤中，客户端发送请求，需要包含这些参数：

| 参数名称      | 参数含义                                       | 是否必须 |
| ------------- | ---------------------------------------------- | -------- |
| response_type | 授权类型，此处值为`token`。                    | 必须     |
| client_id     | 客户端的ID。                                   | 必须     |
| redirect_uri  | 重定向的URI。                                  | 可选     |
| scope         | 权限范围。                                     | 可选     |
| state         | 客户端的当前状态。指定后服务器会原封不动返回。 | 推荐     |

如：

```bash
    GET /authorize?response_type=token&client_id=s6BhdRkqt3&state=xyz
        &redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb HTTP/1.1
    Host: server.example.com
```

2. 在C中，认证服务器返回的URI中，参数主要有：

| 参数名称     | 参数含义                               | 是否必须 |
| ------------ | -------------------------------------- | -------- |
| access_token | 访问令牌。                             | 必须     |
| token_type   | 令牌类型。                             | 必须     |
| expires_in   | 过期时间。                             | 推荐     |
| scope        | 权限范围。                             | 可选     |
| state        | 客户端访问时如果指定了，原封不动返回。 | 可选     |

如：

```bash
     HTTP/1.1 302 Found
     Location: http://example.com/cb#access_token=2YotnFZFEjr1zCsicMWpAA
               &state=xyz&token_type=example&expires_in=3600
```

我们可以看到C中返回的是一个重定向，而重定向的这个网址的`Fragment`部分包含了令牌。

D步骤中就是访问这个重定向指定的URI，而且不带`Fragment`部分，服务器会返回从`Fragment`中提取令牌的脚本，最后浏览器运行脚本获取到令牌发送给客户端。

#### 密码模式 (Resource Owner Password Credentials)

密码模式就是用户直接将用户名密码提供给客户端，客户端使用这些信息到认证服务器请求授权。具体流程如下：

```
     +----------+
     | Resource |
     |  Owner   |
     |          |
     +----------+
          v
          |    Resource Owner
         (A) Password Credentials
          |
          v
     +---------+                                  +---------------+
     |         |>--(B)---- Resource Owner ------->|               |
     |         |         Password Credentials     | Authorization |
     | Client  |                                  |     Server    |
     |         |<--(C)---- Access Token ---------<|               |
     |         |    (w/ Optional Refresh Token)   |               |
     +---------+                                  +---------------+

            Figure 5: Resource Owner Password Credentials Flow
```

- (A) 资源所有者提供用户名密码给客户端；
- (B) 客户端拿着用户名密码去认证服务器请求令牌；
- (C) 认证服务器确认后，返回令牌；

1. 在B中客户端发送的请求中，需要包含这些参数：

| 参数名称   | 参数含义                       | 是否必须 |
| ---------- | ------------------------------ | -------- |
| grant_type | 授权类型，此处值为`password`。 | 必须     |
| username   | 用户名。                       | 必须     |
| password   | 用户的密码。                   | 必须     |
| scope      | 权限范围。                     | 可选     |

如：

```bash
     POST /token HTTP/1.1
     Host: server.example.com
     Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
     Content-Type: application/x-www-form-urlencoded

     grant_type=password&username=johndoe&password=A3ddj3w
```

2. 在C中，认证服务器返回访问令牌。如：

```bash
     HTTP/1.1 200 OK
     Content-Type: application/json;charset=UTF-8
     Cache-Control: no-store
     Pragma: no-cache

     {
       "access_token":"2YotnFZFEjr1zCsicMWpAA",
       "token_type":"example",
       "expires_in":3600,
       "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
       "example_parameter":"example_value"
     }
```

#### 客户端模式 (Client Credentials)

客户端模式，其实就是客户端直接向认证服务器请求令牌。而用户直接在客户端注册即可，一般用于后端 API 的相关操作。其流程如下：

```
     +---------+                                  +---------------+
     |         |                                  |               |
     |         |>--(A)- Client Authentication --->| Authorization |
     | Client  |                                  |     Server    |
     |         |<--(B)---- Access Token ---------<|               |
     |         |                                  |               |
     +---------+                                  +---------------+

                     Figure 6: Client Credentials Flow
```

- (A) 客户端发起身份认证，请求访问令牌；
- (B) 认证服务器确认无误，返回访问令牌。

1. 在A中，客户端发起请求的参数有：

| 参数名称   | 参数含义                                 | 是否必须 |
| ---------- | ---------------------------------------- | -------- |
| grant_type | 授权类型，此处值为`client_credentials`。 | 必须     |
| scope      | 权限范围。                               | 可选     |

如：

```bash
     POST /token HTTP/1.1
     Host: server.example.com
     Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
     Content-Type: application/x-www-form-urlencoded

     grant_type=client_credentials
```

2. 认证服务器认证后，发放访问令牌，如：

```bash
     HTTP/1.1 200 OK
     Content-Type: application/json;charset=UTF-8
     Cache-Control: no-store
     Pragma: no-cache

     {
       "access_token":"2YotnFZFEjr1zCsicMWpAA",
       "token_type":"example",
       "expires_in":3600,
       "example_parameter":"example_value"
     }
```

## 实现一个OAuth 2.0的客户端

说完协议内容后我们来着手实现一个客户端，假定现在我们有着这样一套流程：

> 应用登录——>使用GitHub授权——>获取用户信息

那么我们使用第一种方式来实现一个基本流程。

### 登录页

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
    <a href="https://github.com/login/oauth/authorize?client_id=<id>&redirect_uri=<callback>">
    使用GitHub登录？
    </a>
</body>
</html>
```

在这个简单的登录页中我们加入了一个链接，用于访问GitHub的认证服务器。

其中链接有两个参数：

- `client_id`：应用ID，这和下文的申请有关
- `redirect_uri`：回调地址，这和下文的申请有关

### 申请OAuth应用

打开[申请网站](https://github.com/settings/applications/new)，填写相关信息。

![image-20230726213100533](https://s2.loli.net/2023/07/26/G3jcTgeAlw7Ob8d.png)

1. Application name：应用名
2. Homepage URL：应用的URL
3. Authorization callback URL：回调地址。这和上文的回调地址要保持一致。

为了后文便于理解，回调地址填写为：`http://localhost:8080/oauth/redirect`

申请后你会得到两个参数：`Client ID`和`Client secrets`，其中前者要用于网页链接地址，而后者将用在后续流程中。

### 简单的服务

作者在此处使用到了`Gin`框架，感兴趣的读者可以自行探索，也可以等笔者研究一番。

创建如下的文件结构：

```text
├─go.mod
├─go.sum
├─main.go
├─public
|   └index.html


```

在`main.go`内创建一下简单的路由：

```go
package main

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	file, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile(file+"/public", false)))
	r.Run(":8080")
}
```

运行启动，此时我们便已经可以看到我们的应用了：

![image-20230726214026377](https://s2.loli.net/2023/07/26/e8KIa6pBVLWTtUw.png)

点击链接，我们便可以访问到授权页面：

![image-20230726214223620](https://s2.loli.net/2023/07/26/IN6p9AqgLVlxEFj.png)

但是别急着授权，我们其他的业务逻辑都还没写完呢，下一步就是去实现后续的授权流程了。

当我们授权后，GitHub会将我们重定向至应用内的回调链接，同时会带上一个`code`参数。例如：

```text
http://localhost:8080/oauth/redirect?code=1233211234567
```

回想一下授权码模式：

> (D) 客户端收到授权码，带着前面的重定向地址，向认证服务器申请访问令牌；

这里的`code`便是我们收到的授权码，我们下一步就是要申请访问令牌。

### 重定向路由

我们按照授权码模式，在这个路由需要实现：

1. 解析授权码
2. 访问认证服务器
3. 拿到访问令牌
4. 定位到我们授权成功页面，并附上访问令牌执行后续操作

那我们便补全一下逻辑：

```go
// HandleOAuth 便是我们的重定向路由
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

	// 第四步：从返回信息读取拿到的令牌
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		log.Printf("could not parse JSON response: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}
	// 最后一步：重定向至我们授权成功后的页面
	c.Redirect(http.StatusFound, "/welcome.html?access_token="+t.AccessToken)
}

```

### 授权成功的欢迎页面

在`public`下新建一个`welcome.html`，添加以下内容：

```html
<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>你好</title>
</head>

<body>

</body>
<script>
	// 获取访问令牌
	const query = window.location.search.substring(1)
	const token = query.split('access_token=')[1]

	// 访问资源服务器地址，获取相关资源
	fetch('https://api.github.com/user', {
			headers: {
                // 将token放在Header中
				Authorization: 'token ' + token
			}
		})
		// 解析返回的JSON
		.then(res => res.json())
		.then(res => {
            // 这里我们能得到很多信息
			// 具体看这里 https://developer.github.com/v3/users/#get-the-authenticated-user
			// 这里我们就只展示一下用户名了
			const nameNode = document.createTextNode(`Welcome, ${res.name}`)
			document.body.appendChild(nameNode)
		})
</script>

```

同时完善一下`main.go`：

```go
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
	clientID     = "d27f74811df1d22fb44e"
	clientSecret = "80d9a7a285c5d1c21e80e421bc0b31508a9bb6a7"
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

	r := gin.Default()
    //由于使用了静态资源，这里是一定要这么写的，若你使用：
    //r.StaticFS("/", http.Dir(file+"\\public"))
    //会引起恐慌，具体的原因可以看这篇回答：
    //https://stackoverflow.com/questions/36357791/gin-router-path-segment-conflicts-with-existing-wildcard
	r.Use(static.Serve("/", static.LocalFile(file+"/public", false)))

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

```

让我们重新让他运行起来：

1. 访问`http://localhost:8080/`，进入登录页面：

![image-20230726214026377](https://s2.loli.net/2023/07/26/h9pjItlfoE5ua7z.png)

2. 进入授权页面：

   ![image-20230726214223620](https://s2.loli.net/2023/07/26/IN6p9AqgLVlxEFj.png)

我们点击一下授权，进入重定向页面：

![image-20230726215933032](https://s2.loli.net/2023/07/26/GKw7ul62OJ1dfIa.png)

3. 授权后重定向至欢迎页面：

   ![image-20230726220006847](https://s2.loli.net/2023/07/26/chOCpZaJS7AzPtu.png)

从这里，基础的一个OAuth客户端和流程就结束了。

## 后话

实际上我们可以通过`OAuth`应用获得很多信息，具体可以看GitHub的[官方文档](https://docs.github.com/zh/rest)

这里我们直接通过URL的方式传递令牌实际上是不安全的，我们可以通过`Cookie`来传递令牌

下篇文章我们将尝试实现一个OAuth服务端，并重现授权码模式的流程。




## 参考

- [FRC6749](https://tools.ietf.org/html/rfc6749)
- [Implementing OAuth 2.0 with Go(Golang)](https://www.sohamkamani.com/blog/golang/2018-06-24-oauth-with-golang/)
- [10 分钟理解什么是 OAuth 2.0 协议](https://deepzz.com/post/what-is-oauth2-protocol.html)
- [理解OAuth 2.0](http://www.ruanyifeng.com/blog/2014/05/oauth_2_0.html)