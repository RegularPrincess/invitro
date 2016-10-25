// invitro project main.go
package invitro_parser

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

const urlRoot string = "https://www.invitro.ru"

var store Store

//Need Fix (crutch)
func ScrapeLevel_4(analysis Analysis, urlDesc string) {
	doc, err := goquery.NewDocument(urlRoot + urlDesc)
	if err != nil {
		log.Println("level_4 ", err)
		return
	}
	scriptTag := doc.Find("script").Eq(27)
	textAndGarb, _ := iconv.ConvertString(scriptTag.Text(), "windows-1251", "utf-8")
	re, _ := regexp.Compile("[^А-я ,.\n]")
	text := re.ReplaceAllString(textAndGarb, "")
	analysis.description = text

	store.addAnalysis(&analysis)
}

func ScrapeLevel_3(analysis Analysis, urlAnlysis string) {
	doc, err := goquery.NewDocument(urlRoot + urlAnlysis)
	if err != nil {
		log.Fatalln("level_3 ", err)
	}
	td := doc.Find(".data-table").Find("td")
	td.Each(func(i int, s *goquery.Selection) {
		if s.HasClass("name elem_list_0") || s.HasClass("name elem_list_1") {
			aTag := s.Find("a")
			name := aTag.Text()
			name, _ = iconv.ConvertString(name, "windows-1251", "utf-8")
			link, _ := aTag.Attr("href")
			analysis.name = name

			go ScrapeLevel_4(analysis, link)
		}
	})
}

func ScrapeLevel_2(analysis Analysis, urlSubtype string) {
	doc, err := goquery.NewDocument(urlRoot + urlSubtype)
	if err != nil {
		log.Fatalln("level_2 ", err)
	}
	divSubtype := doc.Find(".list-els")

	parsedDiv := divSubtype.Find("li")
	if parsedDiv.Length() == 0 {
		ScrapeLevel_3(analysis, urlSubtype)
	}
	parsedDiv.Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a")
		subtypeAnalysis := aTag.Text()
		subtypeAnalysis, _ = iconv.ConvertString(subtypeAnalysis, "windows-1251", "utf-8")
		link, _ := aTag.Attr("href")
		//Run the ScrapeLevel_3
		analysis.subtype = subtypeAnalysis

		go ScrapeLevel_3(analysis, link)
	})
}

func Scrape(urlType string, dbConnInfo string) {
	store = GetConnection(19, dbConnInfo)
	doc, err := goquery.NewDocument(urlRoot + urlType)
	if err != nil {
		log.Fatalln(err)
	}
	doc.Find(".group-name-brd").Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a")
		typeAnalysis := aTag.Text()
		typeAnalysis, _ = iconv.ConvertString(typeAnalysis, "windows-1251", "utf-8")
		link, _ := aTag.Attr("href")
		analysis := new(Analysis)
		analysis.kind = typeAnalysis

		go ScrapeLevel_2(*analysis, link)
	})
	var input string
	fmt.Println("Waiting..")
	fmt.Scanln(&input)
}
