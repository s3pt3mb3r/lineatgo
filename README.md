[![GoDoc](https://godoc.org/github.com/s3pt3mb3r/lineatgo?status.svg)](https://godoc.org/github.com/s3pt3mb3r/lineatgo)
[![GoWalker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/s3pt3mb3r/lineatgo)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
# lineatgo
This is an unofficial LINE@ API that was implemented in pure go

## installation
```
go get github.com/s3pt3mb3r/lineatgo
```

## Usage
```go
package main

import (
    "github.com/s3pt3mb3r/lineatgo"
    "log"
    "fmt"
    "os
)

func main()  {
    api := lineatgo.NewApi("MAIL_ADDRESS", "PASSWORD")
    bot, err := api.NewBot("@LINEID")
    if err != nil {
        log.Println(err)
    }
    bot.PostText("This is a post via lineatgo")
    url := bot.GetAuthURL(lineatgo.Administrator)
    //else: lineatgo.Operator, lineatgo.LimitedOperator, lineatgo.Messenger
    fmt.Println(url)

    qr := bot.GetQRCode()

    file, err := os.OpenFile("test.png", os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    file.Write(qr)
}
```

## What is "@LINEID"?
![lineId](https://user-images.githubusercontent.com/33578715/33206339-35ebeb4a-d14d-11e7-94cd-a1b8a69d4031.jpg)

## Todo
- [x] Enable to select authority type in getAuth function
- [x] Enable to Delete paymaster user's clearance
- [x] Enable to Post some text on time line
- [ ] Enable to Post image or video on time line
- [ ] Fix DeletePostAll function



### At last
Probably, being overlook some factors, I can't code Login() function without web driver
(I'm working at dev branch. See login2 function)

If it's possible, please make Login() function more better and send pull request:)
