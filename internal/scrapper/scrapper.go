package internal

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"chat-ak-wikia/internal/model"

	"github.com/gocolly/colly"
)

func updateOrCreateAttribute(attributes map[string]model.Attribute, key string, field string, value interface{}) {

	valueInt := 0
	valueStr := ""
	if v, ok := value.(int); ok {
		valueInt = v
	} else if v, ok := value.(string); ok {
		valueStr = v
	}

	if attr, ok := attributes[key]; ok {
		switch field {
		case "HP":
			attr.HP = valueInt
		case "Atk":
			attr.Atk = valueInt
		case "Def":
			attr.Def = valueInt
		case "Res":
			attr.Res = valueInt
		case "RedeploymentTime":
			attr.RedeploymentTime = valueStr
		case "DPcost":
			attr.DPcost = valueInt
		case "Block":
			attr.Block = valueInt
		case "AttackInterval":
			attr.AttackInterval = valueStr
		}
		attributes[key] = attr
	} else {
		attr := model.Attribute{}
		switch field {
		case "HP":
			attr.HP = valueInt
		case "Atk":
			attr.Atk = valueInt
		case "Def":
			attr.Def = valueInt
		case "Res":
			attr.Res = valueInt
		case "RedeploymentTime":
			attr.RedeploymentTime = valueStr
		case "DPcost":
			attr.DPcost = valueInt
		case "Block":
			attr.Block = valueInt
		case "AttackInterval":
			attr.AttackInterval = valueStr
		}
		attributes[key] = attr
	}
}

func Scrapper(maxVisits int32, link string, c *colly.Collector) ([]model.Operator, error) {
	operators := make([]model.Operator, 0, 10)
	var wg sync.WaitGroup
	var visitedCount int32 = 0

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Clone another collector to scrape operators
		operatorCollector := c.Clone()

		// On every a element which has href attribute call callback
		c.OnHTML("table.mrfz-wtable td[align='center'] > a[href^='/wiki/']", func(e *colly.HTMLElement) {

			link := e.Attr("href")

			if visitedCount < maxVisits { // Check limit before visiting
				fmt.Printf("Link found: %q -> %s\n", e.Text, link)
				visitedCount++ // Increment the counter *before* visiting

				// Only those links are visited which are in AllowedDomains
				operatorCollector.Visit(e.Request.AbsoluteURL(link))
			} else {
				// fmt.Println("Reached maximum visit limit. Skipping:", link)
			}

		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		operatorCollector.OnHTML("main", func(e *colly.HTMLElement) {

			// Define operator
			operator := model.Operator{}

			// Extract operator name
			operator.OperatorName = e.ChildText("span.mw-page-title-main")

			// Extract operator summary
			e.ForEach("tr", func(i int, e *colly.HTMLElement) {
				switch e.ChildText("td:nth-of-type(1)") {
				case "Class":
					operator.Class = e.ChildText("td:nth-of-type(2)")
				case "Branch":
					operator.Branch = e.ChildText("td:nth-of-type(2)")
				case "Faction":
					operator.Faction = e.ChildText("td:nth-of-type(2)")
				case "Position":
					operator.Position = e.ChildText("td:nth-of-type(2)")
				case "Tags":
					tags := strings.Split(e.ChildText("td:nth-of-type(2)"), ",")
					for _, tag := range tags {
						operator.Tags = append(operator.Tags, strings.TrimSpace(tag))
					}
				case "Trait":
					operator.Trait = e.ChildText("td:nth-of-type(2)")
				}
			})

			// Extract operator information
			e.ForEach("aside.portable-infobox", func(i int, e *colly.HTMLElement) {

				e.ForEach("div.pi-item.pi-data.pi-item-spacing.pi-border-color", func(i int, e *colly.HTMLElement) {
					value := e.ChildText("div.pi-data-value")
					switch i {
					case 0:
						operator.CharacterInfo.ThemeSong = value
					case 1:
						operator.CharacterInfo.InternalName = value
					case 2:
						operator.CharacterInfo.BasedOn = value
					case 3:
						operator.CharacterInfo.Etymology = value
					case 7:
						operator.CharacterInfo.FileNo = value
					case 8:
						operator.CharacterInfo.OperatorRecord = value
					case 9:
						operator.CharacterInfo.ParadoxSimulation = value
					case 10:
						operator.CharacterInfo.Illustrator = value
					case 11:
						operator.CharacterInfo.JapaneseCV = value
					case 12:
						operator.CharacterInfo.MandarinCV = value
					case 13:
						operator.CharacterInfo.CantoneseCV = value
					case 14:
						operator.CharacterInfo.EnglishCV = value
					case 15:
						operator.CharacterInfo.KoreanCV = value
					case 16:
						operator.CharacterInfo.Gender = value
					case 17:
						operator.CharacterInfo.CombatExperience = value
					case 18:
						operator.CharacterInfo.PlaceOfBirth = value
					case 19:
						operator.CharacterInfo.DateOfBirth = value
					case 20:
						operator.CharacterInfo.Race = value
					case 21:
						operator.CharacterInfo.Height = value
					case 22:
						operator.CharacterInfo.InfectionStatus = value
					case 23:
						operator.CharacterInfo.PhysicalStrength = value
					case 24:
						operator.CharacterInfo.Mobility = value
					case 25:
						operator.CharacterInfo.PhysicalResilience = value
					case 26:
						operator.CharacterInfo.TacticalAcumen = value
					case 27:
						operator.CharacterInfo.CombatSkill = value
					case 28:
						operator.CharacterInfo.OriginiumArtsAssimilation = value
					}
				})
			})

			// Extract operator potential
			e.ForEach("table.mrfz-wtable#operator-potential-table tr", func(i int, e *colly.HTMLElement) {

				if e.DOM.Find("td").Length() == 0 {
					return
				} else if e.DOM.Find("td").Length() == 1 {
					return
				} else {
					potential := model.Potential{}
					potential.Level = e.ChildAttr("td span", "title")
					potential.Effect = e.ChildText("td")
					operator.Potentials = append(operator.Potentials, potential)
				}
			})

			// Extract operator ascension words
			operator.ShortDescription = e.ChildText("div[style='margin:0 5px; padding:0 1em;'] div:nth-of-type(1)")

			// Extract operator status
			operator.Attributes = make(map[string]model.Attribute)
			e.ForEach("table#operator-attribute-table tr", func(i int, e *colly.HTMLElement) {
				switch e.ChildText("th") {
				case "HP":
					baseHP, _ := strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
					updateOrCreateAttribute(operator.Attributes, "Base", "HP", baseHP)

					baseMaxHP, _ := strconv.Atoi(e.ChildText("td:nth-of-type(2)"))
					updateOrCreateAttribute(operator.Attributes, "BaseMax", "HP", baseMaxHP)

					elite1HP, _ := strconv.Atoi(e.ChildText("td:nth-of-type(3)"))
					updateOrCreateAttribute(operator.Attributes, "Elite1", "HP", elite1HP)

					elite2HP, _ := strconv.Atoi(e.ChildText("td:nth-of-type(4)"))
					updateOrCreateAttribute(operator.Attributes, "Elite2", "HP", elite2HP)

					trustBonusHP, _ := strconv.Atoi(e.ChildText("td:nth-of-type(5)"))
					operator.TrustBonus.HP = trustBonusHP
				case "ATK":
					baseAtk, _ := strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
					updateOrCreateAttribute(operator.Attributes, "Base", "Atk", baseAtk)

					baseMaxAtk, _ := strconv.Atoi(e.ChildText("td:nth-of-type(2)"))
					updateOrCreateAttribute(operator.Attributes, "BaseMax", "Atk", baseMaxAtk)

					elite1Atk, _ := strconv.Atoi(e.ChildText("td:nth-of-type(3)"))
					updateOrCreateAttribute(operator.Attributes, "Elite1", "Atk", elite1Atk)

					elite2Atk, _ := strconv.Atoi(e.ChildText("td:nth-of-type(4)"))
					updateOrCreateAttribute(operator.Attributes, "Elite2", "Atk", elite2Atk)

					trustBonusAtk, _ := strconv.Atoi(e.ChildText("td:nth-of-type(5)"))
					operator.TrustBonus.Atk = trustBonusAtk
				case "DEF":
					baseDef, _ := strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
					updateOrCreateAttribute(operator.Attributes, "Base", "Def", baseDef)

					baseMaxDef, _ := strconv.Atoi(e.ChildText("td:nth-of-type(2)"))
					updateOrCreateAttribute(operator.Attributes, "BaseMax", "Def", baseMaxDef)

					elite1Def, _ := strconv.Atoi(e.ChildText("td:nth-of-type(3)"))
					updateOrCreateAttribute(operator.Attributes, "Elite1", "Def", elite1Def)

					elite2Def, _ := strconv.Atoi(e.ChildText("td:nth-of-type(4)"))
					updateOrCreateAttribute(operator.Attributes, "Elite2", "Def", elite2Def)

					trustBonusDEF, _ := strconv.Atoi(e.ChildText("td:nth-of-type(5)"))
					operator.TrustBonus.DEF = trustBonusDEF
				case "RES":
					baseRes, _ := strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
					updateOrCreateAttribute(operator.Attributes, "Base", "Res", baseRes)
					updateOrCreateAttribute(operator.Attributes, "BaseMax", "Res", baseRes)
					updateOrCreateAttribute(operator.Attributes, "Elite1", "Res", baseRes)
					updateOrCreateAttribute(operator.Attributes, "Elite2", "Res", baseRes)
				case "Redeployment time":
					baseRedeploymentTime := e.ChildText("td:nth-of-type(1)")
					updateOrCreateAttribute(operator.Attributes, "Base", "RedeploymentTime", baseRedeploymentTime)
					updateOrCreateAttribute(operator.Attributes, "BaseMax", "RedeploymentTime", baseRedeploymentTime)
					updateOrCreateAttribute(operator.Attributes, "Elite1", "RedeploymentTime", baseRedeploymentTime)
					updateOrCreateAttribute(operator.Attributes, "Elite2", "RedeploymentTime", baseRedeploymentTime)
				case "DP cost":
					if e.DOM.Find("td").Length() == 3 && e.ChildAttr("td:nth-of-type(1)", "colspan") == "2" {
						baseDPcost, _ := strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
						elite1Cost, _ := strconv.Atoi(e.ChildText("td:nth-of-type(2)"))
						elite2Cost, _ := strconv.Atoi(e.ChildText("td:nth-of-type(3)"))
						updateOrCreateAttribute(operator.Attributes, "Base", "DPcost", baseDPcost)
						updateOrCreateAttribute(operator.Attributes, "BaseMax", "DPcost", baseDPcost)
						updateOrCreateAttribute(operator.Attributes, "Elite1", "DPcost", elite1Cost)
						updateOrCreateAttribute(operator.Attributes, "Elite2", "DPcost", elite2Cost)
					} else {
						baseDPcost, _ := strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
						updateOrCreateAttribute(operator.Attributes, "Base", "DPcost", baseDPcost)
						updateOrCreateAttribute(operator.Attributes, "BaseMax", "DPcost", baseDPcost)
						baseMaxDPcost, _ := strconv.Atoi(e.ChildText("td:nth-of-type(2)"))
						updateOrCreateAttribute(operator.Attributes, "Elite1", "DPcost", baseMaxDPcost)
						updateOrCreateAttribute(operator.Attributes, "Elite2", "DPcost", baseMaxDPcost)
					}
				case "Block count":
					if e.DOM.Find("td").Length() == 2 && e.ChildAttr("td:nth-of-type(1)", "colspan") == "2" {
						baseBlock, _ := strconv.Atoi(e.ChildText("td:nth-of-type(1)"))
						eliteBlock, _ := strconv.Atoi(e.ChildText("td:nth-of-type(2)"))
						updateOrCreateAttribute(operator.Attributes, "Base", "Block", baseBlock)
						updateOrCreateAttribute(operator.Attributes, "BaseMax", "Block", baseBlock)
						updateOrCreateAttribute(operator.Attributes, "Elite1", "Block", eliteBlock)
						updateOrCreateAttribute(operator.Attributes, "Elite2", "Block", eliteBlock)
					} else {
						baseBlock := e.ChildText("td:nth-of-type(1)")
						updateOrCreateAttribute(operator.Attributes, "Base", "Block", baseBlock)
						updateOrCreateAttribute(operator.Attributes, "BaseMax", "Block", baseBlock)
						updateOrCreateAttribute(operator.Attributes, "Elite1", "Block", baseBlock)
						updateOrCreateAttribute(operator.Attributes, "Elite2", "Block", baseBlock)
					}
				case "Attack interval":
					baseAttackInterval := e.ChildText("td:nth-of-type(1)")
					updateOrCreateAttribute(operator.Attributes, "Base", "AttackInterval", baseAttackInterval)
					updateOrCreateAttribute(operator.Attributes, "BaseMax", "AttackInterval", baseAttackInterval)
					updateOrCreateAttribute(operator.Attributes, "Elite1", "AttackInterval", baseAttackInterval)
					updateOrCreateAttribute(operator.Attributes, "Elite2", "AttackInterval", baseAttackInterval)
				}
			})

			// Extract operator promotion
			e.ForEach("table#operator-promotion-table tr", func(i int, e *colly.HTMLElement) {
				promotion := model.Promotion{}
				promotion.Level = e.ChildText("th")

				e.ForEach("tbody td", func(i int, e *colly.HTMLElement) {
					e.ForEach("ul li", func(i int, e *colly.HTMLElement) {
						promotion.GainedEffect = append(promotion.GainedEffect, e.Text)
					})
					e.ForEach("div[style*='display:inline-block']", func(i int, e *colly.HTMLElement) {
						promotion.RequiredMaterials = append(promotion.RequiredMaterials, model.RequiredMaterial{
							Name:   e.ChildAttr("div.item-tooltip", "data-name"),
							Amount: e.ChildText("div.quantity"),
						})
					})
				})
				operator.Promotions = append(operator.Promotions, promotion)
			})

			// Extract operator talent
			e.ForEach("table#operator-talent-table", func(i int, e *colly.HTMLElement) {
				talent := model.Talent{}
				talent.Name = e.ChildText("th")

				// remove "Additional Information" from the talent.Name ("Auxiliary Equipment\nAdditional information")
				if strings.Contains(talent.Name, "\n") {
					talent.Name = strings.Split(talent.Name, "\n")[0]
				}

				e.ForEach("tr", func(i int, e *colly.HTMLElement) {
					// Skip the header row
					if e.ChildText("th") != "" {
						return
					}

					// if the text has words "Additional Information" then skip
					if strings.Contains(e.ChildText("td"), "Additional Information") {
						return
					}

					if e.DOM.Find("ul").Length() > 0 {
						e.ForEach("td ul li", func(i int, e *colly.HTMLElement) {
							talent.AdditionalInfo = append(talent.AdditionalInfo, e.Text)
						})
						return
					}

					// Extract requirement (text inside the first <td>)
					requirement := e.ChildAttr("td span", "title")
					// Extract effect (text inside the second <td>)
					effect := strings.TrimSpace(e.ChildText("td:nth-child(2)"))

					talent.Effect = append(talent.Effect, model.TalentEffect{
						Requirement: requirement,
						Effect:      effect,
					})
				})
				operator.Talents = append(operator.Talents, talent)
			})

			// Extract operator skills
			e.ForEach("div.mw-collapsible[data-expandtext='Click to show details'][data-collapsetext='Click to hide details']", func(i int, e *colly.HTMLElement) {

				skill := model.Skill{}

				e.ForEach("table.mrfz-wtable.skill-info-block", func(i int, e *colly.HTMLElement) {
					// Extract skill name
					skill.Name = strings.TrimSpace(e.ChildText("th div[style='font-size:14px;']"))

					// Extract Recovery Type (Auto or Manual)
					e.ForEach("th div[style*='background']", func(_ int, elem *colly.HTMLElement) {
						if strings.Contains(elem.Attr("style"), "#8EC31F") { // Auto Recovery
							skill.RecoveryType = "Auto Recovery"
						} else if strings.Contains(elem.Attr("style"), "#808080") { // Manual Activation
							skill.RecoveryType = "Manual"
						}
					})

					// Extract Charge Time
					skill.ChargeTime = strings.TrimSpace(e.ChildText("th div[style*='background:lightgray']"))
				})

				e.ForEach("table.mrfz-wtable tbody", func(i int, e *colly.HTMLElement) {

					if e.DOM.HasClass("skill-info-block") {
						return
					}

					e.ForEach("tr", func(i int, e *colly.HTMLElement) {
						if i == 0 {
							return // Skip header row
						}

						if e.DOM.Find("td").Length() == 1 && e.DOM.Find("th").Length() == 0 {
							if e.DOM.Find("ul").Length() > 0 {
								e.ForEach("td ul li", func(i int, e *colly.HTMLElement) {
									skill.Description = append(skill.Description, e.Text)
								})
							}
						}

						levelText := e.ChildAttr("th span", "title")

						// Skipping Skill Description
						if levelText == "" {
							return
						}

						description := e.ChildText("td:nth-child(2)")
						spCost, _ := strconv.Atoi(strings.TrimSpace(e.ChildText("td:nth-child(3)")))
						energyCost, _ := strconv.Atoi(strings.TrimSpace(e.ChildText("td:nth-child(4)")))
						coolDown, _ := strconv.Atoi(strings.TrimSpace(e.ChildText("td:nth-child(5)")))

						skill.Levels = append(skill.Levels, model.SkillLevel{
							Level:       levelText,
							Description: description,
							SPCost:      spCost,
							EnergyCost:  energyCost,
							CoolDown:    coolDown,
						})
					})
				})

				operator.Skills = append(operator.Skills, skill)
			})

			// Append operator to operators
			operators = append(operators, operator)

		})

		c.Visit(link)
	}()

	wg.Wait()
	return operators, nil
}
