# open-jtalk-golang
A Golang wrapper for [Open JTalk] (http://open-jtalk.sourceforge.net/)

## Dependencies

open-jtalk-golang needs:
* libhtsengine-dev
* hts_engine_API-1.09
* open_jtalk-1.08
* voice, ex.: hts-voice-nitech-jp-atr503-m001

## How to use

Set CFLAGS and LDFLAGS before running

```
// #cgo CFLAGS: -I/path/to/open_jtalk/njd -I/path/to/open_jtalk/jpcommon -I/path/to/open_jtalk/text2mecab -I/path/to/open_jtalk/mecab2njd -I/path/to/open_jtalk/njd_set_pronunciation -I/path/to/open_jtalk/njd_set_digit -I/path/to/open_jtalk/njd_set_accent_phrase -I/path/to/open_jtalk/njd_set_accent_type -I/path/to/open_jtalk/njd_set_unvoiced_vowel -I/path/to/open_jtalk/njd_set_long_vowel -I/path/to/open_jtalk/njd2jpcommon -I/path/to/open_jtalk/mecab/src
// #cgo LDFLAGS: -L/path/to/open_jtalk/njd -lnjd -lnjd -lnjd -L/path/to/open_jtalk/jpcommon -ljpcommon -L/path/to/open_jtalk/text2mecab -ltext2mecab -L/path/to/open_jtalk/mecab2njd -lmecab2njd -L/path/to/open_jtalk/njd_set_pronunciation -lnjd_set_pronunciation -L/path/to/open_jtalk/njd_set_digit -lnjd_set_digit -L/path/to/open_jtalk/njd_set_accent_phrase -lnjd_set_accent_phrase -L/path/to/open_jtalk/njd_set_accent_type -lnjd_set_accent_type -L/path/to/open_jtalk/njd_set_unvoiced_vowel -lnjd_set_unvoiced_vowel -L/path/to/open_jtalk/njd_set_long_vowel -lnjd_set_long_vowel -L/path/to/open_jtalk/njd2jpcommon -lnjd2jpcommon -L/path/to/open_jtalk/mecab/src -lmecab -L/path/to/hts_engine_API -lHTSEngine -lm -lstdc++
```

Import and use:

```
var phrase string = "むかしむかし、あるところに、おじいさんとおばあさんが住んでいました。"

/* dictionary directory */
var dn_dict string = "/home/vagrant/open_jtalk_dic_utf_8"

/* HTS voice file name */
var fn_voice string = "/usr/share/hts-voice/nitech-jp-atr503-m001/nitech_jp_atr503_m001.htsvoice"

jt := NewOpenJTalk(dn_dict, fn_voice)
jt.Synthesis(phrase, "voice.wav")
```
