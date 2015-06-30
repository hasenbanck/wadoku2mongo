package main

import (
	"golang.org/x/text/unicode/norm"
	"gopkg.in/mgo.v2"
	"regexp"
	"strings"
)

// Pre-compiled regular expressions
var (
	re_bra    *regexp.Regexp
	re_def    *regexp.Regexp
	re_tok    *regexp.Regexp
	re_tex    *regexp.Regexp
	re_exp    *regexp.Regexp
	re_lit    *regexp.Regexp
	re_fam    *regexp.Regexp
	re_emp    *regexp.Regexp
	re_transl *regexp.Regexp
	re_spe    *regexp.Regexp
	re_iro    *regexp.Regexp
	re_top    *regexp.Regexp
	re_for    *regexp.Regexp
	re_deu    *regexp.Regexp
	re_des    *regexp.Regexp
	re_bir    *regexp.Regexp
	re_tit    *regexp.Regexp
	re_dat    *regexp.Regexp
	re_jap    *regexp.Regexp
	re_transc *regexp.Regexp
)

type Reading struct {
	Hiragana string
	Hatsuon  string
}

type Entry struct {
	Id          int
	GrammType   GrammType
	Orthography []string
	Reading     []Reading
	Translation []string
}

func saveIntoMongo(dict XMLDict, connection string) error {
	session, err := mgo.Dial(connection)
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("wadoku").C("dictionary")
	for _, entry := range dict.Entries {
		orth := make([]string, 0, len(entry.Form.Orthography))
		read := make([]Reading, 0, len(entry.Form.Reading))
		trans := make([]string, 0, len(entry.Sense))

		for _, o := range entry.Form.Orthography {
			// We don't need the midashigo entries, since the are double and often contain non searchable chars like "△原"
			if !o.Midashigo {
				// We normalize just to be sure
				orth = append(orth, norm.NFKC.String(o.Text))
			}
		}

		for _, r := range entry.Form.Reading {
			// We normalize just to be sure
			read = append(read, Reading{norm.NFKC.String(r.Hiragana), norm.NFKC.String(r.Hatsuon)})
		}

		for _, s := range entry.Sense {
			var str string
			if len(s.Translations) >= 1 {
				if len(s.Translations) > 1 {
					for i, t := range s.Translations {
						if i == 0 {
							str = convXML2HTML(t.Text)
						} else {
							str = str + "; " + convXML2HTML(t.Text)
						}
					}
				} else {
					str = convXML2HTML(s.Translations[0].Text)
				}
			}
			trans = append(trans, str)
		}

		err = c.Insert(&Entry{entry.Id,
			parseGrammType(&entry),
			orth,
			read,
			trans})
		if err != nil {
			return err
		}
	}

	// Create indexes
	c.EnsureIndex(mgo.Index{Key: []string{"id"}, Unique: true})
	c.EnsureIndex(mgo.Index{Key: []string{"ortography"}})

	return nil
}

func convXML2HTML(s string) string {
	// Braket
	s = re_bra.ReplaceAllLiteralString(s, "(")
	s = strings.Replace(s, "</bracket>", ") ", -1)

	// Definition
	s = re_def.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</def>", " ", -1)

	// Token
	s = re_tok.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</token>", " ", -1)

	// Text
	s = re_tex.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</text>", " ", -1)

	// Explanation
	s = re_exp.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</expl>", " ", -1)

	// Literal
	s = re_lit.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</literal>", " ", -1)

	// Family Name
	s = re_fam.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</famn>", " ", -1)

	// Emphasis
	s = re_emp.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</emph>", " ", -1)

	// Translation
	s = re_transl.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</transl>", " ", -1)

	// Special character
	s = re_spe.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</specchar>", " ", -1)

	// Irony
	s = re_iro.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</iron>", " ", -1)

	// Topic
	s = re_top.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</topic>", " ", -1)

	// Foreign language
	s = re_for.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</foreign>", " ", -1)

	// German grammatic
	s = re_deu.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</deu_gr>", " ", -1)

	// Description
	s = re_des.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</descr>", " ", -1)

	// Birth and Death
	s = re_bir.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</birthdeath>", " ", -1)

	// Title
	s = re_tit.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</title>", " ", -1)

	// Date
	s = re_dat.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</date>", " ", -1)

	// Japanese
	s = re_jap.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</jap>", " ", -1)

	// Transcription
	s = re_transc.ReplaceAllLiteralString(s, "")
	s = strings.Replace(s, "</transcr>", " ", -1)

	// Clean up
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "  ", " ", -1)
	s = strings.Replace(s, " )", ")", -1)

	return s
}

func init() {
	re_bra = regexp.MustCompile("<bracket([^>])*>")
	re_def = regexp.MustCompile("<def([^>])*>")
	re_tok = regexp.MustCompile("<token([^>])*>")
	re_tex = regexp.MustCompile("<text([^>])*>")
	re_exp = regexp.MustCompile("<expl([^>])*>")
	re_lit = regexp.MustCompile("<literal([^>])*>")
	re_fam = regexp.MustCompile("<famn([^>])*>")
	re_emp = regexp.MustCompile("<emph([^>])*>")
	re_transl = regexp.MustCompile("<transl([^>])*>")
	re_spe = regexp.MustCompile("<specchar([^>])*>")
	re_iro = regexp.MustCompile("<iron([^>])*>")
	re_top = regexp.MustCompile("<topic([^>])*>")
	re_for = regexp.MustCompile("<foreign([^>])*>")
	re_deu = regexp.MustCompile("<deu_gr([^>])*>")
	re_des = regexp.MustCompile("<descr([^>])*>")
	re_bir = regexp.MustCompile("<birthdeath([^>])*>")
	re_tit = regexp.MustCompile("<title([^>])*>")
	re_dat = regexp.MustCompile("<date([^>])*>")
	re_jap = regexp.MustCompile("<jap([^>])*>")
	re_transc = regexp.MustCompile("<transcr([^>])*>")
}
