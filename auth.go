// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"crypto/subtle"
	"encoding/base64"
	"strconv"
)

// AuthUserKey is the cookie name for user credential in basic auth.
const AuthUserKey = "user"

//获取map
// Accounts defines a key/value for user/pass list of authorized logins.
type Accounts map[string]string

//定义基本的authoPair
type authPair struct {
	value string
	user  string
}

//pairs列表
//pairs
type authPairs []authPair

//查找authValue是否在authPairs里面。否则返回user
func (a authPairs) searchCredential(authValue string) (string, bool) {
	if authValue == "" {
		return "", false
	}
	for _, pair := range a {
		if pair.value == authValue {
			return pair.user, true
		}
	}
	return "", false
}

// BasicAuthForRealm returns a Basic HTTP Authorization middleware. It takes as arguments a map[string]string where
// the key is the user name and the value is the password, as well as the name of the Realm.
// If the realm is empty, "Authorization Required" will be used by default.
// (see http://tools.ietf.org/html/rfc2617#section-1.2)
//
//基础安全认证中间件？
func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc {
	//处理Auth
	if realm == "" {
		//设置默认的realm
		realm = "Authorization Required"
	}
	//得到realm
	realm = "Basic realm=" + strconv.Quote(realm)

	//把accounts处理成authPairs

	pairs := processAccounts(accounts)
	return func(c *Context) {
		// Search user in the slice of allowed credentials
		//获取author码
		user, found := pairs.searchCredential(c.requestHeader("Authorization"))
		//未通过
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(401)
			return
		}
		//设置基本的author
		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(AuthUserKey, user)
	}
}

//基本的认证
// BasicAuth returns a Basic HTTP Authorization middleware. It takes as argument a map[string]string where
// the key is the user name and the value is the password.
//基本的认证
func BasicAuth(accounts Accounts) HandlerFunc {
	return BasicAuthForRealm(accounts, "")
}

//处理accounts 把它转成pairs
func processAccounts(accounts Accounts) authPairs {
	//处理账户信息
	assert1(len(accounts) > 0, "Empty list of authorized credentials")
	pairs := make(authPairs, 0, len(accounts))
	for user, password := range accounts {
		assert1(user != "", "User can not be empty")
		//对数据进行加密
		value := authorizationHeader(user, password)
		pairs = append(pairs, authPair{
			value: value, //密码
			user:  user,  //用户名
		})
	}
	return pairs
}

//对value进行加密
func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}

//比较str，比较长度和内容
func secureCompare(given, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	}
	// Securely compare actual to itself to keep constant time, but always return false.
	return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
}
