// Package initcharacter est OBSOLÈTE.
//
// Ce dossier contenait un fichier nommé "main.go" mais déclarant
// `package initcharacter` (pas `package main`) : nom trompeur, à éviter.
// Sa fonction InitCharacter(...) ne faisait par ailleurs rien d'utile :
// elle réassignait ses propres paramètres (name, AP, DP...) sans jamais
// les renvoyer ni les stocker nulle part - un personnage construit avec
// cette fonction restait vide, ses "valeurs par défaut" étaient perdues
// dès la sortie de la fonction.
//
// Toute la création de personnage a été refaite proprement dans le
// package character (voir character/character.go, fonction character.New).
//
// Tu peux supprimer tout le dossier initcharacter/ du projet : plus rien
// n'en dépend. Ces deux fichiers sont laissés vides (juste la déclaration
// de package) pour que le programme continue de compiler même si tu ne les
// supprimes pas tout de suite.
package initcharacter
