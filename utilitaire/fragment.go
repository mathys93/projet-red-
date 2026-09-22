package utilitaire

import (
	"fmt"
	"math/rand"
)

type Equipment struct {
	Head  string
	Chest string
	Feet  string
}

type Character struct {
	Name         string
	Class        string
	Level        int
	MaxHP        int
	CurrentHP    int
	Inventory    []string
	InventoryMax int
	Fragment     int
	FragmentMax  int
	Skill        []string
	Equipment    Equipment

	SkipOpponentNextTurn bool
	StunnedTurns         int
	HasResurrectCharm    bool
	DamageBoost          int
}

type Monster struct {
	Name        string
	MaxHP       int
	CurrentHP   int
	Attack      int
	IsBoss      bool
	FragmentMin int
	FragmentMax int
}

func initCharacter(name, class string, level, maxHP, currentHP int, inventory []string) Character {
	return Character{
		Name:         name,
		Class:        class,
		Level:        level,
		MaxHP:        maxHP,
		CurrentHP:    currentHP,
		Inventory:    inventory,
		InventoryMax: 10,
		Fragment:     100,
		FragmentMax:  200,
		Skill:        []string{"Coup de poing"},
	}
}

func initGoblin() Monster {
	return Monster{
		Name:        "Gobelin d'entrainement",
		MaxHP:       40,
		CurrentHP:   40,
		Attack:      5,
		IsBoss:      false,
		FragmentMin: 2,
		FragmentMax: 5,
	}
}

func initBoss(name string, hp, attack int) Monster {
	return Monster{
		Name:        name,
		MaxHP:       hp,
		CurrentHP:   hp,
		Attack:      attack,
		IsBoss:      true,
		FragmentMin: 20,
		FragmentMax: 40,
	}
}

func addInventory(c *Character, item string) bool {
	if len(c.Inventory) >= c.InventoryMax {
		fmt.Println("Ton inventaire est plein !")
		return false
	}
	c.Inventory = append(c.Inventory, item)
	return true
}

func removeInventory(c *Character, item string) bool {
	for i, it := range c.Inventory {
		if it == item {
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			return true
		}
	}
	return false
}

func accessInventory(c *Character) {
	if len(c.Inventory) == 0 {
		fmt.Println("Ton inventaire est vide.")
		return
	}
	fmt.Println("\n--- INVENTAIRE ---")
	for i, item := range c.Inventory {
		fmt.Printf("%d. %s\n", i+1, item)
	}
}

func giveFragments(c *Character, montant int) {
	c.Fragment += montant
	if c.Fragment > c.FragmentMax {
		fmt.Printf("Tu as atteint la limite de Fragments (%d), le surplus est perdu.\n", c.FragmentMax)
		c.Fragment = c.FragmentMax
	}
}

func dropFragments(c *Character, m Monster) {
	montant := m.FragmentMin + rand.Intn(m.FragmentMax-m.FragmentMin+1)

	if m.IsBoss {
		fmt.Printf("%s s'effondre ! Tu obtiens %d Fragments (butin de boss).\n", m.Name, montant)
	} else {
		fmt.Printf("%s est vaincu ! Tu obtiens %d Fragments.\n", m.Name, montant)
	}

	giveFragments(c, montant)
	fmt.Printf("Fragments : %d / %d\n", c.Fragment, c.FragmentMax)
}

func takePot(c *Character) {
	if !removeInventory(c, "Disque de Pucci") {
		fmt.Println("Vous n'avez pas de Disque de Pucci !")
		return
	}
	c.CurrentHP += 50
	if c.CurrentHP > c.MaxHP {
		c.CurrentHP = c.MaxHP
	}
	fmt.Printf("Vous utilisez Disque de Pucci. PV : %d / %d\n", c.CurrentHP, c.MaxHP)
}

func poisonPot(c *Character) {
	if !removeInventory(c, "Tel Diavolo") {
		fmt.Println("Vous n'avez pas de Tel Diavolo !")
		return
	}
	fmt.Println("Tu ressens le poison du temps se refermer sur toi...")
	for i := 0; i < 3; i++ {
		c.CurrentHP -= 10
		if c.CurrentHP < 0 {
			c.CurrentHP = 0
		}
		fmt.Printf("PV : %d / %d\n", c.CurrentHP, c.MaxHP)
	}
}

func contientSort(c *Character, sort string) bool {
	for _, s := range c.Skill {
		if s == sort {
			return true
		}
	}
	return false
}

func spellBook(c *Character) {
	if contientSort(c, "Onde Solaire") {
		fmt.Println("Tu connais déjà ce sort !")
		return
	}
	c.Skill = append(c.Skill, "Onde Solaire")
	fmt.Println("Tu as appris : Onde Solaire !")
}

var marchandObjets = map[string]int{
	"Disque de Pucci":              3,
	"Tel Diavolo":                  6,
	"Arrow":                        15,
	"Livre de Sort : Onde Solaire": 25,
	"Augmentation d'inventaire":    20,
}

func acheterMarchand(c *Character, item string) {
	prix, existe := marchandObjets[item]
	if !existe {
		fmt.Println("Le marchand ne vend pas cet objet.")
		return
	}

	if c.Fragment < prix {
		fmt.Printf("Tu n'as pas assez de Fragments pour acheter %s.\n", item)
		fmt.Printf("Il faut %d Fragments, tu n'en as que %d.\n", prix, c.Fragment)
		return
	}

	if item == "Augmentation d'inventaire" {
		if !upgradeInventorySlot(c) {
			return
		}
		c.Fragment -= prix
		fmt.Printf("Fragments restants : %d\n", c.Fragment)
		return
	}

	if !addInventory(c, item) {
		return
	}

	c.Fragment -= prix
	fmt.Printf("Tu as acheté %s !\n", item)
	fmt.Printf("Fragments restants : %d\n", c.Fragment)

	if item == "Livre de Sort : Onde Solaire" {
		removeInventory(c, item)
		spellBook(c)
	}
}

func menuMarchand(c *Character) {
	noms := []string{
		"Disque de Pucci", "Tel Diavolo", "Arrow",
		"Livre de Sort : Onde Solaire", "Augmentation d'inventaire",
	}

	for {
		fmt.Println("\n--- MARCHAND ---")
		for i, nom := range noms {
			fmt.Printf("%d. %s - %d Fragments\n", i+1, nom, marchandObjets[nom])
		}
		fmt.Println("0. Retour")

		var choix int
		fmt.Print("Ton choix : ")
		fmt.Scanln(&choix)

		if choix == 0 {
			return
		}
		if choix < 1 || choix > len(noms) {
			fmt.Println("Choix invalide.")
			continue
		}
		acheterMarchand(c, noms[choix-1])
	}
}

var forgeronObjets = map[string]int{
	"Chapeau de Gyro Zeppeli": 8,
	"Haut de Giorno":          10,
	"Bas de Giorno":           8,
}

func acheterForgeron(c *Character, item string) {
	prix, existe := forgeronObjets[item]
	if !existe {
		fmt.Println("Le forgeron ne vend pas cet objet.")
		return
	}

	if c.Fragment < prix {
		fmt.Printf("Tu n'as pas assez de Fragments pour acheter %s.\n", item)
		fmt.Printf("Il faut %d Fragments, tu n'en as que %d.\n", prix, c.Fragment)
		return
	}

	if !addInventory(c, item) {
		return
	}

	c.Fragment -= prix
	fmt.Printf("Tu as acheté %s !\n", item)
	fmt.Printf("Fragments restants : %d\n", c.Fragment)
}

func menuForgeron(c *Character) {
	noms := []string{"Chapeau de Gyro Zeppeli", "Haut de Giorno", "Bas de Giorno"}

	for {
		fmt.Println("\n--- FORGERON ---")
		for i, nom := range noms {
			fmt.Printf("%d. %s - %d Fragments\n", i+1, nom, forgeronObjets[nom])
		}
		fmt.Println("0. Retour")

		var choix int
		fmt.Print("Ton choix : ")
		fmt.Scanln(&choix)

		if choix == 0 {
			return
		}
		if choix < 1 || choix > len(noms) {
			fmt.Println("Choix invalide.")
			continue
		}
		acheterForgeron(c, noms[choix-1])
	}
}

func bonusHPPourItem(item string) int {
	switch item {
	case "Chapeau de Gyro Zeppeli":
		return 10
	case "Haut de Giorno":
		return 25
	case "Bas de Giorno":
		return 15
	}
	return 0
}

func equiperArmure(c *Character, item string) {
	var slot *string

	switch item {
	case "Chapeau de Gyro Zeppeli":
		slot = &c.Equipment.Head
	case "Haut de Giorno":
		slot = &c.Equipment.Chest
	case "Bas de Giorno":
		slot = &c.Equipment.Feet
	default:
		fmt.Println("Cet objet ne peut pas être équipé.")
		return
	}

	if !removeInventory(c, item) {
		fmt.Println("Tu ne possèdes pas cet objet.")
		return
	}

	if *slot != "" {
		ancien := *slot
		addInventory(c, ancien)
		c.MaxHP -= bonusHPPourItem(ancien)
		fmt.Printf("Tu remplaces %s par %s.\n", ancien, item)
	}

	*slot = item
	c.MaxHP += bonusHPPourItem(item)
	fmt.Printf("Tu équipes %s ! PV max : %d\n", item, c.MaxHP)
}

func calculateDamageReduction(c *Character) float64 {
	reduction := 0.0
	if c.Equipment.Head != "" {
		reduction += 0.25
	}
	if c.Equipment.Chest != "" {
		reduction += 0.25
	}
	if c.Equipment.Feet != "" {
		reduction += 0.25
	}
	if reduction > 0.75 {
		reduction = 0.75
	}
	return reduction
}

func applyDamageToCharacter(c *Character, degats int) int {
	reduction := calculateDamageReduction(c)
	degatsReels := int(float64(degats) * (1 - reduction))

	c.CurrentHP -= degatsReels
	if c.CurrentHP < 0 {
		c.CurrentHP = 0
	}

	if reduction > 0 {
		fmt.Printf("(Armure : dégâts réduits de %.0f%%, %d dégâts au lieu de %d)\n",
			reduction*100, degatsReels, degats)
	}

	return degatsReels
}

var upgradesUtilisees int

func upgradeInventorySlot(c *Character) bool {
	if upgradesUtilisees >= 3 {
		fmt.Println("Tu as déjà utilisé les 3 augmentations d'inventaire disponibles.")
		return false
	}
	c.InventoryMax += 10
	upgradesUtilisees++
	fmt.Printf("Capacité d'inventaire augmentée ! Nouvelle limite : %d\n", c.InventoryMax)
	return true
}

type StandEffect struct {
	Name  string
	Type  string
	Apply func(c *Character)
}

var rouletteStand = []StandEffect{
	{
		Name: "Crazy Diamond : soin complet",
		Type: "bonus",
		Apply: func(c *Character) {
			c.CurrentHP = c.MaxHP
		},
	},
	{
		Name: "Star Platinum : +20 PV max",
		Type: "bonus",
		Apply: func(c *Character) {
			c.MaxHP += 20
			c.CurrentHP += 20
		},
	},
	{
		Name: "The World : temps arrêté (l'adversaire passe son tour, +5 dégâts)",
		Type: "bonus",
		Apply: func(c *Character) {
			c.SkipOpponentNextTurn = true
			c.DamageBoost += 5
			fmt.Println("Za Warudo ! L'adversaire ne pourra pas jouer au prochain tour.")
		},
	},
	{
		Name: "King Crimson : le futur est effacé (résurrection assurée)",
		Type: "bonus",
		Apply: func(c *Character) {
			c.HasResurrectCharm = true
			fmt.Println("Si tu tombes à 0 PV, tu reviendras automatiquement une fois.")
		},
	},
	{
		Name: "Gold Experience : +10 dégâts et +15 Fragments",
		Type: "bonus",
		Apply: func(c *Character) {
			c.DamageBoost += 10
			giveFragments(c, 15)
		},
	},
	{
		Name: "Whitesnake : un item volé",
		Type: "malus",
		Apply: func(c *Character) {
			if len(c.Inventory) > 0 {
				perdu := c.Inventory[0]
				removeInventory(c, perdu)
				fmt.Printf("Tu perds : %s\n", perdu)
			} else {
				fmt.Println("Ton inventaire était déjà vide, rien à voler.")
			}
		},
	},
	{
		Name: "Hermit Purple : enraciné, 2 tours sans agir",
		Type: "malus",
		Apply: func(c *Character) {
			c.StunnedTurns += 2
			fmt.Println("Des ronces t'immobilisent pendant 2 tours !")
		},
	},
}

func utiliserArrow(c *Character) {
	if !removeInventory(c, "Arrow") {
		fmt.Println("Tu n'as pas d'Arrow.")
		return
	}

	tirage := rouletteStand[rand.Intn(len(rouletteStand))]
	fmt.Println("\n*** La flèche te transperce... ***")
	fmt.Printf("Stand éveillé : %s [%s]\n", tirage.Name, tirage.Type)
	tirage.Apply(c)

	fmt.Printf("PV : %d / %d | Fragments : %d\n", c.CurrentHP, c.MaxHP, c.Fragment)
}    
