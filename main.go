package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"strings"
	
    "cloud.google.com/go/translate"
    "golang.org/x/text/language"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "golang.org/x/crypto/bcrypt"
    
	 "github.com/jinzhu/gorm"
	  _ "github.com/jinzhu/gorm/dialects/postgres"

)


type LangSet struct {
	Language string
	Tag string `gorm:"type:varchar(10);unique_index"`
	Kotae string
}

type User struct {
	gorm.Model
	Name string
	Email string  `gorm:"type:varchar(100);unique_index"`
	Password []byte
}

type SessionInfo struct {
	Name string //ログインしているユーザの表示名
	Email string //ログインしているユーザのemail
	IsSessionAlive bool //セッションが生きているかどうか
}

type History struct {
	gorm.Model
	Word string
	Lngs string
	Name string
	Email string
}

type MyHis struct {
	Name string
	Email string
	Word string
	Lngs string
}

func (Eigo *LangSet) translater(kotoba string) {
	c := context.Background()
	
	// Clientをつくる
	client, err := translate.NewClient(c)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	target, err := language.Parse(Eigo.Tag)
        if err != nil {
                log.Fatalf("Failed to parse target language: %v", err)
        }

		translations, err := client.Translate(c, []string{kotoba}, target, nil)
		if err != nil {
               log.Fatalf("Failed to translate text: %v", err)
	    }
	    Eigo.Kotae = translations[0].Text
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	log.Println("セッション取得")
	session.Clear()
	log.Println("クリア処理")
	session.Save()
}

func Login ( c *gin.Context, u User) {
	session := sessions.Default(c)
	session.Set("Name", u.Name)
	session.Set("Email", u.Email)
	session.Save()
}

func GetSessionInfo (c *gin.Context) SessionInfo {
	var info SessionInfo
	session := sessions.Default(c)
	user_n := session.Get("Name")
	user_e := session.Get("Email")
	if user_n == nil || user_e == nil {
		info = SessionInfo{
			Name: "", Email: "", IsSessionAlive: false,
		}
	} else {
		info = SessionInfo{
			Name: user_n.(string),
			Email: user_e.(string),
			IsSessionAlive: true,
			}
		log.Println(info)
		}
	return info
}

func Hists(limit int, c *gin.Context) (hists []MyHis) {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=honey dbname=postgres password=password sslmode=disable")
	if err != nil {
    panic("データベースへの接続に失敗しました"+err.Error())
	}
	defer db.Close()
	session := sessions.Default(c)
	//ここまでお約束
	//ここからemailで検索
	email := session.Get("Email")
	rows, err := db.Table("histories").Where("email = ?", email).Select("word, lngs").Rows() //SELECT word, lngs FROM histories WHERE email = 'takuya.kometan@gmail.com' limit 1;
	defer rows.Close()
	//繰り返しで取得
	fmt.Println(rows)
	for rows.Next() {
		var h MyHis
		err = rows.Scan(&h.Word, &h.Lngs)
		if err != nil {
			fmt.Printf("Scan Error!!!! err:%v\n", err)
			return
		}
		hists = append(hists, h)
	}
	fmt.Println(hists)
	return hists
}

func Replace(s string) string {
	s = strings.Replace(s, "&lt;", "<", -1)
	s = strings.Replace(s, "&gt;", "<", -1)
	return s
}

func main() {
	en := LangSet{"英語", "en", ""}
	de := LangSet{"ドイツ語", "de", ""}
	irish := LangSet{"アイルランド語", "ga", ""}
	ar := LangSet{"アラビア語", "ar", ""}
	el := LangSet{"ギリシャ語", "el", ""}
	espe := LangSet{"エスペラント", "eo", ""}
	espanol := LangSet{"スペイン語", "es", ""}
	fr := LangSet{"フランス語", "fr", ""}
	ita := LangSet{"イタリア語", "it", ""}
	dut := LangSet{"オランダ語", "nl", ""}
	afr := LangSet{"アフリカーンス語", "af", ""}
	fin := LangSet{"フィンランド語", "fi", ""}
	sweden := LangSet{"スウェーデン語", "sv", ""}
	nor := LangSet{"ノルウェー語", "no", ""}
	ice := LangSet{"アイスランド語", "is", ""}
	ru := LangSet{"ロシア語", "ru", ""}
	polish := LangSet{"ポーランド語", "pl", ""}
	bul := LangSet{"ブルガリア語", "bg", ""}
	ukr := LangSet{"ウクライナ語", "uk", ""}
	kyrgyz := LangSet{"キルギス語", "ky", ""}
	turkish := LangSet{"トルコ語", "tr", ""}
	swa := LangSet{"スワヒリ語", "sw", ""}
	
	languages := []LangSet{en, irish, ar, el, espe, espanol, fr, ita, dut, afr, fin, de, sweden, nor, ice, ru, polish, bul, ukr, kyrgyz, turkish, swa}
	
	router := gin.Default()
	router.Static("/templates", "./templates")
    router.LoadHTMLGlob("templates/*.html") // 事前にテンプレートをロード 相対パス
    store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
    
    db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=honey dbname=postgres password=password sslmode=disable")
	if err != nil {
    panic("データベースへの接続に失敗しました"+err.Error())
	}
	defer db.Close()
	
	db.AutoMigrate(&User{})
	db.AutoMigrate(&History{})

    router.GET("/", func(c *gin.Context){
    	info := GetSessionInfo(c)
    	
    	c.HTML(200, "index.html", gin.H{
    		"SessionInfo": info,
    	})
    })
    
    router.POST("/new", func(c *gin.Context){
    	text := c.PostForm("word")
    	
    	for i := 0; i < len(languages); i++{
    		languages[i].translater(text)
    		fmt.Println(languages[i].Language, "：", languages[i].Kotae)
    	}
    	// 履歴に残すために
    	session := sessions.Default(c)
    	var h History
    	h.Word = text
    	for i := 0; i < len(languages); i++ {
    		h.Lngs = h.Lngs + "<tr><td>" + languages[i].Language + "</td>" +"<td>" + languages[i].Kotae + "</td></tr>"
    	}
    	h.Name = session.Get("Name").(string)
    	h.Email = session.Get("Email").(string)
    	now := time.Now()
    	h.CreatedAt = now
    	h.UpdatedAt = now
    	log.Println(h)
    	
    	db.NewRecord(h)
    	db.Create(&h)
    	if db.NewRecord(h) == false {
			log.Printf("History %d Recorded\n", h.ID)
		} else {
			log.Println("History not created") //エラー内容を表示したい http://gorm.io/ja_JP/docs/error_handling.html
		}
    	// ここまで
    	info := GetSessionInfo(c)
       	c.HTML(200, "translations.html", gin.H{
       		"text": text,
       		"languages": languages,
       		"SessionInfo": info,
       	})
    })
    
    router.GET("/login", func(c *gin.Context){
    	info := GetSessionInfo(c)
    	c.HTML(200, "login.html", gin.H{
    		"SessionInfo": info,
    		})
    })
    
    router.POST("/login", func(c *gin.Context){
    	var n, m, user User
    	n.Email = c.PostForm("email")
    	pass := c.PostForm("password")
    	db.Where("email = ?", n.Email).First(&user).Scan(&m)
    	err = bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(pass))
    	if n.Email == m.Email && err == nil {
    		Login(c, m)
    		info := GetSessionInfo(c)
    		c.HTML(200, "index.html", gin.H{
    			"SessionInfo": info,
    		})
    	} else {
    		// ログイン情報が違うとメッセージ出す
    		info := GetSessionInfo(c)
    		c.HTML(400, "error.html", gin.H{
    			"SessionInfo": info,
    		})
    	}
    })
    
    router.GET("/logout", func(c *gin.Context){
    	Logout(c)
    	info := GetSessionInfo(c)
    	c.HTML(200, "index.html", gin.H{
    		"SessionInfo": info,
    	})
    })
    
    router.GET("/signup", func (c *gin.Context){
    	info := GetSessionInfo(c)
    	c.HTML(200, "signup.html", gin.H{
    		"SessionInfo": info,
    	})
    })
    
    router.POST("/signup", func (c *gin.Context){
    	var n User
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
    	
		db.NewRecord(n)
		db.Create(&n)
		if db.NewRecord(n) == false {
			log.Printf("User %s Created\n", n.Name)
		} else {
			log.Println("User not created") //エラー内容を表示したい http://gorm.io/ja_JP/docs/error_handling.html
		}
			info := GetSessionInfo(c)
	 		c.HTML(200, "index.html", gin.H{
	 			"SessionInfo": info,
	 		})
    })
    
    router.GET("/my", func(c *gin.Context){
    	info := GetSessionInfo(c)
    	hst := Hists(10, c)
    	// 文字実体参照の問題を解決。強引に変換
    	for i :=0 ; i < len(hst); i++ {
    		hst[i].Lngs = Replace(hst[i].Lngs)
    	}
    	c.HTML(200, "my.html", gin.H{
    		"SessionInfo": info,
    		"Hst": hst,
    	})
    })

    router.Run()
}