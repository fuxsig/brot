// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fuxsig/brot/di"
)

type JWT struct {
	Path     string                 `brot:"path"`
	Method   string                 `brot:"method"`
	Iss      string                 `brot:"iss"`
	Sub      string                 `brot:"sub"`
	Aud      string                 `brot:"aud"`
	Custom   map[string]interface{} `brot:"custom"`
	ID       string                 `brot:"clientId"`
	Secret   string                 `brot:"clientSecret"`
	URL      string                 `brot:"url"`
	Database ProvidesDatabase       `brot:"database"`
	pk       *rsa.PrivateKey
	sm       jwt.SigningMethod
	token    *JWTToken
}

type JWTToken struct {
	Expires time.Time
	Token   string
}

type AECAccessToken struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	Expires   int64  `json:"expires_in"`
}

func (jt *JWTToken) Valid() bool {
	return time.Now().Add(time.Minute * 10).Before(jt.Expires)
}

func (j *JWT) InitFunc() (err error) {
	pemKey, err := ioutil.ReadFile(j.Path)
	if err != nil {
		return
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemKey)
	if err != nil {
		return
	}
	j.pk = privateKey

	method := jwt.GetSigningMethod(j.Method)
	if method == nil {
		log.Fatalf("Signing method %s is unknown", j.Method)
	}
	j.sm = method
	return
}

func (j *JWT) Retry() bool {
	return false
}

func (j *JWT) Token() string {
	if j.token != nil && j.token.Valid() {
		return j.token.Token
	}

	t := &JWTToken{}
	if err := j.Database.Get("JWTToken4AEC", t); err == nil {
		if t.Valid() {
			j.token = t
			fmt.Println("Authorization: ", "<"+t.Token+">")
			return t.Token
		}
	} else {
		if Warning {
			Warning.Printf("JWT.Token: could not load token from database: %s", err.Error())
		}
	}

	if token := jwt.New(j.sm); token != nil {
		claims := token.Claims.(jwt.MapClaims)
		claims["iss"] = j.Iss
		claims["sub"] = j.Sub
		claims["aud"] = j.Aud
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		for key, value := range j.Custom {
			claims[key] = value
		}

		tokenString, err := token.SignedString(j.pk)
		if err != nil {
			log.Fatal(err.Error())
		}

		values := url.Values{
			"client_id":     {j.ID},
			"client_secret": {j.Secret},
			"jwt_token":     {tokenString}}
		urlStr := j.URL
		res, err := http.PostForm(urlStr, values)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		switch {
		case res.StatusCode < 200:
		case res.StatusCode >= 300:
			log.Fatalf("[%d %s]%s", res.StatusCode, res.Status, string(b))
		}
		if err != nil {
			log.Fatal(err.Error())
		}

		at := AECAccessToken{}
		err = json.Unmarshal(b, &at)
		if err != nil {
			log.Fatal(err.Error())
		}
		t.Token = fmt.Sprint(strings.Title(at.TokenType), " ", at.Token)
		t.Expires = time.Now().Add(time.Millisecond * time.Duration(at.Expires))
		fmt.Println("Authorization: ", "<"+t.Token+">", t.Expires.String())
		j.token = t

		j.Database.Set("JWTToken4AEC", t)
		return t.Token

	}
	return ""
}

var _ di.ProvidesInit = (*JWT)(nil)
var _ = di.GlobalScope.Declare((*JWT)(nil))
