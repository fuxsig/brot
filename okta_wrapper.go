package brot

import (
	cr "crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

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

func (h *OktaWrapper) exchangeCode(code string, r *http.Request) (exchange *Exchange, err error) {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(h.ClientID + ":" + h.ClientSecret))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("redirect_uri", h.RedirectURI)

	url := h.Issuer + "/v1/token?" + q.Encode()

	var (
		req  *http.Request
		resp *http.Response
		body []byte
	)

	if req, err = http.NewRequest("POST", url, nil); err != nil {
		log.Printf("[OktaWrapper.exchangeCode] request error: %s", err.Error())
		return
	}

	header := req.Header
	header.Add("Authorization", "Basic "+authHeader)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	header.Add("Connection", "close")
	header.Add("Content-Length", "0")

	// todo: this needs to be removed
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	if resp, err = client.Do(req); err != nil {
		log.Printf("[OktaWrapper.exchangeCode] http error: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Printf("[OktaWrapper.exchangeCode] read body error: %s", err.Error())
		return
	}
	exchange = new(Exchange)
	err = json.Unmarshal(body, exchange)
	if err != nil {
		log.Printf("[OktaWrapper.exchangeCode] json unmarshall error: %s", err.Error())
		return nil, nil
	}
	return
}

func (h *OktaWrapper) verifyToken(nonce, t string) *verifier.Jwt {
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = h.ClientID
	jv := verifier.JwtVerifier{
		Issuer:           h.Issuer,
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)
	if err != nil {
		log.Printf("[OktaWrapper.verifyToken] verify error: %s", err.Error())
		return nil
	}
	return result
}

func (h *OktaWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	if state != "ApplicationState" {
		log.Printf("[OktaWrapper.ServeHTTP]: Wrong state %s in request %s\n", state, r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := q.Get("code")
	if code == "" {
		log.Printf("[OktaWrapper.ServeHTTP]: Missing code in request %s\n", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exchange, err := h.exchangeCode(code, r)
	if err != nil {
		log.Printf("[OktaWrapper.ServeHTTP]: Exchange error: %s\n", err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}

	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")
	if err != nil {
		log.Printf("[OktaWrapper.ServeHTTP]: Session error: %s\n", err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	/*nonce, ok := session.Values["nonce"].(string)
	if !ok {
		log.Printf("[OktaWrapper.ServeHTTP]: Missing nonce value in session\n")
		w.WriteHeader(http.StatusForbidden)
		return
	}*/
	url, ok := session.Values["url"].(string)
	if !ok {
		url = "/"
	}

	// the external lib does not work well at the moment :(
	/*if _, err = h.verifyToken(nonce, exchange.IdToken); err != nil {
		log.Printf("[OktaWrapper.ServeHTTP]: Verify token error: %s\n", err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	} else {*/
	session.Values["id_token"] = exchange.IdToken
	session.Values["access_token"] = exchange.AccessToken

	session.Save(r, w)
	//}

	// retrieve user information
	var (
		req  *http.Request
		resp *http.Response
	)

	reqURL := h.Issuer + "/v1/userinfo"
	req, err = http.NewRequest("GET", reqURL, nil)
	head := req.Header
	head.Add("Authorization", "Bearer "+exchange.AccessToken)
	head.Add("Accept", "application/json")

	// todo: this needs to be removed
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// this is the trick
	defer http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	if resp, err = client.Do(req); err != nil {
		log.Printf("[OktaWrapper.ServeHTTP] user info request error: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Printf("[OktaWrapper.ServeHTTP] body read error: %s", err.Error())
		return
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal(body, &m); err != nil {
		log.Printf("[OktaWrapper.ServeHTTP] json unmarshall error: %s", err.Error())
		return
	}
	session.Values["email"] = m["email"]
	session.Values["given_name"] = m["given_name"]
	if err = session.Save(r, w); err != nil {
		log.Printf("[OktaWrapper.ServeHTTP] could not save seesion: %s", err.Error())
		return
	}
	log.Printf("[OktaWrapper.ServeHTTP] retrieved user profile for %s", m["email"])
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
		log.Printf("[OktaWrapper.ServeChain]: Session error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	nonce := generateNonce()
	session.Values["nonce"] = nonce
	session.Values["url"] = r.URL.String()
	err = session.Save(r, w)
	if err != nil {
		log.Printf("[OktaWrapper.ServeChain]: Session error: %s", err.Error())
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
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

func (h *OktaWrapper) HandlerFunc() http.Handler {
	return h
}

type OktaLogout struct {
	URL    string `brot:"url,mandatory"`
	Issuer string `brot:"issuer,mandatory"`
}

func (o *OktaLogout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := url.Values{"post_logout_redirect_uri": {o.URL}}
	reqURL := o.Issuer + "/v1/logout"
	session, err := sessionStore.Get(r, "okta-hosted-login-session-store")
	if err != nil {
		reqURL := reqURL + "?" + values.Encode()
		http.Redirect(w, r, reqURL, http.StatusTemporaryRedirect)
		return
	}
	idToken, ok := session.Values["id_token"].(string)
	if ok {
		values.Add("id_token_hint", idToken)
		reqURL = reqURL + "?" + values.Encode()
	}
	defer http.Redirect(w, r, reqURL, http.StatusTemporaryRedirect)

	delete(session.Values, "id_token")
	delete(session.Values, "access_token")
	delete(session.Values, "nonce")
	if err = session.Save(r, w); err != nil {
		log.Printf("[OktaLogout.ServeHTTP] could not save session: %s", err.Error())
		return
	}
	log.Printf("[OktaLogout.ServeHTTP] successfully deleted session for id_token %s", idToken)
}

func (o *OktaLogout) HandlerFunc() http.Handler {
	return o
}

var _ wrapper.Handler = (*OktaWrapper)(nil)
var _ ProvidesHandler = (*OktaWrapper)(nil)
var _ = di.GlobalScope.Declare((*OktaWrapper)(nil))

var _ ProvidesHandler = (*OktaLogout)(nil)
var _ = di.GlobalScope.Declare((*OktaLogout)(nil))
