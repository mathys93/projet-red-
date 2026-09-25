Projet RED

Un RPG de combat en ligne de commande, écrit en Go, sur le thème de "JoJo's Bizarre Adventure". Le jeu emprunte son système de combat à Undertale (menus FIGHT / ACT / ITEM / MERCY et mini-jeu d'esquive en temps réel) et l'habille avec des boss, des Stands et de l'ASCII art tirés de la license JoJo.

Aperçu

Tu incarnes un personnage qui explore une carte de zones, chacune gardée par un boss (DIO, Diavolo, Yoshikage Kira...). Chaque combat alterne entre une phase d'action (attaquer, utiliser une capacité, consulter l'inventaire, tenter la pitié) et une phase d'esquive façon bullet hell où il faut survivre à l'attaque du boss dans une arène ASCII, avec un pattern d'attaque propre à chaque adversaire.

En progressant, le joueur gagne de l'XP (montée de niveau), des Fragments (monnaie du jeu) et peut débloquer un Stand aléatoire (objet "Arrow") qui ajoute une capacité de coup de grâce puissante et un bonus d'attaque permanent.

 Fonctionnalités

Le jeu propose une carte du monde avec des zones à débloquer dans l'ordre, une zone boutique et une zone "farm" (combat répétable à volonté pour l'XP et les Fragments. Le combat se joue au tour par tour avec un rendu ASCII animé : boîtes, barres de vie, de MP et d'XP, art du boss recadré à la taille du terminal.

Chaque boss propose un mini-jeu d'esquive en temps réel avec un pattern d'attaque distinct : arrêt du temps pour DIO, effacement pour Diavolo, bombes pour Kira, accélération pour Pucci, et ainsi de suite.

Le système de Stands permet, à l'utilisation d'une "Arrow", de tirer aléatoirement un Stand parmi trois raretés (Commun, Rare, Légendaire), chacun avec un bonus d'attaque et une capacité unique (Star Platinum, The World, Crazy Diamond, Killer Queen...).

La boutique se compose de deux volets : le Marchand, qui vend des potions, des disques de Pucci, du poison et des Arrow, et le Forgeron, qui fabrique de l'équipement contre des Fragments.

Le joueur dispose également d'un inventaire et d'objets consommables (soin, MP, poison, résurrection via King Crimson), d'une musique de fond par boss (lecture MP3 sur Windows), et d'un outil interne de génération d'ASCII art à partir d'images pour produire les visuels des boss.

## Architecture du projet

```
projet-red-/
├── main.go
├── assets/          ressources (dessins, logos, musiques)
├── internal/        code du jeu
│   ├── terminal/    écran, couleurs, clavier, moteur des fonds animés
│   ├── fond/        les fonds animés de chaque page
│   ├── ui/          menu principal, paramètres, crédits, menu pause, carte
│   ├── combat/      combat, esquive, objets
│   ├── boutique/    boutique, Marchand, Forgeron
│   ├── boss/        boss et leurs patterns
│   ├── character/   personnage et Stands
│   ├── world/       zones de la carte
│   ├── ascii/       rendu des dessins et du logo
│   └── audio/       musique
└── tools/img2ascii/ convertisseur image vers ASCII/braille
```

| Dossier | Rôle |
|---|---|
| [main.go](main.go) | Point d'entrée : menu principal, création du joueur, boucle carte, combat, récompenses, crédits de fin. |
| [assets/](assets/) | Dessins des boss ("art/*.txt"), logos du titre ("logo/*.png"), tous deux intégrés à l'exécutable, et musiques ("musique/*.mp3"). |
| [internal/terminal/](internal/terminal/) | Bas niveau : taille du terminal, couleurs, cadres, lecture clavier multiplateforme, et l'écran qui fait passer les fonds animés derrière le contenu. |
| [internal/fond/](internal/fond/) | Les fonds animés : ciel étoilé (menus), carte, boutique, et un décor par boss. |
| [internal/ui/](internal/ui/) | Écrans hors combat : menu principal, saisie du nom, paramètres, crédits, menu pause (Échap), carte. |
| [internal/combat/](internal/combat/) | Boucle de combat ("battle.go"), mini-jeu d'esquive ("dodge.go"), objets ("items.go"), sous-menus ("submenu.go"). |
| [internal/boutique/](internal/boutique/) | Écran de la boutique, Marchand et Forgeron. |
| [internal/boss/](internal/boss/) | Modèle des boss, leurs statistiques et leurs patterns ("patterns.go"). |
| [internal/character/](internal/character/) | Personnage (PV, AP, DP, MP, XP, inventaire) et Stands ("stand.go"). |
| [internal/world/](internal/world/) | Zones de la carte et déverrouillage. |
| [internal/ascii/](internal/ascii/) | Mise à l'échelle des dessins ("Fit") et accès aux logos. |
| [internal/audio/](internal/audio/) | Piste musicale de chaque boss et lecture MP3 (Windows). |
| [tools/img2ascii/](tools/img2ascii/) | Convertisseur image vers ASCII/braille pour produire les fichiers "assets/art/*.txt". |

Prérequis

Go 1.26 ou supérieur (voir [go.mod](go.mod)), disponible sur https://go.dev/, ainsi qu'un terminal qui supporte les couleurs ANSI et l'affichage Unicode (recommandé : Windows Terminal).

## Installation et lancement

```powershell
git clone <url-du-repo>
cd projet-red-
go run .
```

Ou pour compiler un exécutable :

```powershell
go build -o projet-red.exe .
./projet-red.exe
```

## Comment jouer

La navigation se fait aux flèches directionnelles ou avec "Z/Q/S/D" (disposition AZERTY), "Entrée" pour valider et "X" pour revenir en arrière ou annuler. "Échap" ouvre le menu pause (Paramètres, Crédits, Quitter) : c'est par là qu'on quitte la partie.

En combat, quatre actions sont disponibles, à la Undertale : FIGHT pour une attaque physique ou une capacité de Stand (coûte des MP), ACT pour des options spécifiques au boss (observer, tenter un dialogue...), ITEM pour utiliser un objet de l'inventaire (potion, Disque de Pucci, Arrow...), et MERCY pour épargner le boss si ses PV sont assez bas.

Vient ensuite la phase d'esquive : après ton tour, le boss attaque et il faut déplacer ton cœur dans l'arène ASCII pour éviter les projectiles, selon le pattern propre au boss.

 Zones et boss

| Zone | Boss | Particularité |
|---|---|---|
| Zone 1, Le Manoir de DIO | DIO | Arrête le temps ("AttackTimeStop") |
| Zone 2, La Villa de Diavolo | Diavolo | Efface ses traces ("AttackErase") |
| Zone 3, Morioh | Yoshikage Kira | Pose des bombes ("AttackBombs") |
| Zone 4, La Boutique | aucun | Marchand et Forgeron, pas de combat |
| Zone 5, Le Terrain Vague | Iggy | Zone "farm", combat répétable à volonté |
| Zone bonus | Enrico Pucci | Débloqué une fois les 3 zones principales terminées, accélère le temps ("AttackAccelerate") |

Outils annexes

La commande "go run ./tools/img2ascii -in image.png -out assets/art/nom.txt -cols 60" convertit une image en art ASCII/braille (dithering, détection de contours, contraste réglables). Le fichier doit ensuite être ajouté dans [assets/assets.go](assets/assets.go).

 Notes

La lecture audio et la maximisation de la fenêtre console sont implémentées spécifiquement pour Windows ("*_windows.go"), avec des variantes neutres pour les autres systèmes d'exploitation ("*_other.go").

Les musiques de boss sont configurables dans [internal/audio/tracks.go](internal/audio/tracks.go). Chaque boss a sa piste dans "assets/musique/" : ce dossier doit rester à côté de l'exécutable pour avoir le son.

Chaque page a son fond animé, dessiné uniquement dans les cases libres : il ne recouvre jamais le texte, les cadres ni les dessins.
