package auth

import (
	"ops/pkg/ctl"
	"net/http"
	"ops/pkg/cfg"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	JwtSecret string
	jwtSecret = []byte(JwtSecret)
)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
func GetState(c *gin.Context) *cfg.State {
	var (
		state *cfg.State
		si    interface{}
		ok    bool
	)
	si, ok = c.Get("state")
	ctl.Debug(si)
	ctl.Debug(ok)
	if ok {
		state = si.(*cfg.State)
	} else {
		c.Data(200, "", []byte("State error"))
	}
	return state
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctl.Debug(c.Request.URL)
		var (
			token  string
			err    error
			data   interface{}
			claims *Claims
			code   int
		)
		code = SUCCESS
		token, err = c.Cookie("token")
		if err != nil {
			code = ERROR_AUTH_CHECK_TOKEN_FAIL
		}

		if token == "" {
			code = INVALID_PARAMS
		} else {
			claims, err = ParseToken(token)
			if err != nil {
				code = ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}
		if code != SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  GetMsg(code),
				"data": data,
			})
			c.Abort()
			return
		}
		c.Set("user", claims.Username)

		url := c.Request.URL.String()
		if len(url) > 10 {
			if url[:9] == "/v1/game/" {
				ctl.Debug(url[9:])
				gameType := strings.Split(url[9:], "/")[0] //获取URL中的game_type
				if _, ok := cfg.PdtState[gameType]; ok {
					c.Set("state", cfg.PdtState[gameType])
				} else {
					ctl.Debug(url)
					return
				}
			}
			if url[:8] == "/v1/plt/" {
				ctl.Debug(url[8:])
				gameType := strings.Split(url[8:], "/")[0] //获取URL中的game_type
				if _, ok := cfg.PdtState[gameType]; ok {
					c.Set("state", cfg.PdtState[gameType])
				} else {
					ctl.Debug(url)
					return
				}
			}
		}
		if len(url) > 15 {
			if url[:14] == "/v1/pop2/game/" {
				ctl.Debug(url[14:])
				gameType := strings.Split(url[14:], "/")[0]
				if _, ok := cfg.PdtState[gameType]; ok {
					//	return pdtState[gameType]
					c.Set("state", cfg.PdtState[gameType])
				} else {
					ctl.Debug(url)
					return
				}
			}
			if url[:13] == "/v1/pop/game/" {
				ctl.Debug(url[13:])
				gameType := strings.Split(url[13:], "/")[0]
				if _, ok := cfg.PdtState[gameType]; ok {
					//	return pdtState[gameType]
					c.Set("state", cfg.PdtState[gameType])
				} else {
					ctl.Debug(url)
					return
				}
			}
			if url[:13] == "/v1/pop2/plt/" {
				ctl.Debug(url[13:])
				gameType := strings.Split(url[13:], "/")[0]
				if _, ok := cfg.PdtState[gameType]; ok {
					//	return pdtState[gameType]
					c.Set("state", cfg.PdtState[gameType])
				} else {
					ctl.Debug(url)
					return
				}
			}
			if url[:13] == "/v1/pop3/plt/" {
				ctl.Debug(url[13:])
				gameType := strings.Split(url[13:], "/")[0]
				if _, ok := cfg.PdtState[gameType]; ok {
					//	return pdtState[gameType]
					c.Set("state", cfg.PdtState[gameType])
				} else {
					ctl.Debug(url)
					return
				}
			}
			if url[:13] == "/v1/pop4/plt/" {
				ctl.Debug(url[13:])
				gameType := strings.Split(url[13:], "/")[0]
				if _, ok := cfg.PdtState[gameType]; ok {
					//	return pdtState[gameType]
					c.Set("state", cfg.PdtState[gameType])
				} else {
					ctl.Debug(url)
					return
				}
			}

		}
		if len(url) > 13 {
			if url[:12] == "/v1/pop/plt/" {
				ctl.Debug(url[12:])
				gameType := strings.Split(url[12:], "/")[0]
				if _, ok := cfg.PdtState[gameType]; ok {
					//	return pdtState[gameType]
					c.Set("state", cfg.PdtState[gameType])
				} else {
					ctl.Debug(url)
					return
				}
			}

		}

		ctl.Debug(url)

		c.Next()
	}
}

type Game struct {
	Id         int
	GameId     int `xorm:"unique"`
	TypeId     int
	Name       string `xorm:"unique"`
	NameImg    string
	Sort       int
	Status     int
	ZipId      int
	UpdateTime string `xorm:"DATETIME"`
	CreateTime string `xorm:"DATETIME"`
	LastZipId  int
}
