package openjtalk

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <math.h>

// Main headers
#include "mecab.h"
#include "njd.h"
#include "jpcommon.h"
#include "HTS_engine.h"

// Sub headers
#include "text2mecab.h"
#include "mecab2njd.h"
#include "njd_set_pronunciation.h"
#include "njd_set_digit.h"
#include "njd_set_accent_phrase.h"
#include "njd_set_accent_type.h"
#include "njd_set_unvoiced_vowel.h"
#include "njd_set_long_vowel.h"
#include "njd2jpcommon.h"

typedef const char* const_char_ptr;
void fprintfW(FILE * f, const char * c) {
  fprintf(f, c);
}

*/
import "C"
import "fmt"
import "unsafe"

type Open_JTalk struct {
	mecab    C.Mecab
	njd      C.NJD
	jpcommon C.JPCommon
	engine   C.HTS_Engine
}

const MAXBUFLEN C.int = 1024

func (open_jtalk *Open_JTalk) Initialize() {
	C.Mecab_initialize(&open_jtalk.mecab)
	C.NJD_initialize(&open_jtalk.njd)
	C.JPCommon_initialize(&open_jtalk.jpcommon)
	C.HTS_Engine_initialize(&open_jtalk.engine)
}

func (open_jtalk *Open_JTalk) Clear() {
	C.Mecab_clear(&open_jtalk.mecab)
	C.NJD_clear(&open_jtalk.njd)
	C.JPCommon_clear(&open_jtalk.jpcommon)
	C.HTS_Engine_clear(&open_jtalk.engine)
}

func (open_jtalk *Open_JTalk) Load(dn_mecab string, fn_voice string) bool {
	dn_mecab_char := (*C.char)(C.CString(dn_mecab))
	fn_voice_char := (*C.char)(C.CString(fn_voice))

	if C.Mecab_load(&open_jtalk.mecab, dn_mecab_char) != C.TRUE {
		open_jtalk.Clear()
		return false
	}
	if C.HTS_Engine_load(&open_jtalk.engine, &fn_voice_char, 1) != C.TRUE {
		open_jtalk.Clear()
		return false
	}
	base_cmp := (*C.char)(C.CString("HTS_TTS_JPN"))
	if C.strcmp(C.HTS_Engine_get_fullcontext_label_format(&open_jtalk.engine), base_cmp) != 0 {
		open_jtalk.Clear()
		return false
	}
	return true
}

func (open_jtalk *Open_JTalk) Set_sampling_frequency(i uint) {
	C.HTS_Engine_set_sampling_frequency(&open_jtalk.engine, C.size_t(i))
}

func (open_jtalk *Open_JTalk) Set_fperiod(i uint) {
	C.HTS_Engine_set_fperiod(&open_jtalk.engine, C.size_t(i))
}

func (open_jtalk *Open_JTalk) Set_alpha(f float64) {
	C.HTS_Engine_set_alpha(&open_jtalk.engine, C.double(f))
}

func (open_jtalk *Open_JTalk) Set_beta(f float64) {
	C.HTS_Engine_set_beta(&open_jtalk.engine, C.double(f))
}

func (open_jtalk *Open_JTalk) Set_speed(f float64) {
	C.HTS_Engine_set_speed(&open_jtalk.engine, C.double(f))
}

func (open_jtalk *Open_JTalk) Add_half_tone(f float64) {
	C.HTS_Engine_add_half_tone(&open_jtalk.engine, C.double(f))
}

func (open_jtalk *Open_JTalk) Set_msd_threshold(i uint, f float64) {
	C.HTS_Engine_set_msd_threshold(&open_jtalk.engine, C.size_t(i), C.double(f))
}

func (open_jtalk *Open_JTalk) Set_gv_weight(i uint, f float64) {
	C.HTS_Engine_set_gv_weight(&open_jtalk.engine, C.size_t(i), C.double(f))
}

func (open_jtalk *Open_JTalk) Set_volume(f float64) {
	C.HTS_Engine_set_volume(&open_jtalk.engine, C.double(f))
}

func (open_jtalk *Open_JTalk) Set_audio_buff_size(i uint) {
	C.HTS_Engine_set_audio_buff_size(&open_jtalk.engine, C.size_t(i))
}

/**
* @todo: log to string instead of file: logbuff := var logbuffer [1024]C.char; (*C.char)(unsafe.Pointer(&logbuffer[0])); C.NJD_sprint(&open_jtalk.njd, logbuff, (C.const_char_ptr)(C.CString(" | ")))
 */
func (open_jtalk *Open_JTalk) synthesis(text string, wavfp *C.FILE, logfp *C.FILE) (result bool) {
	txt := (C.const_char_ptr)(C.CString(text))
	result = false
	var buffer [MAXBUFLEN]C.char
	buff := (*C.char)(unsafe.Pointer(&buffer[0]))

	C.text2mecab(buff, txt)
	C.Mecab_analysis(&open_jtalk.mecab, buff)
	C.mecab2njd(&open_jtalk.njd, C.Mecab_get_feature(&open_jtalk.mecab),
		C.Mecab_get_size(&open_jtalk.mecab))
	C.njd_set_pronunciation(&open_jtalk.njd)
	C.njd_set_digit(&open_jtalk.njd)
	C.njd_set_accent_phrase(&open_jtalk.njd)
	C.njd_set_accent_type(&open_jtalk.njd)
	C.njd_set_unvoiced_vowel(&open_jtalk.njd)
	C.njd_set_long_vowel(&open_jtalk.njd)
	C.njd2jpcommon(&open_jtalk.jpcommon, &open_jtalk.njd)
	C.JPCommon_make_label(&open_jtalk.jpcommon)
	if C.JPCommon_get_label_size(&open_jtalk.jpcommon) > 2 {
		if C.HTS_Engine_synthesize_from_strings(&open_jtalk.engine,
			C.JPCommon_get_label_feature(&open_jtalk.jpcommon),
			C.size_t(C.JPCommon_get_label_size(&open_jtalk.jpcommon))) == C.TRUE {
			result = true
		}
		if wavfp != nil {
			C.HTS_Engine_save_riff(&open_jtalk.engine, wavfp)
		}
		if logfp != nil {
			C.fprintfW(logfp, (C.const_char_ptr)(C.CString("[Text analysis result]\n")))
			C.NJD_fprint(&open_jtalk.njd, logfp)
			C.fprintfW(logfp, (C.const_char_ptr)(C.CString("\n[Output label]\n")))
			C.HTS_Engine_save_label(&open_jtalk.engine, logfp)
			C.fprintfW(logfp, (C.const_char_ptr)(C.CString("\n")))
			C.HTS_Engine_save_information(&open_jtalk.engine, logfp)
		}
		C.HTS_Engine_refresh(&open_jtalk.engine)
	}
	C.JPCommon_refresh(&open_jtalk.jpcommon)
	C.NJD_refresh(&open_jtalk.njd)
	C.Mecab_refresh(&open_jtalk.mecab)

	return result
}

func (open_jtalk *Open_JTalk) Synthesis(phrase string, wavFilename string) {

	/* output file pointers */
	var wavfp *C.FILE = C.fopen(C.CString(wavFilename), C.CString("wb"))

	/* synthesize */
	if !open_jtalk.synthesis(phrase, wavfp, nil) {
		fmt.Println("Error: waveform cannot be synthesized.\n")
		open_jtalk.Clear()
		panic("Error: waveform cannot be synthesized.\n")
	}

	/* close files */
	if wavfp != nil {
		C.fclose(wavfp)
	}
}

func (open_jtalk *Open_JTalk) SynthesisAndLog(phrase string, wavFilename string, logFilname string) {

	/* output file pointers */
	var logfp *C.FILE = C.fopen(C.CString(logFilname), C.CString("wt"))
	var wavfp *C.FILE = C.fopen(C.CString(wavFilename), C.CString("wb"))

	/* synthesize */
	if !open_jtalk.synthesis(phrase, wavfp, logfp) {
		fmt.Println("Error: waveform cannot be synthesized.\n")
		open_jtalk.Clear()
		panic("Error: waveform cannot be synthesized.\n")
	}

	/* close files */
	if wavfp != nil {
		C.fclose(wavfp)
	}
	if logfp != nil {
		C.fclose(logfp)
	}
}

/**
* @todo: implement options
 */
func NewOpenJTalk(dn_dict string, fn_voice string) (open_jtalk Open_JTalk) {
	/* initialize Open JTalk */
	open_jtalk.Initialize()

	/* load dictionary and HTS voice */
	if !open_jtalk.Load(dn_dict, fn_voice) {
		fmt.Println("Error: Dictionary or HTS voice cannot be loaded.\n")
		open_jtalk.Clear()
		panic("Error: Dictionary or HTS voice cannot be loaded.\n")
	}

	/* skip options */

	return open_jtalk
}
