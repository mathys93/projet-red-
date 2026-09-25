package assets

import _ "embed"

//go:embed logo/title.png
var TitlePNG []byte

//go:embed logo/subtitle.png
var SubtitlePNG []byte

//go:embed art/dio.txt
var DioArt string

//go:embed art/diavolo.txt
var DiavoloArt string

//go:embed art/kira.txt
var KiraArt string

//go:embed art/disque_de_pucci.txt
var PucciArt string

//go:embed art/iggy.txt
var IggyArt string

//go:embed art/shop.txt
var ShopArt string
