package img_generator

import (
	"fmt"
	"strings"
)

// type VCardContact struct {
// 	N     string //Name
// 	FN    string
// 	TITLE string
// 	ORG   []string
// 	TEL   []string
// 	EMAIL []string
// 	KEY   []string //Optional
// 	URL   []string
// }

type VCardContact map[string][]string

type WIFIQR struct {
	T string //Type enc
	S string //SSID
	P string //Password "nopass" omnitted
	H bool   //Hidden false
}

var ()

func stringArrTovCard(arr []string) string {
	return strings.Join(arr, ";")
}

func telToString(tels map[string]string) string {
	ret := ""
	for k, v := range tels {
		ret += fmt.Sprintf("TEL;VALUE=uri;TYPE=%v:tel:%v\n", k, v)
	}
	return ret
}

func NewVCardContact(FN []string, N []string, TITLE []string, ORG []string, TEL []string, EMAIL []string, URL []string, KEY []string) *VCardContact {
	return &VCardContact{
		"N":     N,
		"FN":    FN,
		"TITLE": TITLE,
		"ORG":   ORG,
		"TEL":   TEL,
		"EMAIL": EMAIL,
		"URL":   URL,
		"KEY":   KEY,
	}
}

func NewWIFIQR(T string, //Type enc
	S string, //SSID
	P string, //Password "nopass" omnitted
	H bool) *WIFIQR {
	return &WIFIQR{
		T: T,
		S: S,
		P: P,
		H: H,
	}
}

/*
https://github.com/zxing/zxing/wiki/Barcode-Contents
https://tools.ietf.org/html/rfc6350


BEGIN:VCARD
VERSION:4.0
N:Owen;Sean;<pre titles>;<post titles>;
FN:Sean Owen
TITLE:Software Engineer
ORG:<types delimited by ;>
ROLE:
#TEL;VALUE=uri;PREF=1;TYPE="voice,home":tel:+1-555-555-5555;ext=5555
TEL;VALUE=uri;TYPE=home:tel:+33-01-23-45-67
EMAIL;TYPE=WORK:srowen@google.com
URL:https://example.com
#KEY:data:application/pgp-keys;base64,MIICajCCAdOgAwIBAgICBE
      UwDQYJKoZIhvcNAQEEBQAwdzELMAkGA1UEBhMCVVMxLDAqBgNVBAoTI05l
      <... remainder of base64-encoded data ...>
END:VCARD
*/
func (c VCardContact) ToQRFormat() (string, error) {
	vCard := "BEGIN:VCARD\nVERSION:4.0\n"
	//log.Debug().Msgf("vCard=%v", c)
	for k, vals := range c {
		for _, v := range vals {
			vCard += fmt.Sprintf("%v:%v\n", k, v)
		}
	}
	vCard += "END:VCARD"
	return vCard, nil
}

func (c VCardContact) ToLinesFormat() ([]string, error) {
	lines := make([]string, 1)
	for k, vals := range c {
		for _, v := range vals {
			line := k
			lines = append(lines, line+":"+v)
		}

	}
	return lines, nil
}

/*
WIFI:T:WPA;S:mynetwork;P:mypass;H:true;

Parameter	Example	Description
T	WPA	Authentication type; can be WEP or WPA, or 'nopass' for no password. Or, omit for no password.
S	mynetwork	Network SSID. Required. Enclose in double quotes if it is an ASCII name, but could be interpreted as hex (i.e. "ABCD")
P	mypass	Password, ignored if T is "nopass" (in which case it may be omitted). Enclose in double quotes if it is an ASCII name, but could be interpreted as hex (i.e. "ABCD")
H	true	Optional. True if the network SSID is hidden.
*/

func (wifi WIFIQR) ToQRFormat() string {
	str := "WIFI:"
	if len(wifi.T) > 0 {
		str += "T:" + wifi.T + ";"
	}
	str += "S:" + wifi.S + ";"
	if len(wifi.P) <= 0 {
		str += "P:" + "\"nopass\"" + ";"
	} else {
		str += "P:" + wifi.P + ";"
	}
	if wifi.H {
		str += "H:true;"
	}
	return str
}
