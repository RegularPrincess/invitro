// invitro_parser project invitro_parser.go
package invitro_parser

import (
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
	parser.Store = db.GetConnection(9, dbConnInfo)
	return parser
}

/* На четвертом уровне парсинга(уровень оддельного анализа) происходит парсинг
   описания и запись, собранной в структуру Analysis, информации в базу */
func (this *Parser) scrapeLevel_4(analysis db.Analysis, urlDesc string) {
	doc, err := goquery.NewDocument(urlRoot + urlDesc)
	if err != nil {
		//На уровне отдельного исследования допустимы ошибки связанные с получением данных
		log.Println("level_4 ", err)
	}
	scriptTag := doc.Find("script").Eq(27)
	textAndGarb, _ := iconv.ConvertString(scriptTag.Text(), "windows-1251", "utf-8")
	re, _ := regexp.Compile("[^А-я ,.\n]")
	text := re.ReplaceAllString(textAndGarb, "")
	analysis.Description = text

	this.Store.AddAnalysis(&analysis)
}

/* Уровень 3 - уровень определения названия исследования */
func (this *Parser) scrapeLevel_3(analysis db.Analysis, urlAnlysis string) {
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
	tdNew.Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a")
		name := aTag.Text()
		name, _ = iconv.ConvertString(name, "windows-1251", "utf-8")
		link, _ := aTag.Attr("href")
		analysis.Name = name
		go this.scrapeLevel_4(analysis, link)
	})
}

/* Уровень определения подтипа исследования*/
func (this *Parser) scrapeLevel_2(analysis db.Analysis, urlSubtype string) {
	doc, err := goquery.NewDocument(urlRoot + urlSubtype)
	if err != nil {
		log.Fatalln("level_2 ", err)
	}
	divSubtype := doc.Find(".list-els")
	parsedDiv := divSubtype.Find("li")
	length := parsedDiv.Length()
	if length == 0 {
		this.scrapeLevel_3(analysis, urlSubtype)
	}
	parsedDiv.Each(func(i int, s *goquery.Selection) {
		aTag := s.Find("a")
		subtypeAnalysis := aTag.Text()
		subtypeAnalysis, _ = iconv.ConvertString(subtypeAnalysis, "windows-1251", "utf-8")
		link, _ := aTag.Attr("href")
		//Run the ScrapeLevel_3
		analysis.Subtype = subtypeAnalysis
		go this.scrapeLevel_3(analysis, link)
	})
}

/*Парсинг сайта работает подобно обходу дерева
  На каждом уровне запускается паралельнай обход низлежащих "ветвей"
  На каждом уровне работает своё правило парсинга*/
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
		go this.scrapeLevel_2(*analysis, link)
	})
	//	_ = <-this.finalChan
	//	time.Sleep(time.Second * 3)
	//	fmt.Println("Completed!")
}

//Проверка необходимости парсинга
func (this *Parser) NeedParse() bool {
	fill, err := this.Store.DBFilled()
	if err != nil {
		log.Fatalln(err)
	}
	return !fill
}
