package handler

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	translate "cloud.google.com/go/translate/apiv3"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"

	"github.com/tuckKome/konjac/db"
)

type LangSet struct {
	Language string
	Tag      string `gorm:"type:varchar(10);unique_index"`
	Kotae    string
}

type SessionInfo struct {
	UserID         int
	Name           string //ログインしているユーザの表示名
	IsSessionAlive bool   //セッションが生きているかどうか
}

type MyHis struct {
	Name  string
	Email string
	Word  string
	Lngs  string
}

type pair struct {
	Language string
	code     string
	Response string
}

type translationSet struct {
	text      string
	Languages []string
	codes     []string
	Responses []string
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	log.Println("セッション取得")
	session.Clear()
	log.Println("クリア処理")
	session.Save()
	c.Redirect(302, "/")
}

func Login(c *gin.Context, u db.User) {
	session := sessions.Default(c)
	session.Set("Name", u.Name)
	session.Set("Email", u.Email)
	session.Save()
}

func GetSessionInfo(c *gin.Context) SessionInfo {
	var info SessionInfo
	session := sessions.Default(c)
	userID := session.Get("UserID")
	userName := session.Get("Name")
	if userID == nil {
		info = SessionInfo{
			UserID: -1, Name: "", IsSessionAlive: false,
		}
	} else {
		info = SessionInfo{
			UserID:         userID.(int),
			Name:           userName.(string),
			IsSessionAlive: true,
		}
		log.Println(info)
	}
	return info
}

func Replace(s string) string {
	s = strings.Replace(s, "&lt;", "<", -1)
	s = strings.Replace(s, "&gt;", "<", -1)
	return s
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Index(c *gin.Context) {
	info := GetSessionInfo(c)

	c.HTML(200, "index.html", gin.H{
		"SessionInfo": info,
	})
}

func New(c *gin.Context) {
	text := c.PostForm("word")
	uri := fmt.Sprintf("/translate/%s", text)
	c.Redirect(302, uri)
}

func GetTranslation(c *gin.Context) {
	pairs := []pair{
		{"日本語", "ja", ""},
		{"英語", "en", ""},
		{"ドイツ語", "de", ""},
		{"アイルランド語", "ga", ""},
		{"アラビア語", "ar", ""},
		{"ギリシャ語", "el", ""},
		{"エスペラント", "eo", ""},
		{"スペイン語", "es", ""},
		{"フランス語", "fr", ""},
		{"イタリア語", "it", ""},
		{"オランダ語", "nl", ""},
		{"アフリカーンス語", "af", ""},
		{"フィンランド語", "fi", ""},
		{"スウェーデン語", "sv", ""},
		{"ノルウェー語", "no", ""},
		{"アイスランド語", "is", ""},
		{"ロシア語", "ru", ""},
		{"ポーランド語", "pl", ""},
		{"ブルガリア語", "bg", ""},
		{"ウクライナ語", "uk", ""},
		{"キルギス語", "ky", ""},
		{"トルコ語", "tr", ""},
		{"スワヒリ語", "sw", ""},
	}

	languages := make([]string, len(pairs))
	for i := range pairs {
		languages[i] = pairs[i].language
	}
	codes := make([]string, len(pairs))
	for i := range pairs {
		codes[i] = pairs[i].code
	}
	responses := make([]string, len(pairs))
	transSet := translationSet{Languages: languages, codes: codes, Responses: responses}

	text := c.Param("text")
	transSet.text = text
	texts := []string{text}
	fmt.Println(transSet)
	// Clientをつくる
	client, err := translate.NewTranslationClient(c)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	for i := range pairs {

		request := &translatepb.TranslateTextRequest{
			Contents:           texts,
			TargetLanguageCode: pairs[i].code,
			Parent:             "projects/honyac-konjac",
		}

		translation, err := client.TranslateText(c, request)
		if err != nil {
			log.Fatalf("Failed to translate text: %v", err)
		}
		pairs[i].Response = translation.Translations[0].GetTranslatedText()
	}
	fmt.Println(pairs) //デバッグ用
	// 履歴に残すために
	info := GetSessionInfo(c)
	if info.IsSessionAlive == true {
		var history db.History
		history.Word = text
		history.UserID = info.UserID
		now := time.Now()
		history.CreatedAt = now
		history.UpdatedAt = now
		log.Println(history)
		db.CreateHistory(history)
	}

	c.HTML(200, "translations.html", gin.H{
		"text":        text,
		"Translation": pairs,
		"SessionInfo": info,
	})
}

func GetLogin(c *gin.Context) {
	info := GetSessionInfo(c)
	c.HTML(200, "login.html", gin.H{
		"SessionInfo": info,
	})
}

func PostLogin(c *gin.Context) {
	email := c.PostForm("email")
	pass := c.PostForm("password")
	m, err := db.FindUser(email, pass)

	if err == nil {
		Login(c, m)
		c.Redirect(302, "/")
	} else {
		// ログイン情報が違うとメッセージ出す
		c.Redirect(302, "/error")
	}
}

func GetSignup(c *gin.Context) {
	info := GetSessionInfo(c)
	c.HTML(200, "signup.html", gin.H{
		"SessionInfo": info,
	})
}

func PostSignup(c *gin.Context) {
	var n db.User
	n.Name = c.PostForm("newName")
	n.Email = c.PostForm("newEmail")
	pass := c.PostForm("password")
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	n.Password = hash
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
	log.Println(n)

	db.CreateUser(n)

	c.Redirect(302, "/")
}

func GetMy(c *gin.Context) {
	info := GetSessionInfo(c)
	history := db.GetHistory(info.UserID)

	c.HTML(200, "my.html", gin.H{
		"SessionInfo": info,
		"History":     history,
	})
}

func GetError(c *gin.Context) {
	info := GetSessionInfo(c)

	c.HTML(200, "error.html", gin.H{
		"SessionInfo": info,
		"Error":       c.Errors,
	})
}
