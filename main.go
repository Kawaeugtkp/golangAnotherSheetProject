package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ここにURLを入力
	doc, err := goquery.NewDocument("https://tokubai.co.jp/news/backyard/articles/6076/edit")
	if err != nil {
		log.Fatal(err)
	}

	// ボタンが表示されるまで待機
	if err := waitForButton(ctx, doc, "公開取り下げ"); err != nil {
		log.Fatal(err)
	}

	// ボタンをクリック
	if err := clickButton(ctx, doc, "公開取り下げ"); err != nil {
		log.Fatal(err)
	}
}

// ボタンが表示されるまで待機する関数
func waitForButton(ctx context.Context, doc *goquery.Document, buttonText string) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if doc.Find("button:contains('" + buttonText + "')").Length() > 0 {
				return nil
			}
		}
	}
}

// ボタンをクリックする関数
func clickButton(ctx context.Context, doc *goquery.Document, buttonText string) error {
	button := doc.Find("button:contains('" + buttonText + "')").First()
	if button.Length() == 0 {
		return fmt.Errorf("ボタンが見つかりませんでした: %s", buttonText)
	}

	href, exists := button.Attr("href")
	if exists {
		// ボタンがリンクである場合は遷移する
		if err := chromedp.Run(ctx, chromedp.Navigate(href)); err != nil {
			return err
		}
	} else {
		// ボタンがJavaScriptを実行する場合はクリックする
		if err := chromedp.Run(ctx, chromedp.Click(button)); err != nil {
			return err
		}
	}

	return nil
}