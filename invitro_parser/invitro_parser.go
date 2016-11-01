// invitro_parser project invitro_parser.go
package invitro_parser

import (
	"fmt"
	"log"
	"regexp"

	db "github.com/RegularPrincess/invitro/invitro_model"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

const urlRoot string = "https://www.invitro.ru"

type Parser struct {
	Store     db.Store
	finalChan chan bool
}

func GetParser(dbConnInfo string) *Parser {
	var parser = new(Parser)
	parser.finalChan = make(chan bool)
	parser.Store = db.GetConnection(15, dbConnInfo)
	return parser
}

func (this *Parser) CloseConn() {

}

func (this *Parser) scrapeLevel_4(analysis db.Analysis, urlDesc string, isLast bool) {
	doc, err := goquery.NewDocument(urlRoot + urlDesc)
	if err != nil {
		log.Println("level_4 ", err)
		if isLast {
			this.finalChan <- true
		}
		return
	}
	scriptTag := doc.Find("script").Eq(27)
	textAndGarb, _ := iconv.ConvertString(scriptTag.Text(), "windows-1251", "utf-8")
	re, _ := regexp.Compile("[^А-я ,.\n]")
	text := re.ReplaceAllString(textAndGarb, "")
	analysis.Description = text

	this.Store.AddAnalysis(&analysis)
	if isLast {
		fmt.Print("0")
		this.finalChan <- true
	}
}

func (this *Parser) scrapeLevel_3(analysis db.Analysis, urlAnlysis string, isLast bool) {
	doc, err := goquery.NewDocument(urlRoot + urlAnlysis)
	if err != nil {
		log.Fatalln("level_3 ", err)
	}
	td := doc.Find(".data-table").Find("td")

	tdNew := new(goquery.Selection)
	td.Each(func(i int, s *goquery.Selection) {
		if s.HasClass("name elem_list_0") || s.HasClass("name elem_list_1") {
			tdNew = tdNew.AddSelection(s)
		}
	})
	length := tdNew.Length()
	length--
	tdNew.Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a")
		name := aTag.Text()
		name, _ = iconv.ConvertString(name, "windows-1251", "utf-8")
		link, _ := aTag.Attr("href")
		analysis.Name = name
		if isLast {
			if i < length {
				go this.scrapeLevel_4(analysis, link, false)
			} else {
				go this.scrapeLevel_4(analysis, link, true)
			}
		} else {
			go this.scrapeLevel_4(analysis, link, false)
		}
	})
}

func (this *Parser) scrapeLevel_2(analysis db.Analysis, urlSubtype string, isLast bool) {
	doc, err := goquery.NewDocument(urlRoot + urlSubtype)
	if err != nil {
		log.Fatalln("level_2 ", err)
	}
	divSubtype := doc.Find(".list-els")
	parsedDiv := divSubtype.Find("li")
	length := parsedDiv.Length()
	if length == 0 {
		this.scrapeLevel_3(analysis, urlSubtype, isLast)
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
				go this.scrapeLevel_3(analysis, link, false)
			} else {
				go this.scrapeLevel_3(analysis, link, true)
			}
		} else {
			go this.scrapeLevel_3(analysis, link, false)
		}
	})
}

func (this *Parser) Scrape(urlType string) {
	//fmt.Println("Waiting..")

	err := this.Store.Clean()
	if err != nil {
		log.Fatalln(err)
	}
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
			go this.scrapeLevel_2(*analysis, link, false)
		} else {
			go this.scrapeLevel_2(*analysis, link, true)
		}
	})
	//	_ = <-this.finalChan
	//	time.Sleep(time.Second * 3)
	//	fmt.Println("Completed!")
}

func (this *Parser) NeedParse() bool {
	fill, err := this.Store.DBFilled()
	if err != nil {
		log.Fatalln(err)
	}
	return !fill
}
