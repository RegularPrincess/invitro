// invitro_parser project invitro_parser.go
package invitro_parser

import (
	"fmt"
	db "invitro_model"
	"log"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

const urlRoot string = "https://www.invitro.ru"

var store db.Store

var finalChan chan bool

func scrapeLevel_4(analysis db.Analysis, urlDesc string, isLast bool) {
	doc, err := goquery.NewDocument(urlRoot + urlDesc)
	if err != nil {
		log.Println("level_4 ", err)
		return
	}
	scriptTag := doc.Find("script").Eq(27)
	textAndGarb, _ := iconv.ConvertString(scriptTag.Text(), "windows-1251", "utf-8")
	re, _ := regexp.Compile("[^А-я ,.\n]")
	text := re.ReplaceAllString(textAndGarb, "")
	analysis.Description = text

	store.AddAnalysis(&analysis)
	if isLast {
		finalChan <- true
	}
}

func scrapeLevel_3(analysis db.Analysis, urlAnlysis string, isLast bool) {
	doc, err := goquery.NewDocument(urlRoot + urlAnlysis)
	if err != nil {
		log.Fatalln("level_3 ", err)
	}
	td := doc.Find(".data-table").Find("td")
	length := td.Length()
	length = length - 2
	td.Each(func(i int, s *goquery.Selection) {
		if s.HasClass("name elem_list_0") || s.HasClass("name elem_list_1") {
			aTag := s.Find("a")
			name := aTag.Text()
			name, _ = iconv.ConvertString(name, "windows-1251", "utf-8")
			link, _ := aTag.Attr("href")
			analysis.Name = name
			if isLast {
				if i < length {
					go scrapeLevel_4(analysis, link, false)
				} else {
					go scrapeLevel_4(analysis, link, true)
				}
			} else {
				go scrapeLevel_4(analysis, link, false)
			}
		}
	})
}

func scrapeLevel_2(analysis db.Analysis, urlSubtype string, isLast bool) {
	doc, err := goquery.NewDocument(urlRoot + urlSubtype)
	if err != nil {
		log.Fatalln("level_2 ", err)
	}
	divSubtype := doc.Find(".list-els")
	parsedDiv := divSubtype.Find("li")
	length := parsedDiv.Length()
	if length == 0 {
		scrapeLevel_3(analysis, urlSubtype, isLast)
	}
	length--
	parsedDiv.Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a")
		subtypeAnalysis := aTag.Text()
		subtypeAnalysis, _ = iconv.ConvertString(subtypeAnalysis, "windows-1251", "utf-8")
		link, _ := aTag.Attr("href")
		//Run the ScrapeLevel_3
		analysis.Subtype = subtypeAnalysis
		if isLast {
			if i < length {
				go scrapeLevel_3(analysis, link, false)
			} else {
				go scrapeLevel_3(analysis, link, true)
			}
		} else {
			go scrapeLevel_3(analysis, link, false)
		}
	})
}

func Scrape(urlType string, dbConnInfo string) {
	fmt.Println("Waiting..")
	finalChan = make(chan bool)
	store = db.GetConnection(19, dbConnInfo)
	doc, err := goquery.NewDocument(urlRoot + urlType)
	if err != nil {
		log.Fatalln(err)
	}
	parsedDoc := doc.Find(".group-name-brd")
	length := parsedDoc.Length()
	length--
	parsedDoc.Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a")
		typeAnalysis := aTag.Text()
		typeAnalysis, _ = iconv.ConvertString(typeAnalysis, "windows-1251", "utf-8")
		link, _ := aTag.Attr("href")
		analysis := new(db.Analysis)
		analysis.Kind = typeAnalysis
		if i < length {
			go scrapeLevel_2(*analysis, link, false)
		} else {
			go scrapeLevel_2(*analysis, link, true)
		}
	})
	_ = <-finalChan
	time.Sleep(time.Second)
	fmt.Println("Completed!")
}
