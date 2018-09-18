package brot

import (
	"bytes"
	cr "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fuxsig/brot/di"
	"github.com/fuxsig/brot/wrapper"
	"github.com/gorilla/sessions"
	verifier "github.com/okta/okta-jwt-verifier-golang"
)

var sessionStore = sessions.NewCookieStore([]byte("okta-hosted-login-session-store"))

type OktaWrapper struct {
	ClientID     string `brot:"client_id,mandatory"`
	ClientSecret string `brot:"client_secret,mandatory"`
	RedirectURI  string `brot:"redirect_uri,mandatory"`
	Issuer       string `brot:"issuer,mandatory"`
}

type Exchange struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
	IdToken          string `json:"id_token,omitempty"`
}

func (h *OktaWrapper) exchangeCode(code string, r *http.Request) Exchange {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(h.ClientID + ":" + h.ClientSecret))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("redirect_uri", h.RedirectURI)

	url := h.Issuer + "/v1/token?" + q.Encode()

	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
	header := req.Header
	header.Add("Authorization", "Basic "+authHeader)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	header.Add("Connection", "close")
	header.Add("Content-Length", "0")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var exchange Exchange
	json.Unmarshal(body, &exchange)

	return exchange

}

func (h *OktaWrapper) verifyToken(nonce, t string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = h.ClientID
	jv := verifier.JwtVerifier{
		Issuer:           h.Issuer,
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified")
}

func (h *OktaWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	if state != "ApplicationState" {
		log.Printf("OktaWrapper: Wrong state %s in request %s\n", state, r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := q.Get("code")
	if code == "" {
		log.Printf("OktaWrapper: Missing code in request %s\n", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exchange := h.exchangeCode(code, r)
	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")
	if err != nil {
		log.Printf("OktaWrapper: Session error: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nonce, ok := session.Values["nonce"].(string)
	if !ok {
		log.Printf("OktaWrapper: Missing nonce value in session\n")
	}
	url, ok := session.Values["url"].(string)
	if !ok {
		url = "/"
	}

	if _, err = h.verifyToken(nonce, exchange.IdToken); err != nil {
		log.Printf("OktaWrapper: Verify token error: %s\n", err.Error())
	} else {
		session.Values["id_token"] = exchange.IdToken
		session.Values["access_token"] = exchange.AccessToken

		session.Save(r, w)
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func isAuthenticated(r *http.Request) bool {
	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")

	if err != nil || session.Values["id_token"] == nil || session.Values["id_token"] == "" {
		return false
	}

	return true
}

func generateNonce() string {
	nonceBytes := make([]byte, 32)
	cr.Read(nonceBytes)
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

func (h *OktaWrapper) ServeChain(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if isAuthenticated(r) {
		next.ServeHTTP(w, r)
		return
	}
	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")
	if err != nil {
		log.Printf("OktaWrapper: Session error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	nonce := generateNonce()
	session.Values["nonce"] = nonce
	session.Values["url"] = r.URL.String()
	err = session.Save(r, w)
	if err != nil {
		log.Printf("OktaWrapper: Session error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	q := r.URL.Query()
	q.Add("client_id", h.ClientID)
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", h.RedirectURI)
	q.Add("state", "ApplicationState")
	q.Add("nonce", nonce)
	redirect := h.Issuer + "/v1/authorize?" + q.Encode()
	http.Redirect(w, r, redirect, http.StatusMovedPermanently)
}

func (h *OktaWrapper) HandlerFunc() http.Handler {
	return h
}

var _ wrapper.Handler = (*OktaWrapper)(nil)
var _ ProvidesHandler = (*OktaWrapper)(nil)
var _ = di.GlobalScope.Declare((*OktaWrapper)(nil))
