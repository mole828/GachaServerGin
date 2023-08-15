package src

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

func GetLimitedPools() []string {
	response, err := http.Get("https://prts.wiki/w/%E5%8D%A1%E6%B1%A0%E4%B8%80%E8%A7%88")
	if err != nil {
		return []string{"云间清醒梦", "真理孑然", "万象伶仃"}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Error(err)
		}
	}(response.Body)
	if response.StatusCode != 200 {
		Logger.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
		return nil
	}
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		Logger.Error(err)
		return nil
	}
	var labels []string
	document.Find("table tbody tr").Each(func(_ int, selection *goquery.Selection) {
		selection.Find("td").Find("a").Each(func(_ int, selection *goquery.Selection) {
			full := selection.Text()
			title, exist := selection.Attr("title")
			if strings.Index(full, "限定寻访") != -1 && exist {
				labels = append(labels, title)
			}
		})
	})
	return labels

}
