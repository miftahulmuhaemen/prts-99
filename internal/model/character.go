package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qdrant/go-client/qdrant"
)

type Attribute struct {
	HP               int    `json:"hp"`
	Atk              int    `json:"atk"`
	Def              int    `json:"def"`
	Res              int    `json:"res"`
	RedeploymentTime string `json:"redeployment_time"`
	DPcost           int    `json:"dp_cost"`
	Block            int    `json:"block"`
	AttackInterval   string `json:"attack_interval"`
}

type TrustExtraStatus struct {
	HP  int `json:"hp"`
	Atk int `json:"atk"`
	DEF int `json:"def"`
}

type CharacterInfo struct {
	ThemeSong                 string `json:"themesong"`
	InternalName              string `json:"internalname"`
	BasedOn                   string `json:"basedon"`
	Etymology                 string `json:"etymology"`
	FileNo                    string `json:"fileno"`
	OperatorRecord            string `json:"operatorrecord"`
	ParadoxSimulation         string `json:"simulation"`
	Illustrator               string `json:"illustrator"`
	JapaneseCV                string `json:"japanesecv"`
	MandarinCV                string `json:"mandarincv"`
	CantoneseCV               string `json:"cantonesecv"`
	EnglishCV                 string `json:"englishcv"`
	KoreanCV                  string `json:"koreancv"`
	Gender                    string `json:"gender"`
	CombatExperience          string `json:"combatexperience"`
	PlaceOfBirth              string `json:"placeofbirth"`
	DateOfBirth               string `json:"dateofbirth"`
	Race                      string `json:"race"`
	Height                    string `json:"height"`
	InfectionStatus           string `json:"infectionstatus"`
	PhysicalStrength          string `json:"physicalstrength"`
	Mobility                  string `json:"mobility"`
	PhysicalResilience        string `json:"physicalresilience"`
	TacticalAcumen            string `json:"tacticalacumen"`
	CombatSkill               string `json:"combatskill"`
	OriginiumArtsAssimilation string `json:"originiumartsassimilation"`
}

type Potential struct {
	Level  string `json:"level"`
	Effect string `json:"effect"`
}

type RequiredMaterial struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type Promotion struct {
	Level             string             `json:"level"`
	GainedEffect      []string           `json:"gained_effect"`
	RequiredMaterials []RequiredMaterial `json:"required_materials"`
}

type TalentEffect struct {
	Requirement string `json:"requirement"`
	Effect      string `json:"effect"`
}

type Talent struct {
	Name           string         `json:"name"`
	Effect         []TalentEffect `json:"effect"`
	AdditionalInfo []string       `json:"additional_info"`
}

type SkillLevel struct {
	Level       string `json:"level"`       // Skill Level
	Description string `json:"description"` // Description of the skill at this level
	SPCost      int    `json:"sp_cost"`     // SP Cost (Play)
	EnergyCost  int    `json:"energy_cost"` // Energy Cost
	CoolDown    int    `json:"cool_down"`   // Cooldown Time
}

type Skill struct {
	Name         string       `json:"name"`          // Name of the skill
	RecoveryType string       `json:"recovery_type"` // Auto Recovery or Manual
	ChargeTime   string       `json:"charge_time"`   // Skill Charge Time
	Levels       []SkillLevel `json:"levels"`        // Array of Skill Levels
	Description  []string     `json:"description"`   // Description of the skill
}

type Operator struct {
	OperatorName     string               `json:"operator_name"`
	ShortDescription string               `json:"short_description"`
	AscensionWords   []string             `json:"ascension_words"`
	Class            string               `json:"class"`
	Branch           string               `json:"branch"`
	Faction          string               `json:"faction"`
	Position         string               `json:"position"`
	Tags             []string             `json:"tags"`
	Trait            string               `json:"trait"`
	CharacterInfo    CharacterInfo        `json:"character_info"`
	Attributes       map[string]Attribute `json:"attributes"`
	TrustBonus       TrustExtraStatus     `json:"trust_bonus"`
	Potentials       []Potential          `json:"potentials"`
	Promotions       []Promotion          `json:"promotions"`
	Talents          []Talent             `json:"talents"`
	Skills           []Skill              `json:"skills"`
}

// // Create function to convert struct to JSON string
func JSONToString(data interface{}) string {
	// Convert struct to JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func (o Operator) OperatorToStrings() ([][]string, []map[string]*qdrant.Value) {

	retStr := [][]string{}
	retVal := []map[string]*qdrant.Value{}
	metadata := fmt.Sprintf(`%s`, o.OperatorName)

	base := fmt.Sprintf(`
		ShortDescription: %s
		AscensionWords: %s
		Class: %s
		Branch: %s
		Faction: %s
		Position: %s
		Tags: %s
		Trait: %s
		Theme Song: %s
		Internal Name: %s
		Based On: %s
		Etymology: %s
		File No: %s
		Operator Record: %s
		Paradox Simulation: %s
		Illustrator: %s
		Japanese CV: %s
		Mandarin CV: %s
		Cantonese CV: %s
		English CV: %s
		Korean CV: %s
		Gender: %s
		Combat Experience: %s
		Place of Birth: %s
		Date of Birth: %s
		Race: %s
		Height: %s
		Infection Status: %s
		Physical Strength: %s
		Mobility: %s
		Physical Resilience: %s
		Tactical Acumen: %s
		Combat Skill: %s
		Originium Arts Assimilation: %s
		`, o.ShortDescription, strings.Join(o.AscensionWords, ", "), o.Class, o.Branch, o.Faction, o.Position, strings.Join(o.Tags, ", "), o.Trait,
		o.CharacterInfo.ThemeSong, o.CharacterInfo.InternalName, o.CharacterInfo.BasedOn, o.CharacterInfo.Etymology, o.CharacterInfo.FileNo, o.CharacterInfo.OperatorRecord, o.CharacterInfo.ParadoxSimulation, o.CharacterInfo.Illustrator,
		o.CharacterInfo.JapaneseCV, o.CharacterInfo.MandarinCV, o.CharacterInfo.CantoneseCV, o.CharacterInfo.EnglishCV, o.CharacterInfo.KoreanCV, o.CharacterInfo.Gender, o.CharacterInfo.CombatExperience, o.CharacterInfo.PlaceOfBirth, o.CharacterInfo.DateOfBirth,
		o.CharacterInfo.Race, o.CharacterInfo.Height, o.CharacterInfo.InfectionStatus, o.CharacterInfo.PhysicalStrength, o.CharacterInfo.Mobility, o.CharacterInfo.PhysicalResilience, o.CharacterInfo.TacticalAcumen, o.CharacterInfo.CombatSkill, o.CharacterInfo.OriginiumArtsAssimilation)
	retStr = append(retStr, []string{"Base", metadata, base})
	retVal = append(retVal, map[string]*qdrant.Value{
		"operator_name":               qdrant.NewValueString(o.OperatorName),
		"short_description":           qdrant.NewValueString(o.ShortDescription),
		"ascension_words":             qdrant.NewValueString(strings.Join(o.AscensionWords, ", ")),
		"class":                       qdrant.NewValueString(o.Class),
		"branch":                      qdrant.NewValueString(o.Branch),
		"faction":                     qdrant.NewValueString(o.Faction),
		"position":                    qdrant.NewValueString(o.Position),
		"tags":                        qdrant.NewValueString(strings.Join(o.Tags, ", ")),
		"trait":                       qdrant.NewValueString(o.Trait),
		"theme_song":                  qdrant.NewValueString(o.CharacterInfo.ThemeSong),
		"internal_name":               qdrant.NewValueString(o.CharacterInfo.InternalName),
		"based_on":                    qdrant.NewValueString(o.CharacterInfo.BasedOn),
		"etymology":                   qdrant.NewValueString(o.CharacterInfo.Etymology),
		"file_no":                     qdrant.NewValueString(o.CharacterInfo.FileNo),
		"operator_record":             qdrant.NewValueString(o.CharacterInfo.OperatorRecord),
		"paradox_simulation":          qdrant.NewValueString(o.CharacterInfo.ParadoxSimulation),
		"illustrator":                 qdrant.NewValueString(o.CharacterInfo.Illustrator),
		"japanese_cv":                 qdrant.NewValueString(o.CharacterInfo.JapaneseCV),
		"mandarin_cv":                 qdrant.NewValueString(o.CharacterInfo.MandarinCV),
		"cantonesecv":                 qdrant.NewValueString(o.CharacterInfo.CantoneseCV),
		"english_cv":                  qdrant.NewValueString(o.CharacterInfo.EnglishCV),
		"korean_cv":                   qdrant.NewValueString(o.CharacterInfo.KoreanCV),
		"gender":                      qdrant.NewValueString(o.CharacterInfo.Gender),
		"combat_experience":           qdrant.NewValueString(o.CharacterInfo.CombatExperience),
		"place_of_birth":              qdrant.NewValueString(o.CharacterInfo.PlaceOfBirth),
		"date_of_birth":               qdrant.NewValueString(o.CharacterInfo.DateOfBirth),
		"race":                        qdrant.NewValueString(o.CharacterInfo.Race),
		"height":                      qdrant.NewValueString(o.CharacterInfo.Height),
		"infection_status":            qdrant.NewValueString(o.CharacterInfo.InfectionStatus),
		"physical_strength":           qdrant.NewValueString(o.CharacterInfo.PhysicalStrength),
		"mobility":                    qdrant.NewValueString(o.CharacterInfo.Mobility),
		"physical_resilience":         qdrant.NewValueString(o.CharacterInfo.PhysicalResilience),
		"tactical_acumen":             qdrant.NewValueString(o.CharacterInfo.TacticalAcumen),
		"combatskill":                 qdrant.NewValueString(o.CharacterInfo.CombatSkill),
		"originium_arts_assimilation": qdrant.NewValueString(o.CharacterInfo.OriginiumArtsAssimilation),
	})

	for key, attribute := range o.Attributes {
		atbStr := fmt.Sprintf(`%s
			HP: %d
			Atk: %d
			Def: %d
			Res: %d
			Redeployment Time: %s
			DP Cost: %d
			Block: %d
			Attack Interval: %s
		`, key, attribute.HP, attribute.Atk, attribute.Def, attribute.Res, attribute.RedeploymentTime, attribute.DPcost, attribute.Block, attribute.AttackInterval)
		retStr = append(retStr, []string{key, metadata, atbStr})
		retVal = append(retVal, map[string]*qdrant.Value{
			"hp":                qdrant.NewValueString(fmt.Sprint(attribute.HP)),
			"atk":               qdrant.NewValueString(fmt.Sprint(attribute.Atk)),
			"def":               qdrant.NewValueString(fmt.Sprint(attribute.Def)),
			"res":               qdrant.NewValueString(fmt.Sprint(attribute.Res)),
			"redeployment_time": qdrant.NewValueString(attribute.RedeploymentTime),
			"dp_cost":           qdrant.NewValueString(fmt.Sprint(attribute.DPcost)),
			"block":             qdrant.NewValueString(fmt.Sprint(attribute.Block)),
			"attack_interval":   qdrant.NewValueString(attribute.AttackInterval),
		})
	}

	// potential := fmt.Sprintf(`Potential: %s`, JSONToString(o.Potentials))

	// promotion := fmt.Sprintf(`Promotion: %s`, JSONToString(o.Promotions))

	// talent := fmt.Sprintf(`Talent: %s`, JSONToString(o.Talents))

	for iter, skill := range o.Skills {
		for iterLevel, level := range skill.Levels {
			str := fmt.Sprintf(`Skill %d - Lvl %d: {
				Name: %s
				RecoveryType: %s
				ChargeTime:  %s
				Description:  %s
				Level: %s
				Description Level: %s
				SPCost: %d
				EnergyCost: %d
				CoolDown: %d
			}`, iter, iterLevel, skill.Name, skill.RecoveryType, skill.ChargeTime, strings.Join(skill.Description, ", "),
				level.Level, level.Description, level.SPCost, level.EnergyCost, level.CoolDown)
			retStr = append(retStr, []string{fmt.Sprintf("Skill %d", iterLevel), metadata, str})
			retVal = append(retVal, map[string]*qdrant.Value{
				"recovery_type":     qdrant.NewValueString(fmt.Sprint(skill.RecoveryType)),
				"charge_time":       qdrant.NewValueString(fmt.Sprint(skill.ChargeTime)),
				"description":       qdrant.NewValueString(strings.Join(skill.Description, ", ")),
				"level":             qdrant.NewValueString(fmt.Sprint(level.Level)),
				"description_level": qdrant.NewValueString(fmt.Sprint(level.Description)),
				"sp_cost":           qdrant.NewValueString(fmt.Sprint(level.SPCost)),
				"energy_cost":       qdrant.NewValueString(fmt.Sprint(level.EnergyCost)),
				"cool_down":         qdrant.NewValueString(fmt.Sprint(level.CoolDown)),
			})
		}
	}

	return retStr, retVal
}

// func (o Operator) OperatorToStrings() string {

// 	base := fmt.Sprintf(`%s
// 		%s;
// 		%s;
// 		%s;
// 		%s;
// 		%s;
// 		%s;
// 		%s;
// 		%s;`, o.OperatorName, o.ShortDescription, strings.Join(o.AscensionWords, "; "), o.Class, o.Branch, o.Faction, o.Position, strings.Join(o.Tags, ", "), o.Trait)

// 	// info := fmt.Sprintf(`%s;`, JSONToString(o.CharacterInfo))
// 	// attribute := fmt.Sprintf(`%s; %s;`, JSONToString(o.Attributes), JSONToString(o.TrustBonus))
// 	// potential := fmt.Sprintf(`%s;`, JSONToString(o.Potentials))
// 	// promotion := fmt.Sprintf(`%s;`, JSONToString(o.Promotions))
// 	// talent := fmt.Sprintf(`%s;`, JSONToString(o.Talents))

// 	// skills := []string{}
// 	// for iter, skill := range o.Skills {
// 	// 	for iterLevel, level := range skill.Levels {
// 	// 		str := fmt.Sprintf(`%d - Lvl %d: {
// 	// 			%s
// 	// 			%s
// 	// 			%s
// 	// 			%s
// 	// 			%s
// 	// 		}`, iter, iterLevel, skill.Name, skill.RecoveryType, skill.ChargeTime, skill.Description,
// 	// 			JSONToString(level))
// 	// 		skills = append(skills, str)
// 	// 	}
// 	// }

// 	return base
// }

// func (o Operator) OperatorToQdrantValue() map[string]*qdrant.Value {

// 	info := fmt.Sprintf(`Detail Information: %s`, JSONToString(o.CharacterInfo))
// 	attribute := fmt.Sprintf(`Attribute: %s; Status Bonus: %s`, JSONToString(o.Attributes), JSONToString(o.TrustBonus))
// 	potential := fmt.Sprintf(`Potential: %s: `, JSONToString(o.Potentials))
// 	promotion := fmt.Sprintf(`Promotion: %s: `, JSONToString(o.Promotions))
// 	talent := fmt.Sprintf(`Talent: %s: `, JSONToString(o.Talents))

// 	qdrantValue := map[string]any{
// 		"operator_name":     o.OperatorName,
// 		"short_description": o.ShortDescription,
// 		"ascension_words":   strings.Join(o.AscensionWords, "; "),
// 		"class":             o.Class,
// 		"branch":            o.Branch,
// 		"faction":           o.Faction,
// 		"position":          o.Position,
// 		"tags":              strings.Join(o.Tags, ", "),
// 		"trait":             o.Trait,
// 		"character_info":    info,
// 		"attributes":        attribute,
// 		"potentials":        potential,
// 		"promotions":        promotion,
// 		"talents":           talent,
// 	}

// 	for iter, skill := range o.Skills {
// 		for iterLevel, level := range skill.Levels {
// 			str := fmt.Sprintf(`
// 				Name: %s;
// 				RecoveryType: %s;
// 				ChargeTime:  %s;
// 				Description:  %s;
// 				Additional: %s;`, skill.Name, skill.RecoveryType, skill.ChargeTime, strings.Join(skill.Description, ", "),
// 				JSONToString(level))
// 			qdrantValue[fmt.Sprintf("Skill %d - Lvl %d", iter, iterLevel)] = str
// 		}
// 	}

// 	return qdrant.NewValueMap(qdrantValue)
// }
