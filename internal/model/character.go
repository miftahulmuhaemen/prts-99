package model

type Attribute struct {
	HP               int
	Atk              int
	Def              int
	Res              int
	RedeploymentTime string
	DPcost           int
	Block            int
	AttackInterval   string
}

type TrustExtraStatus struct {
	HP  int
	Atk int
	DEF int
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
	Level  string
	Effect string
}

type RequiredMaterial struct {
	Name   string
	Amount string
}

type Promotion struct {
	Level             string
	GainedEffect      []string
	RequiredMaterials []RequiredMaterial
}

type TalentEffect struct {
	Requirement string
	Effect      string
}

type Talent struct {
	Name           string
	Effect         []TalentEffect
	AdditionalInfo []string
}

type SkillLevel struct {
	Level       string `json:"Level"`       // Skill Level
	Description string `json:"Description"` // Description of the skill at this level
	SPCost      int    `json:"SPCost"`      // SP Cost (Play)
	EnergyCost  int    `json:"EnergyCost"`  // Energy Cost
	CoolDown    int    `json:"CoolDown"`    // Cooldown Time
}

type Skill struct {
	Name         string       `json:"Name"`         // Name of the skill
	RecoveryType string       `json:"RecoveryType"` // Auto Recovery or Manual
	ChargeTime   string       `json:"ChargeTime"`   // Skill Charge Time
	Levels       []SkillLevel `json:"Levels"`       // Array of Skill Levels
	Description  []string     `json:"Description"`  // Description of the skill
}

type Operator struct {
	OperatorName     string
	ShortDescription string
	AscensionWords   []string
	Class            string
	Branch           string
	Faction          string
	Position         string
	Tags             []string
	Trait            string
	CharacterInfo    CharacterInfo
	Attributes       map[string]Attribute
	TrustBonus       TrustExtraStatus
	Potentials       []Potential
	Promotions       []Promotion
	Talents          []Talent
	Skills           []Skill
}
