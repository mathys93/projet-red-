# ProjetRED

Un RPG en ligne de commande façon **Undertale**, sur un thème **JoJo's Bizarre Adventure** : tu incarnes un personnage qui traverse une carte de zones gardées par des boss (DIO, Diavolo, Yoshikage Kira, Enrico Pucci...), avec combats au tour par tour, mini-jeu d'esquive de projectiles, Stands à débloquer, boutique, et une interface texte redessinée à chaque frame façon boîte de dialogue Undertale.

Le jeu tourne dans un terminal. Le rendu avancé (esquive en temps réel, plein écran automatique, musique) est pensé pour **Windows** ; sur les autres OS, le jeu reste jouable mais bascule sur des solutions de repli plus simples (voir [Plateforme & terminal](#plateforme--terminal)).

## Sommaire

- [Démarrage rapide](#démarrage-rapide)
- [Boucle de jeu](#boucle-de-jeu)
- [Carte et zones](#carte-et-zones)
- [Combat](#combat)
- [Le personnage](#le-personnage)
- [Les boss](#les-boss)
- [Objets et économie](#objets-et-économie)
- [Musique](#musique)
- [Rendu ASCII des artworks](#rendu-ascii-des-artworks)
- [Plateforme & terminal](#plateforme--terminal)
- [Arborescence du projet](#arborescence-du-projet)
- [Points morts et choses à finir](#points-morts-et-choses-à-finir)
- [Contrôles](#contrôles)

## Démarrage rapide

Prérequis : Go (voir la version exacte dans [go.mod](go.mod)).

```
go run .          # lance le jeu directement
go build .        # produit un exécutable ProjetRED(.exe)
```

Le jeu s'exécute dans le terminal courant. Sur Windows, il agrandit automatiquement la fenêtre de la console au démarrage (voir `combat.MaximizeConsoleWindow`).

## Boucle de jeu

Tout part de [main.go](main.go) :

1. Création du joueur (`character.New`), avec deux Potions de vie en poche.
2. Construction de la carte (`world.New`).
3. Boucle principale : `combat.SelectZone` affiche la carte et renvoie la zone choisie (ou `nil` si le joueur quitte).
   - Zone boutique → `combat.RunShop`.
   - Zone de combat → `combat.RunBattle` contre le boss de la zone.
4. À la victoire (ou à l'épargne via MERCY) : la zone est marquée terminée (sauf zone "farm", voir plus bas), le joueur gagne des Fragments et de l'XP (avec messages de montée de niveau).
5. À la défaite : Game Over, le programme s'arrête.
6. Une fois les 3 zones de boss "normales" terminées, `world.AllCleared()` débloque le boss bonus (Enrico Pucci) ; `combat.ConfirmBonusBoss` demande confirmation avant de le lancer.

## Carte et zones

Définie dans [world/world.go](world/world.go). Une `Zone` a un nom, une description, éventuellement un `*boss.Boss`, et deux drapeaux : `IsShop` et `IsFarm`.

| Zone | Boss | Particularité |
|---|---|---|
| Zone 1 - Le Manoir de DIO | DIO | Toujours déverrouillée, c'est la première zone |
| Zone 2 - La Villa de Diavolo | Diavolo | Verrouillée tant que la Zone 1 n'est pas terminée |
| Zone 3 - Morioh | Yoshikage Kira | Verrouillée tant que la Zone 2 n'est pas terminée |
| Zone 4 - La Boutique | — | `IsShop` : toujours accessible, jamais "verrouillée" |
| Zone 5 - Le Terrain Vague | Iggy | `IsFarm` : toujours accessible, ses PV sont remis au max à chaque entrée, et elle n'est jamais marquée "terminée" — c'est la zone d'entraînement pour farmer XP/Fragments sans limite |
| (bonus, hors carte) | Enrico Pucci | Débloqué quand les zones 1/2/3 sont toutes terminées |

`World.Locked(zone)` implémente le verrouillage séquentiel : une zone de boss (hors boutique/farm) est verrouillée tant que la zone de boss précédente sur la carte n'est pas `Cleared`. `combat.SelectZone` affiche `[verrouillée]`, `[zone terminée]` ou `[toujours ouverte]` en conséquence et refuse d'entrer dans une zone verrouillée ou déjà nettoyée en réaffichant un message au lieu de lancer le combat.

## Combat

Le cœur du jeu, dans le package `combat` (essentiellement [combat/battle.go](combat/battle.go)).

### Mise en page à l'écran

Toute la mise en page est calculée dynamiquement en fonction de la taille réelle du terminal (`layout()` dans battle.go, basé sur `terminalSize()`), avec un repli sur 80×24 si la taille ne peut pas être lue. Chaque écran :

- affiche l'artwork du boss **au-dessus** du cadre de combat (`DrawArt`), redimensionné à la volée à la taille disponible (`ascii.Fit`) ;
- dessine un cadre (`DrawEmptyBox`) qui reste **vide** pendant les échanges de tours (façon boîte Undertale) ;
- réutilise ce même cadre comme **grille d'options à deux colonnes** (`DrawOptionGrid`) dès qu'un sous-menu est ouvert (FIGHT/ACT/ITEM, ou les menus de la boutique) — chaque option occupe une cellule (ex: "Potion de vie" et "Disque de Pucci" côte à côte), navigable au clavier en haut/bas/gauche/droite, avec la description de l'option survolée affichée sous le cadre ;
- affiche les 4 actions FIGHT/ACT/ITEM/MERCY comme de gros boutons encadrés et répartis sur toute la largeur du terminal (`DrawActionButtons`), plutôt qu'une ligne de texte compacte ;
- centre verticalement tout le contenu dans le terminal (`printPadding`) au lieu de le coller en haut de l'écran.

### Boucle d'un combat (`RunBattle`)

Pour chaque tour :

1. **Tour du joueur** : `RenderBattleScreen` affiche le menu FIGHT/ACT/ITEM/MERCY. Navigation avec les flèches (ou Z/S/Q/D), validation par Entrée, ou raccourci direct (1/F, 2/A, 3/I, 4/M).
   - **FIGHT** : choix d'une attaque (`Move`) dans la grille — coup de poing de base, plus l'attaque de Stand si un Stand est équipé. Coûte des MP, consomme les PV du boss.
   - **ACT** : sous-menu "lore" propre à chaque boss (Observer, Défier, mentionner un personnage...) — n'inflige pas de dégâts, sert au fan-service/texte.
   - **ITEM** : ouvre l'inventaire (objets identiques regroupés "Potion de vie x2") et applique l'effet de l'objet choisi.
   - **MERCY** : épargne le boss immédiatement, termine le combat en victoire "douce" (`ResultSpared`, traité comme une victoire côté récompenses).
   - Si le joueur est enraciné (`StunnedTurns > 0`, effet du Stand malus Hermit Purple), le tour est automatiquement sauté avec un message dédié.
2. **Tour du boss**, affiché à part (`RenderTurnMessage`) pour bien séparer visuellement le tour joueur/boss :
   - Sauf si `SkipBossNextTurn` est actif (effet de The World), le joueur passe d'abord par le **mini-jeu d'esquive** `RunDodge` (voir plus bas) ; s'il est touché, `b.Turn()` applique le pattern d'attaque du boss (dégâts, effets spéciaux) ; sinon un message "Tu esquives..." s'affiche sans dégâts.
   - Les changements de PV sont animés progressivement (`animateHPChange`), pas de saut instantané à la valeur finale.
3. Si le joueur tombe à 0 PV mais possède `HasResurrectCharm` (effet du Stand King Crimson), il ressuscite automatiquement avec un tiers de ses PV max au lieu de perdre (`tryResurrect`), à consommation unique.
4. La musique du boss démarre au début du combat et s'arrête à la fin (`audio.Play`/`audio.Stop`, voir [Musique](#musique)).

### Mini-jeu d'esquive (`combat/dodge.go`)

Pendant le tour du boss, un mini-jeu de type bullet-hell façon Undertale s'affiche dans le cadre de combat : le joueur déplace un cœur (`❤`) dans une arène avec les flèches, doit éviter les projectiles (`◆`/`◉`) pendant une durée fixe (`dodgeDuration` ticks). S'il est touché, le tour d'attaque normal du boss s'applique ; s'il esquive tout, le boss ne fait aucun dégât.

Le pattern des projectiles dépend du `boss.AttackStyle` du boss (`Boss.Style`), pensé pour coller au Stand de chacun :

| Style | Boss | Comportement des projectiles |
|---|---|---|
| `AttackStroll` (par défaut) | Iggy | Projectiles simples venant de la droite |
| `AttackTimeStop` | DIO | Le temps se fige périodiquement (vitesse à 0) puis explose en rafale à la reprise — référence à "The World" |
| `AttackErase` | Diavolo | Les projectiles deviennent invisibles par intermittence — référence à "King Crimson" qui efface le temps |
| `AttackBombs` | Kira | Bombes à retardement (`fuse`) qui explosent en étoile, plus des projectiles téléguidés (`homing`) — référence aux bombes de Killer Queen |
| `AttackAccelerate` | Pucci | La vitesse des projectiles augmente progressivement sur toute la durée — référence à "Made in Heaven" qui accélère le temps |

Le déplacement gère une accélération progressive tant qu'une touche est maintenue (lecture d'évènements clavier bas-niveau, voir [Plateforme & terminal](#plateforme--terminal)) ; sur un terminal qui ne permet pas de sondage clavier en temps réel (`canPollInput() == false`), le mini-jeu est court-circuité et compte comme une esquive automatique.

## Le personnage

[character/character.go](character/character.go) définit `Character`, utilisé aussi bien pour le joueur que pour les boss (le type `boss.Boss` l'embarque par composition).

Statistiques principales : `Level`, `AP` (Attack Power), `DP` (Defense Power, réduit les dégâts subis dans `TakeDamage`), `LP`/`MaxLP` (PV), `MP`/`MaxMP`, `XP`, `Fragment` (monnaie), `Inventory`, `Moves` (actions disponibles en FIGHT).

### Expérience et niveaux

`XPForLevel(level)` donne le seuil d'XP nécessaire pour passer au niveau suivant (`40 + (level-1)*30`, donc de plus en plus long). `GainXP(amount)` ajoute de l'XP et fait monter le niveau en chaîne si plusieurs paliers sont franchis d'un coup ; chaque niveau augmente `AP`, `DP`, `MaxLP`, `MaxMP`, puis restaure PV/MP au maximum, et renvoie un message par niveau gagné (affiché par `main.go` après chaque victoire).

### Les Stands ([character/stand.go](character/stand.go))

L'objet consommable **Arrow** (vendu chez le marchand) équipe un Stand aléatoire (`RollStand`, pondéré par rareté : 10% Légendaire, 30% Rare, 60% Commun). Équiper un Stand (`EquipStand`) :

- ajoute un bonus d'`AP` selon la rareté (`Rarity.Bonus()` : +2 Commun, +5 Rare, +9 Légendaire) ;
- remplace la liste de `Moves` du personnage par les moves par défaut + l'attaque signature du Stand ;
- si un Stand était déjà équipé, retire son bonus d'AP avant d'appliquer le nouveau (un seul Stand actif à la fois).

Stands disponibles (`standPool`) : Hermit Purple, Silver Chariot, Whitesnake (Communs) ; Crazy Diamond, Killer Queen, Sticky Fingers (Rares) ; Star Platinum, The World, King Crimson (Légendaires). Certains ont des effets spéciaux au-delà des dégâts bruts : Crazy Diamond soigne, Whitesnake vole du MP à la cible, Sticky Fingers vole soigne le lanceur, The World pose `SkipBossNextTurn`, King Crimson pose `HasResurrectCharm`.

## Les boss

[boss/boss.go](boss/boss.go) définit `Boss` (un `*character.Character` + artwork + `Pattern` de comportement + couleur ANSI + récompenses). Les patterns de combat (ce qui se passe au tour du boss et les options ACT) sont dans [boss/patterns.go](boss/patterns.go), un type par boss implémentant l'interface `Pattern` (`Act` + `ActOptions`).

| Boss | Zone | Niveau | AP | DP | PV | Récompense XP | Style d'attaque | Particularité de pattern |
|---|---|---|---|---|---|---|---|---|
| DIO | 1 | 5 | 10 | 4 | 70 | 50 | TimeStop | 1 coup sur 3 est "Za Warudo" : dégâts bruts (ignore la DP), impossible à esquiver via l'ACT normal |
| Diavolo | 2 | 8 | 14 | 6 | 110 | 90 | Erase | Se soigne périodiquement, "enrage" (dégâts +50%) sous 1/3 de PV |
| Yoshikage Kira | 3 | 11 | 18 | 8 | 150 | 140 | Bombs | "Sheer Heart Attack" (1 tour sur 4) inflige des dégâts bruts croissants avec la durée du combat |
| Enrico Pucci | Bonus | 15 | 24 | 10 | 220 | 260 | Accelerate | Se soigne, enrage sous 1/3 de PV, et inflige parfois des dégâts bruts ("Made in Heaven") |
| Iggy | 5 (farm) | 3 | 5 | 1 | 30 | 18 | Stroll | Boss faible et répétable, récompense en Fragments réduite (6 au lieu de 15 par défaut) |

`FragmentReward` vaut 15 par défaut pour tout boss créé via `boss.New`, sauf Iggy qui le réduit explicitement à 6 (cohérent avec le fait qu'il est combattable à l'infini).

## Objets et économie

Le package `utilitaire` fournit les catalogues et la logique d'achat/fabrication, utilisés par [combat/shop.go](combat/shop.go) :

- **Marchand** ([utilitaire/marchand.go](utilitaire/marchand.go)) : Potion de vie (4 Fragments, +20 PV), Potion de MP (4, +15 MP), Disque de Pucci (6, +50 PV), Tel Diavolo (3, dégâts sur soi ignorant la DP — objet piège/RP), Arrow (15, équipe un Stand aléatoire).
- **Forgeron** ([utilitaire/forgeron.go](utilitaire/forgeron.go)) : fabrique des équipements ("Chapeau/Tunique/Bottes de l'aventurier") contre des Fragments **et** des ressources (Plume de Corbeau, Cuir de Sanglier, Fourrure de Loup, Peau de Troll). ⚠️ voir [Points morts](#points-morts-et-choses-à-finir) : ces ressources ne sont actuellement obtenables nulle part dans le jeu.

Les effets des objets utilisables en combat (ITEM) sont centralisés dans [combat/items.go](combat/items.go) (`itemDescription`, `itemMenuItems`, `useItem`), pas dans le package `utilitaire` — c'est le code réellement utilisé par la boucle de jeu pour les items.

## Musique

[audio/tracks.go](audio/tracks.go) associe un nom de boss à un fichier `.mp3` dans `musique/` (`BossTracks`), avec repli sur `BattleTrack` (vide par défaut = silence). Sur Windows, la lecture passe par l'API MCI de `winmm.dll` (`mciSendStringW`, voir [audio/play_windows.go](audio/play_windows.go)) en boucle, avec volume réglable (`audio.Volume`, sur 1000). Sur les autres OS, `Play`/`Stop` sont des no-op ([audio/play_other.go](audio/play_other.go)).

`RunBattle` lance la musique du boss au début du combat et l'arrête à la fin.

## Rendu ASCII des artworks

[ascii/ascii.go](ascii/ascii.go) embarque les artworks (`go:embed`) des boss (`ascii/*.txt`) et expose `Fit(art, maxW, maxH)` pour les redimensionner à la volée à la taille réellement disponible à l'écran (recalculée à chaque frame, pas figée à la création du boss). Deux modes de rendu selon le contenu du fichier source :

- **rampe de caractères** (` .:-=+*#%@` ou blocs `░▒▓█`) : rééchantillonnage par densité de caractère, avec étirement de contraste par percentile ;
- **braille Unicode** (détecté automatiquement si le fichier contient des caractères U+2800–U+28FF) : chaque caractère braille encode une grille de 2×4 points, ce qui permet une résolution ASCII bien plus fine, avec tramage d'erreur (Floyd-Steinberg) façon image en niveaux de gris.

Deux outils CLI accompagnent ce pipeline (non utilisés par le jeu lui-même, juste pour produire/tester les artworks) :

- [tools/img2ascii](tools/img2ascii/main.go) : convertit une image (PNG/JPEG/WebP) en art ASCII ou braille, avec plein de réglages (contours, tons, contraste, tramage, recadrage, inversion, prévisualisation PNG). Exemple : `go run ./tools/img2ascii -in dio.png -out ascii/dio.txt -cols 70 -braille`.
- [tools/fitcheck](tools/fitcheck/main.go) : recharge un fichier `ascii/*.txt` existant, applique `ascii.Fit` à une taille donnée et affiche le résultat (+ une prévisualisation PNG optionnelle), pour vérifier rapidement le rendu sans lancer le jeu.

## Plateforme & terminal

Le comportement diffère entre Windows et les autres OS via des fichiers `_windows.go` / `_other.go` (build tags Go) :

| Fonctionnalité | Windows | Autres OS |
|---|---|---|
| Taille du terminal (`combat/termsize_*.go`) | Lue via `GetConsoleScreenBufferInfo` (API Win32) | Repli fixe sur 80×24 |
| Plein écran au démarrage | Fenêtre de la console maximisée (`ShowWindow`/`GetConsoleWindow`) | Aucun effet |
| Entrée clavier (`combat/input_*.go`) | Lecture bas niveau (`ReadConsoleInputW`), touches fléchées, et **sondage** d'évènements haut/bas pour le mini-jeu d'esquive en temps réel | Repli ligne par ligne (`bufio.Scanner` sur stdin : taper une lettre/mot + Entrée) ; le mini-jeu d'esquive est court-circuité (esquive automatique) |
| Musique (`audio/play_*.go`) | Lecture via `winmm.dll` | Silencieux (no-op) |

Le mode "entrée ligne par ligne" (`readLineKey`) reste toujours disponible comme repli si la console n'est pas détectée (utile aussi pour scripter/tester le jeu via un pipe stdin), avec les raccourcis Z/W (haut), S (bas), Q/A (gauche), D (droite), X (quitter), Entrée (valider).

## Arborescence du projet

```
main.go                    Point d'entrée, boucle de jeu haut niveau

character/
  character.go              Type Character (stats, PV/MP, inventaire), XP/niveaux
  stand.go                   Stands (rareté, pool, tirage, équipement)

boss/
  boss.go                    Type Boss, constructeurs des 5 boss, styles d'attaque
  patterns.go                Comportement de tour + options ACT, un type par boss

world/
  world.go                   Carte (Zone, World), verrouillage séquentiel des zones

combat/
  battle.go                  Rendu de l'écran de combat, boucle RunBattle, gros boutons/cadre
  dodge.go                   Mini-jeu d'esquive de projectiles (bullet-hell)
  submenu.go                 Sous-menu générique en grille (FIGHT/ACT/ITEM)
  items.go                   Effets des objets utilisables en combat
  arrow.go                   Effet de l'objet Arrow (équipe un Stand)
  mapmenu.go                 Écran de carte (sélection de zone) et confirmation boss bonus
  shop.go                    Écran boutique (marchand/forgeron)
  input.go                   Types clavier communs + repli "ligne par ligne"
  input_windows.go           Entrée clavier bas niveau (API console Win32)
  input_other.go             Repli non-Windows (pas de sondage temps réel)
  termsize_windows.go        Taille terminal + plein écran (API Win32)
  termsize_other.go          Repli non-Windows (taille fixe, pas de plein écran)

ascii/
  ascii.go                   Chargement des artworks (go:embed) + redimensionnement (Fit)
  *.txt                      Artworks des boss (rampe de caractères ou braille)

audio/
  tracks.go                  Association boss → fichier musique, résolution de chemin
  play_windows.go            Lecture MP3 via winmm.dll (MCI)
  play_other.go               No-op hors Windows

utilitaire/
  marchand.go                Catalogue + logique d'achat du marchand
  forgeron.go                 Catalogue + logique de fabrication (avec ressources)
  inventaire.go, disque_de_pucci.go, poison_pot.go   Code legacy, voir ci-dessous
  fragment.go                  Prototype legacy indépendant, voir ci-dessous
  (les effets des objets réellement utilisés en combat vivent dans combat/items.go, pas ici)

tools/
  img2ascii/main.go          CLI : image → art ASCII/braille
  fitcheck/main.go            CLI : recharge/redimensionne un art existant pour vérification
```

## Points morts et choses à finir

Pour une lecture honnête du code tel qu'il est aujourd'hui :

- **[utilitaire/fragment.go](utilitaire/fragment.go)** est un prototype autonome et plus ancien : il redéfinit son propre `Character`/`Monster` et toute une économie parallèle (équipement, sorts, roulette de Stand...) totalement déconnectée du vrai `character.Character` utilisé par le jeu. Aucune de ses fonctions n'est appelée ailleurs dans le projet — c'est du code mort, probablement un premier jet à nettoyer ou à supprimer.
- **[utilitaire/inventaire.go](utilitaire/inventaire.go)**, **[disque_de_pucci.go](utilitaire/disque_de_pucci.go)** et **[poison_pot.go](utilitaire/poison_pot.go)** définissent des fonctions (`PrintInventory`, `UseInventory`, `ApplyDisqueDePucci`, `BuyDisqueDePucci`, `BuyPoisonPotion`...) qui ne sont appelées nulle part : la logique réellement utilisée pour les objets en combat vit dans [combat/items.go](combat/items.go). Probablement une itération antérieure du menu ITEM, avant son passage à la grille en cadre.
- **Ressources du forgeron introuvables** : `ForgeronObjets` demande des ressources (Plume de Corbeau, Cuir de Sanglier, Fourrure de Loup, Peau de Troll) qu'aucun boss, aucun marchand et aucune autre source du jeu ne donne actuellement — la fabrication est donc bloquée en l'état.
- **Musiques non branchées** : `musique/` contient déjà `diavolo.mp3`, `kira.mp3`, `pucci.mp3` et `iggy.mp3`, mais `audio.BossTracks` ne les référence pas encore (seul DIO a une piste assignée, `KwikFlip.mp3`) — il suffit de compléter la map dans [audio/tracks.go](audio/tracks.go).
- **Fichiers orphelins à la racine** : `arrow.txt` (un art ASCII brut, non chargé par le package `ascii`, qui ne connaît que les fichiers dans `ascii/`) et `tel_diavolo` (fichier vide) semblent être des restes de travail à ranger ou supprimer.

## Contrôles

| Contexte | Touches |
|---|---|
| Navigation générale | Flèches directionnelles (ou Z/S/Q/D en repli ligne par ligne) |
| Valider | Entrée |
| Quitter / annuler | X (ou Échap en mode console Windows) |
| Menu FIGHT/ACT/ITEM/MERCY | Flèches + Entrée, ou raccourci direct 1/F, 2/A, 3/I, 4/M |
| Sous-menus en grille (attaques, ACT, objets) | Flèches (déplacement 2D dans la grille) + Entrée, ou un chiffre pour choisir et valider directement |
| Mini-jeu d'esquive | Flèches (maintenir accélère le déplacement) |
