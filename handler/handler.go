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

type sessionInfo struct {
	UserID         int
	Name           string //ログインしているユーザの表示名
	IsSessionAlive bool   //セッションが生きているかどうか
}

type pair struct {
	Language string
	Code     string
	Parent   string
	Response string
}

//Logout let user log out. Clear session
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	log.Println("セッション取得")
	session.Clear()
	log.Println("クリア処理")
	session.Save()
	c.Redirect(302, "/")
}

// Login writes user info to session
func Login(c *gin.Context, u db.User) {
	session := sessions.Default(c)
	session.Set("Name", u.Name)
	session.Set("Email", u.Email)
	session.Save()
}

//getSessionInfo get current user info
func getSessionInfo(c *gin.Context) sessionInfo {
	var info sessionInfo
	session := sessions.Default(c)
	userID := session.Get("UserID")
	userName := session.Get("Name")
	if userID == nil {
		info = sessionInfo{
			UserID: -1, Name: "", IsSessionAlive: false,
		}
	} else {
		info = sessionInfo{
			UserID:         userID.(int),
			Name:           userName.(string),
			IsSessionAlive: true,
		}
		log.Println(info)
	}
	return info
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//Index show index page
func Index(c *gin.Context) {
	info := getSessionInfo(c)
	c.HTML(200, "index.html", gin.H{
		"sessionInfo": info,
	})
}

//New tell GetTranslation what is the word to translate
func New(c *gin.Context) {
	text := c.PostForm("word")
	uri := fmt.Sprintf("/translate/%s", text)
	c.Redirect(302, uri)
}

//GetTranslation translate word using google cloud translate.
//If you update languages, update template file also.
func GetTranslation(c *gin.Context) {
	pairs := []pair{
		{"日本語", "ja", "root", ""},
		{"英語", "en", "germanic", ""},
		{"ドイツ語", "de", "germanic", ""},
		{"アイルランド語", "ga", "celtic", ""},
		{"アラビア語", "ar", "afro-asiatic", ""},
		{"ギリシャ語", "el", "hellenic", ""},
		{"エスペラント", "eo", "root", ""},
		{"スペイン語", "es", "italic", ""},
		{"フランス語", "fr", "italic", ""},
		{"イタリア語", "it", "italic", ""},
		{"オランダ語", "nl", "germanic", ""},
		{"アフリカーンス語", "af", "germanic", ""},
		{"フィンランド語", "fi", "uralic", ""},
		{"スウェーデン語", "sv", "germanic", ""},
		{"ノルウェー語", "no", "germanic", ""},
		{"アイスランド語", "is", "germanic", ""},
		{"ロシア語", "ru", "balt-slavic", ""},
		{"ポーランド語", "pl", "balt-slavic", ""},
		{"ブルガリア語", "bg", "balt-slavic", ""},
		{"ウクライナ語", "uk", "balt-slavic", ""},
		{"キルギス語", "ky", "turkic", ""},
		{"トルコ語", "tr", "turkic", ""},
		{"スワヒリ語", "sw", "niger-congo", ""},
	}

	text := c.Param("text")
	texts := []string{text}
	fmt.Println(pairs)
	// Clientをつくる
	client, err := translate.NewTranslationClient(c)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	for i := range pairs {

		request := &translatepb.TranslateTextRequest{
			Contents:           texts,
			TargetLanguageCode: pairs[i].Code,
			Parent:             "projects/honyac-konjac",
		}

		translation, err := client.TranslateText(c, request)
		if err != nil {
			log.Fatalf("Failed to translate text: %v", err)
		}
		t := translation.Translations[0].GetTranslatedText()
		pairs[i].Response = strings.ReplaceAll(t, "&#39;", "'") //フランス語でL'EgypteがL&#39;Egupteとなる問題
		fmt.Println(pairs[i])                                   //デバッグ用
	}

	// 履歴に残すために
	info := getSessionInfo(c)
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
		"jsmind":      true,
	})
}

func GetLogin(c *gin.Context) {
	info := getSessionInfo(c)
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
	info := getSessionInfo(c)
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
	info := getSessionInfo(c)
	history := db.GetHistory(info.UserID)

	c.HTML(200, "my.html", gin.H{
		"SessionInfo": info,
		"History":     history,
	})
}

func GetError(c *gin.Context) {
	info := getSessionInfo(c)

	c.HTML(200, "error.html", gin.H{
		"SessionInfo": info,
		"Error":       c.Errors,
	})
}
