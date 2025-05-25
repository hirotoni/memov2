package search

import (
	"bufio"
	"bytes"
	"strings"
)

type IRomajiConverter interface {
	Convert(query string) []string
}

// RomajiConverter handles conversion from romaji to Japanese using SKK dictionary
type RomajiConverter struct {
	dict map[string][]string // hiragana -> kanji mapping
}

// Initialize katakana mapping
var hiraganaToKatakana = map[rune]rune{
	'あ': 'ア', 'い': 'イ', 'う': 'ウ', 'え': 'エ', 'お': 'オ',
	'か': 'カ', 'き': 'キ', 'く': 'ク', 'け': 'ケ', 'こ': 'コ',
	'さ': 'サ', 'し': 'シ', 'す': 'ス', 'せ': 'セ', 'そ': 'ソ',
	'た': 'タ', 'ち': 'チ', 'つ': 'ツ', 'て': 'テ', 'と': 'ト',
	'な': 'ナ', 'に': 'ニ', 'ぬ': 'ヌ', 'ね': 'ネ', 'の': 'ノ',
	'は': 'ハ', 'ひ': 'ヒ', 'ふ': 'フ', 'へ': 'ヘ', 'ほ': 'ホ',
	'ま': 'マ', 'み': 'ミ', 'む': 'ム', 'め': 'メ', 'も': 'モ',
	'や': 'ヤ', 'ゆ': 'ユ', 'よ': 'ヨ',
	'ら': 'ラ', 'り': 'リ', 'る': 'ル', 'れ': 'レ', 'ろ': 'ロ',
	'わ': 'ワ', 'を': 'ヲ', 'ん': 'ン',
	'が': 'ガ', 'ぎ': 'ギ', 'ぐ': 'グ', 'げ': 'ゲ', 'ご': 'ゴ',
	'ざ': 'ザ', 'じ': 'ジ', 'ず': 'ズ', 'ぜ': 'ゼ', 'ぞ': 'ゾ',
	'だ': 'ダ', 'ぢ': 'ヂ', 'づ': 'ヅ', 'で': 'デ', 'ど': 'ド',
	'ば': 'バ', 'び': 'ビ', 'ぶ': 'ブ', 'べ': 'ベ', 'ぼ': 'ボ',
	'ぱ': 'パ', 'ぴ': 'ピ', 'ぷ': 'プ', 'ぺ': 'ペ', 'ぽ': 'ポ',
	'ゃ': 'ャ', 'ゅ': 'ュ', 'ょ': 'ョ',
	'っ': 'ッ', 'ー': 'ー',
}

// romajiToHiragana maps romaji to hiragana
var romajiToHiragana = map[string]string{
	// Basic hiragana
	"a": "あ", "i": "い", "u": "う", "e": "え", "o": "お",
	"ka": "か", "ki": "き", "ku": "く", "ke": "け", "ko": "こ",
	"sa": "さ", "si": "し", "shi": "し", "su": "す", "se": "せ", "so": "そ",
	"ta": "た", "ti": "ち", "chi": "ち", "tu": "つ", "tsu": "つ", "te": "て", "to": "と",
	"na": "な", "ni": "に", "nu": "ぬ", "ne": "ね", "no": "の",
	"ha": "は", "hi": "ひ", "fu": "ふ", "hu": "ふ", "he": "へ", "ho": "ほ",
	"ma": "ま", "mi": "み", "mu": "む", "me": "め", "mo": "も",
	"ya": "や", "yu": "ゆ", "yo": "よ",
	"ra": "ら", "ri": "り", "ru": "る", "re": "れ", "ro": "ろ",
	"wa": "わ", "wo": "を", "n": "ん", "nn": "ん",
	// Voiced sounds (dakuten)
	"ga": "が", "gi": "ぎ", "gu": "ぐ", "ge": "げ", "go": "ご",
	"za": "ざ", "zi": "じ", "ji": "じ", "zu": "ず", "ze": "ぜ", "zo": "ぞ",
	"da": "だ", "di": "ぢ", "du": "づ", "de": "で", "do": "ど",
	"ba": "ば", "bi": "び", "bu": "ぶ", "be": "べ", "bo": "ぼ",
	"pa": "ぱ", "pi": "ぴ", "pu": "ぷ", "pe": "ぺ", "po": "ぽ",
	// Combinations
	"kya": "きゃ", "kyu": "きゅ", "kyo": "きょ",
	"sha": "しゃ", "sya": "しゃ", "shu": "しゅ", "syu": "しゅ", "sho": "しょ", "syo": "しょ",
	"cha": "ちゃ", "tya": "ちゃ", "chu": "ちゅ", "tyu": "ちゅ", "cho": "ちょ", "tyo": "ちょ",
	"nya": "にゃ", "nyu": "にゅ", "nyo": "にょ",
	"hya": "ひゃ", "hyu": "ひゅ", "hyo": "ひょ",
	"mya": "みゃ", "myu": "みゅ", "myo": "みょ",
	"rya": "りゃ", "ryu": "りゅ", "ryo": "りょ",
	"gya": "ぎゃ", "gyu": "ぎゅ", "gyo": "ぎょ",
	"ja": "じゃ", "ju": "じゅ", "jo": "じょ",
	"zya": "じゃ", "zyu": "じゅ", "zyo": "じょ",
	"dya": "ぢゃ", "dyu": "ぢゅ", "dyo": "ぢょ",
	"bya": "びゃ", "byu": "びゅ", "byo": "びょ",
	"pya": "ぴゃ", "pyu": "ぴゅ", "pyo": "ぴょ",
	// Small tsu for double consonants
	// Basic hiragana with っ
	"kka": "っか", "kki": "っき", "kku": "っく", "kke": "っけ", "kko": "っこ",
	"ssa": "っさ", "sshi": "っし", "ssi": "っし", "ssu": "っす", "sse": "っせ", "sso": "っそ",
	"tta": "った", "tti": "っち", "tchi": "っち", "ttsu": "っつ", "ttu": "っつ", "tte": "って", "tto": "っと",
	"ppa": "っぱ", "ppi": "っぴ", "ppu": "っぷ", "ppe": "っぺ", "ppo": "っぽ",
	"hha": "っは", "hhi": "っひ", "hhu": "っふ", "hhe": "っへ", "hho": "っほ",
	// Voiced sounds with っ
	"gga": "っが", "ggi": "っぎ", "ggu": "っぐ", "gge": "っげ", "ggo": "っご",
	"zza": "っざ", "zzi": "っじ", "jji": "っじ", "zzu": "っず", "zze": "っぜ", "zzo": "っぞ",
	"dda": "っだ", "ddi": "っぢ", "ddu": "っづ", "dde": "っで", "ddo": "っど",
	"bba": "っば", "bbi": "っび", "bbu": "っぶ", "bbe": "っべ", "bbo": "っぼ",
	// Combinations with っ
	"kkya": "っきゃ", "kkyu": "っきゅ", "kkyo": "っきょ",
	"ssha": "っしゃ", "ssya": "っしゃ", "sshu": "っしゅ", "ssyu": "っしゅ", "ssho": "っしょ", "ssyo": "っしょ",
	"tcha": "っちゃ", "ttya": "っちゃ", "tchu": "っちゅ", "ttyu": "っちゅ", "tcho": "っちょ", "ttyo": "っちょ",
	"hhya": "っひゃ", "hhyu": "っひゅ", "hhyo": "っひょ",
	"ggya": "っぎゃ", "ggyu": "っぎゅ", "ggyo": "っぎょ",
	"jja": "っじゃ", "jju": "っじゅ", "jjo": "っじょ",
	"zzya": "っじゃ", "zzyu": "っじゅ", "zzyo": "っじょ",
	"ddya": "っぢゃ", "ddyu": "っぢゅ", "ddyo": "っぢょ",
	"bbya": "っびゃ", "bbyu": "っびゅ", "bbyo": "っびょ",
	"ppya": "っぴゃ", "ppyu": "っぴゅ", "ppyo": "っぴょ",
	// Additional rows with っ
	"nna": "っな", "nni": "っに", "nnu": "っぬ", "nne": "っね", "nno": "っの",
	"mma": "っま", "mmi": "っみ", "mmu": "っむ", "mme": "っめ", "mmo": "っも",
	"yya": "っや", "yyu": "っゆ", "yyo": "っよ",
	"rra": "っら", "rri": "っり", "rru": "っる", "rre": "っれ", "rro": "っろ",
	"wwa": "っわ", "wwo": "っを",
	// Additional combinations with っ
	"nnya": "っにゃ", "nnyu": "っにゅ", "nnyo": "っにょ",
	"mmya": "っみゃ", "mmyu": "っみゅ", "mmyo": "っみょ",
	"rrya": "っりゃ", "rryu": "っりゅ", "rryo": "っりょ",
	// Long vowels
	"-": "ー",
}

// NewRomajiConverter creates a new RomajiConverter instance and loads the dictionary
func NewRomajiConverter() (*RomajiConverter, error) {
	rc := &RomajiConverter{
		dict: make(map[string][]string),
	}
	if err := rc.loadDictionary(); err != nil {
		return nil, err
	}
	return rc, nil
}

// loadDictionary loads the SKK dictionary from embedded data
func (rc *RomajiConverter) loadDictionary() error {
	reader := bytes.NewReader(skkDictData)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments and empty lines
		if strings.HasPrefix(line, ";") || line == "" {
			continue
		}

		// SKK dictionary format: hiragana /kanji1;comment/kanji2;comment/...
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			continue
		}

		hiragana := parts[0]
		// Split by slash and remove comments
		rawCandidates := strings.Split(strings.Trim(parts[1], "/"), "/")
		candidates := make([]string, 0, len(rawCandidates))
		for _, c := range rawCandidates {
			if semicolonIdx := strings.Index(c, ";"); semicolonIdx != -1 {
				candidates = append(candidates, c[:semicolonIdx])
			} else {
				candidates = append(candidates, c)
			}
		}
		rc.dict[hiragana] = candidates
	}

	return scanner.Err()
}

// Convert converts romaji query to possible Japanese candidates
func (rc *RomajiConverter) Convert(query string) []string {
	// If query is already in Japanese, still try to get variations
	var results []string

	// Add original query
	results = append(results, query)

	// Convert romaji to hiragana if input is not Japanese
	var hiragana string
	if containsJapanese(query) {
		hiragana = query
	} else {
		hiragana = rc.toHiragana(query)
		if hiragana != "" {
			results = append(results, hiragana)
		}
	}

	// Convert hiragana to katakana
	if katakana := rc.toKatakana(hiragana); katakana != "" {
		results = append(results, katakana)
	}

	// Look up kanji in SKK dictionary
	if candidates, ok := rc.dict[hiragana]; ok {
		results = append(results, candidates...)
	}

	// Remove duplicates while preserving order
	return rc.removeDuplicates(results)
}

// toHiragana converts romaji string to hiragana
func (rc *RomajiConverter) toHiragana(romaji string) string {
	romaji = strings.ToLower(romaji)
	var result strings.Builder
	for i := 0; i < len(romaji); {
		// Try to match longest possible romaji sequence
		matched := false
		for j := min(len(romaji), i+4); j > i; j-- {
			if hiragana, ok := romajiToHiragana[romaji[i:j]]; ok {
				result.WriteString(hiragana)
				i = j
				matched = true
				break
			}
		}
		if !matched {
			// If no match found, skip this character
			i++
		}
	}
	return result.String()
}

// toKatakana converts hiragana to katakana
func (rc *RomajiConverter) toKatakana(hiragana string) string {
	var result strings.Builder
	for _, r := range hiragana {
		if kata, ok := hiraganaToKatakana[r]; ok {
			result.WriteRune(kata)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// removeDuplicates removes duplicate strings while preserving order
func (rc *RomajiConverter) removeDuplicates(strs []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(strs))
	for _, str := range strs {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}
	return result
}

// containsJapanese checks if the string contains Japanese characters
func containsJapanese(s string) bool {
	for _, r := range s {
		if (r >= 0x3040 && r <= 0x309F) || // Hiragana
			(r >= 0x30A0 && r <= 0x30FF) || // Katakana
			(r >= 0x4E00 && r <= 0x9FFF) { // Kanji
			return true
		}
	}
	return false
}
