// Fichier input.go : lecture des commandes clavier du package combat.
//
// Avant cette réorganisation, main.go activait le mode "une touche = une
// action" en appelant la commande système "stty" via exec.Command. Cette
// commande n'existe pas sous Windows (cmd.exe / PowerShell) : sttyRun()
// échouait silencieusement (l'erreur était ignorée), le terminal restait
// en mode ligne classique, et readKey() lisait des séquences d'échappement
// qui n'arrivent jamais dans ce mode. Le combat ne fonctionnait donc pas du
// tout pour un joueur Windows (le cas ici) sans WSL ni Git Bash.
//
// Plutôt que de dépendre d'un paquet externe (golang.org/x/term) ou de
// coder à la main les appels système Windows/Linux/macOS (beaucoup de code
// spécifique à chaque OS pour un résultat fragile), le jeu utilise ici un
// mode de saisie "une touche + Entrée", identique sur toutes les
// plateformes avec seulement la bibliothèque standard de Go (bufio). C'est
// un compromis assumé : on perd le côté "flèche pressée = action
// immédiate", mais le jeu tourne de façon identique et fiable sous
// Windows, macOS et Linux.
package combat

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Key représente une commande "logique" reconnue par le jeu.
type Key string

const (
	KeyUp    Key = "UP"
	KeyDown  Key = "DOWN"
	KeyLeft  Key = "LEFT"
	KeyRight Key = "RIGHT"
	KeyEnter Key = "ENTER"
	KeyQuit  Key = "QUIT"
	KeyOther Key = "OTHER"
)

// stdinScanner est UNIQUE et partagé par tous les écrans du jeu (carte,
// combats successifs) pendant toute la durée du programme.
//
// Piège évité ici : bufio.Scanner lit par blocs (pas ligne par ligne au
// niveau système). Si chaque écran créait son propre bufio.Scanner(os.Stdin)
// (comme c'était le cas dans une version précédente de ce fichier), une
// entrée tapée un peu à l'avance - ou fournie d'un coup, par exemple via un
// script/pipe de test - pouvait être engloutie dans le buffer interne d'un
// scanner puis carrément perdue dès qu'un nouveau scanner était créé pour
// l'écran suivant (le combat suivant recevait alors un EOF immédiat et se
// fermait tout seul). Un unique scanner, créé une fois, élimine le problème.
var stdinScanner = bufio.NewScanner(os.Stdin)

// terminalSession représente une session de lecture clavier pour un écran
// donné (carte ou combat). Elle ne possède plus son propre scanner : elle
// s'appuie sur stdinScanner, partagé.
//
// raw garde la dernière ligne lue (en minuscules, sans espaces) même
// quand elle ne correspond à aucune Key connue (KeyOther) : ça sert aux
// raccourcis directs des menus de combat (taper "2" ou "a" puis Entrée
// pour choisir ET valider une option en une seule saisie, plutôt que de
// naviguer option par option avant de valider - voir combat/submenu.go et
// RunBattle dans battle.go).
type terminalSession struct {
	raw string
}

func newTerminalSession() *terminalSession {
	return &terminalSession{}
}

// restore réaffiche le curseur (rien d'autre à restaurer : on n'a jamais
// mis le terminal dans un mode spécial).
func (ts *terminalSession) restore() {
	fmt.Print("\033[?25h")
}

// WaitEnter attend une simple pression sur Entrée (utilisé par exemple pour
// l'écran de bienvenue). Passe par le même stdinScanner partagé que
// readKey, pour la même raison : ne jamais lire os.Stdin via deux buffers
// différents.
func WaitEnter() {
	stdinScanner.Scan()
}

// readKey lit une ligne au clavier et la traduit en commande logique.
// Contrôles : Z/W ou "up" (haut), S ou "down" (bas), Q/A ou "left" (gauche),
// D ou "right" (droite), Entrée seule pour valider, X/"quit" pour quitter.
func (ts *terminalSession) readKey() Key {
	if !stdinScanner.Scan() {
		return KeyQuit
	}
	line := strings.ToLower(strings.TrimSpace(stdinScanner.Text()))
	ts.raw = line
	switch line {
	case "":
		return KeyEnter
	case "z", "w", "up", "haut":
		return KeyUp
	case "s", "down", "bas":
		return KeyDown
	case "q", "a", "left", "gauche":
		return KeyLeft
	case "d", "right", "droite":
		return KeyRight
	case "x", "quit", "exit", "quitter":
		return KeyQuit
	}
	return KeyOther
}
