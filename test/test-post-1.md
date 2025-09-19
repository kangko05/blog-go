# Go 언어로 REST API 만들기

Go 언어의 Gin 프레임워크를 사용하여 간단한 REST API를 만드는 방법을 알아보겠습니다.

## 프로젝트 설정

먼저 새로운 Go 모듈을 생성합니다:

```bash
go mod init api-server
go get github.com/gin-gonic/gin
```

## 기본 서버 구조

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run()
}
```

## 미들웨어 활용

Gin은 다양한 미들웨어를 제공합니다. CORS, 로깅, 인증 등을 쉽게 추가할 수 있습니다.

이런 식으로 단계적으로 API 서버를 구축해나가면 됩니다.
