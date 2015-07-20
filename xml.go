package main

type XMLReading struct {
	Hiragana string `xml:"hira"`
	Hatsuon  string `xml:"hatsuon"`
}

type XMLOrthography struct {
	Text      string `xml:",chardata"`
	Midashigo bool   `xml:"midashigo,attr"` // Schlagwort; Eintrag sollte immer "doppel" vorhanden sein -> Prüfen!
}

type XMLForm struct {
	Orthography []XMLOrthography `xml:"orth"` // Später in Mongo DB -> Multikey!
	Reading     []XMLReading     `xml:"reading"`
}

type XMLTranslation struct {
	Text string `xml:",innerxml"`
}

type XMLSense struct {
	Translations []XMLTranslation `xml:"trans>tr"`
}

type XMLGrammGroup struct {
	Doushi           *struct{} `xml:"doushi"`
	Meishi           *struct{} `xml:"meishi"`
	Keiyoushi        *struct{} `xml:"keiyoushi"`
	Keiyoudoushi     *struct{} `xml:"keiyoudoushi"`
	Fukushi          *struct{} `xml:"fukushi"`
	Prefix           *struct{} `xml:"prefix"`
	Suffix           *struct{} `xml:"suffix"`
	Rentaishi        *struct{} `xml:"rentaishi"`
	Setsuzokushi     *struct{} `xml:"setsuzokushi"`
	Daimeishi        *struct{} `xml:"daimeishi"`
	Kandoushi        *struct{} `xml:"kandoushi"`
	Jodoushi         *struct{} `xml:"jodoushi"`
	Jokeiyoushi      *struct{} `xml:"jokeiyoushi"`
	Joshi            *struct{} `xml:"joshi"`
	Fukujoshi        *struct{} `xml:"fukujoshi"`
	Kakarijoshi      *struct{} `xml:"kakarijoshi"`
	Shuujoshi        *struct{} `xml:"shuujoshi"`
	Wordcomponent    *struct{} `xml:"wordcomponent"`
	Specialcharacter *struct{} `xml:"specialcharacter"`
	Rengo            *struct{} `xml:"rengo"`
	Kanji            *struct{} `xml:"kanji"`
}

type XMLEntry struct {
	Id         int           `xml:"id,attr"`
	Form       XMLForm       `xml:"form"`
	GrammGroup XMLGrammGroup `xml:"gramGrp"`
	Sense      []XMLSense    `xml:"sense"`
}

type XMLDict struct {
	Entries []XMLEntry `xml:"entry"`
}

func parseGrammType(entry *XMLEntry) grammType {
	switch {
	case entry.GrammGroup.Daimeishi != nil:
		return Daimeishi
	case entry.GrammGroup.Doushi != nil:
		return Doushi
	case entry.GrammGroup.Fukujoshi != nil:
		return Fukujoshi
	case entry.GrammGroup.Fukushi != nil:
		return Fukushi
	case entry.GrammGroup.Jodoushi != nil:
		return Jodoushi
	case entry.GrammGroup.Jokeiyoushi != nil:
		return Jokeiyoushi
	case entry.GrammGroup.Joshi != nil:
		return Joshi
	case entry.GrammGroup.Kakarijoshi != nil:
		return Kakarijoshi
	case entry.GrammGroup.Kandoushi != nil:
		return Kandoushi
	case entry.GrammGroup.Kanji != nil:
		return Kanji
	case entry.GrammGroup.Keiyoudoushi != nil:
		return Keiyoudoushi
	case entry.GrammGroup.Keiyoushi != nil:
		return Keiyoushi
	case entry.GrammGroup.Meishi != nil:
		return Meishi
	case entry.GrammGroup.Prefix != nil:
		return Prefix
	case entry.GrammGroup.Rengo != nil:
		return Rengo
	case entry.GrammGroup.Rentaishi != nil:
		return Rentaishi
	case entry.GrammGroup.Setsuzokushi != nil:
		return Setsuzokushi
	case entry.GrammGroup.Shuujoshi != nil:
		return Shuujoshi
	case entry.GrammGroup.Specialcharacter != nil:
		return Specialcharacter
	case entry.GrammGroup.Suffix != nil:
		return Suffix
	case entry.GrammGroup.Wordcomponent != nil:
		return Wordcomponent
	default:
		return NotDefined
	}
}
