// Make a function to  scrape the data from the website (https://arknights.wiki.gg/wiki/Aak)
// 1. The function should execute a background job
// 2. The function should return object as follows: {
//    "operator_name": "",
//    "short_description": "",}
//    "class": "",
//    "branch": "",
//    "faction": "",}
//    "position": "",}
//    "tags": [""],}
//    "trait": "",}

package internal

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type Operator struct {
	OperatorName     string
	ShortDescription string
	Class            string
	Branch           string
	Faction          string
	Position         string
	Tags             []string
	Trait            string
}

func Scrapper() (Operator, error) {
	var operator Operator
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		c := colly.NewCollector()

		c.OnHTML("table", func(e *colly.HTMLElement) {
			e.ForEach("tr", func(i int, e *colly.HTMLElement) {
				switch i {
				case 1:
					operator.OperatorName = e.ChildText("td")
				case 2:
					operator.ShortDescription = e.ChildText("td")
				case 3:
					operator.Class = e.ChildText("td")
				case 4:
					operator.Branch = e.ChildText("td")
				case 5:
					operator.Faction = e.ChildText("td")
				case 6:
					operator.Position = e.ChildText("td")
				case 7:
					tags := strings.Split(e.ChildText("td"), ",")
					for _, tag := range tags {
						operator.Tags = append(operator.Tags, strings.TrimSpace(tag))
					}
				case 8:
					operator.Trait = e.ChildText("td")
				}
			})
		})

		c.Visit("https://arknights.wiki.gg/wiki/Aak")
	}()

	wg.Wait()
	return operator, nil
}

func PrintOperator(operator Operator) {
	fmt.Println("Operator Name:", operator.OperatorName)
	fmt.Println("Short Description:", operator.ShortDescription)
	fmt.Println("Class:", operator.Class)
	fmt.Println("Branch:", operator.Branch)
	fmt.Println("Faction:", operator.Faction)
	fmt.Println("Position:", operator.Position)
	fmt.Println("Tags:", operator.Tags)
	fmt.Println("Trait:", operator.Trait)
}
