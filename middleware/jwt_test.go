package middleware

import (
	"encoding/json"
	"testing"
)

func TestJwt(t *testing.T) {
	jwt := InitJwt("SecretKey")
	token := jwt.GenerateToken("qiqi") // 生成有效期为24小时的 JWT
	t.Log(token)

	claims := jwt.ParseToken(token)
	b, _ := json.Marshal(claims)
	t.Log(string(b))

	newToken := jwt.CreateTokenByOldToken(token)
	t.Log(newToken)
	newClaims := jwt.ParseToken(newToken)
	b, _ = json.Marshal(newClaims)
	t.Log(string(b))

}
