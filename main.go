package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var cookiestr string
var cookiesarr http.CookieJar

func main() {
	var err error

	// 创建内容
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建chrome实例
	c, err := chromedp.New(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// 运行任务

	// 导航
	err = c.Run(ctxt, chromedp.Navigate(`https://www.alimama.com/member/login.htm?forward=http%3A%2F%2Fpub.alimama.com%2Fmyunion.htm%3Fspm%3Da219t.7900221%2F1.a214tr8.2.446dfb5b8vg0Sx`))
	if err != nil {
		log.Fatal(err)
	}

	var site string
	for {
		err = c.Run(ctxt, chromedp.Location(&site))
		if err != nil {
			log.Fatal(err)
		}
		// 循环判断网址是否是登陆成功后的网址
		if string([]byte(site)[:34]) == "http://pub.alimama.com/myunion.htm" {
			break
		}
		time.Sleep(3 * time.Second)
	}
	err = c.Run(ctxt, getcookies())
	if err != nil {
		log.Fatal(err)
	}

	// 关闭浏览器
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// 等待浏览器完全关闭
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cookiestr)
	log.Println(cookiesarr)

}

func getcookies() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
			cookies, err := network.GetAllCookies().Do(ctxt, h)
			if err != nil {
				return err
			}
			cookiestr = ""
			var b []*http.Cookie

			for _, c := range cookies {
				if c.Domain == ".alimama.com" { //筛选作用域
					cookiestr += c.Name + "=" + c.Value + ";"
					var d http.Cookie
					d.Domain = c.Domain
					d.Path = c.Path
					d.Name = c.Name
					d.Value = c.Value
					b = append(b, &d)
				}
			}
			u, _ := url.Parse("http://*.alimama.com/")
			cookiesarr.SetCookies(u, b)
			return nil
		}),
	}
}
