package lineatgo

import (
    "net/http"
    "time"
    "fmt"
    "log"
    "strings"
    "context"
    "github.com/sclevine/agouti"
    "net/http/cookiejar"
    "net/url"
    "crypto/tls"
    "io/ioutil"
    "encoding/json"
    "github.com/mattn/go-scan"
)

type Bot struct {
    Name string
    AccountId string
    BotId string
}

type Api struct {
    MailAddress string
    Password string
    Client *http.Client
    XRT string
    CsrfToken string
    Bots []Bot
}

func NewApi(mail, pass string) *Api {
    return &Api{MailAddress: mail, Password: pass}
}

func (a *Api) GetBotIdByName(AccountName string) string {
    var bi string
    for _, b := range a.Bots{
        if b.Name == AccountName {
            bi = b.BotId
        }
    }
    return bi
}

func (a *Api) Login()  {
    driver := agouti.ChromeDriver(agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu"}), )
    if err := driver.Start(); err != nil {
        log.Fatalf("Failed to start driver:%v", err)
    }
    defer driver.Stop()

    page, err := driver.NewPage()
    if err != nil {
        log.Fatalf("Failed to open page:%v", err)
    }

    if err := page.Navigate("https://admin-official.line.me/"); err != nil { //ログイン画面を開く
        log.Fatalf("Failed to navigate:%v", err)
    }

    mailBox := page.FindByID("id")     //メアドボックス
    passBox := page.FindByID("passwd") //パスワードボックス
    mailBox.Fill(a.MailAddress) //メアドボックスにメアドを入力
    passBox.Fill(a.Password) //パスワードボックスにメアドを入力
    if err := page.FindByClass("MdBtn03Login").Submit(); err != nil { //ログインサブミットのボタン
        log.Fatalf("Failed to login:%v", err)
    }

    time.Sleep(1000 * time.Millisecond)
    PINcode, err := page.FindByClass("mdLYR04PINCode").Text()
    if err != nil {
        log.Println("メールアドレスまたはパスワードが間違っています。")
    }

    var limit bool
    ctx := context.Background()
    ctx, cancelTimer := context.WithCancel(ctx)
    fmt.Println(fmt.Sprintf("携帯のLINEで以下のPINコードを入力してください: %v", PINcode))//PINコードを表示
    go timer(140000, ctx, &limit)
    for {
        title, _ := page.Title()
        if strings.Contains(title, "LINE@ MANAGER") {
            limit = false
            cancelTimer()
            break
        }

        if limit {
            log.Println("時間切れです。")
            limit = false
        }
    }
    c, _ := page.GetCookies()
    a.createClient(c)
    a.getBotInfo()
}

func (a *Api) createClient(c []*http.Cookie) {
    jar, _ := cookiejar.New(nil)
    u, _ := url.Parse("https://admin-official.line.me/")
    jar.SetCookies(u, c)
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{ServerName: "*.line.me"},
    }
    a.Client = &http.Client{
        Transport: tr,
        Jar: jar,
    }
}

func (a *Api) getBotInfo() {
    request, _ := http.NewRequest("GET", "https://admin-official.line.me/api/basic/bot/list?_=1510425201579&count=10&page=1", nil)
    request.Header.Set("Content-Type", "application/json;charset=UTF-8")
    resp, _ := a.Client.Do(request)
    defer resp.Body.Close()
    cont, _ := ioutil.ReadAll(resp.Body)
    var ij interface{}
    json.Unmarshal([]byte(string(cont)), &ij)
    var v Bot
    var b []Bot
    for i:=0; i<strings.Count(string(cont), "botId"); i++ {
        scan.ScanTree(ij, fmt.Sprintf("/list[%v]/displayName", i), &v.Name)
        scan.ScanTree(ij, fmt.Sprintf("/list[%v]/botId", i), &v.BotId)
        b = append(b, v)
    }
    a.Bots = b
}

func timer(wait int, ctx context.Context, l *bool) {
    time.Sleep(time.Duration(wait) * time.Millisecond)
    *l = true
}