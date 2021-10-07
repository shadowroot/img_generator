package img_generator

type VCardContactReq struct {
	N     []string `json:"name" query:"name"` //Name
	FN    []string `json:"fn" query:"fn"`
	TITLE []string `json:"title" query:"title"`
	ORG   []string `json:"org" query:"org"`
	TEL   []string `json:"tel" query:"tel"`
	EMAIL []string `json:"email" query:"email"`
	KEY   []string `json:"cryptkey" query:"cryptkey"` //Optional
	URL   []string `json:"url" query:"url"`
}

type QRRequestParam struct {
	IMG         IMGRequest
	X           int    `json:"x" query:"x"`
	Y           int    `json:"y" query:"y"`
	Location    string `json:"location" query:"location"`
	BlockSize   int    `json:"bs" query:"bs"`
	FillColor   string `json:"fillColor" query:"fillColor"`
	UnfillColor string `json:"unfillColor" query:"unfillColor"`
	Transparent bool   `json:"t" query:"t"`
	Opacity     bool   `json:"op" query:"op"`
	Text        string `json:"text" query:"text"`
}

type IMGRequest struct {
	W     int    `json:"w" query:"w"`
	H     int    `json:"h" query:"h"`
	Color string `json:"color" query:"color"`
}
type VCardRequest struct {
	QR         QRRequestParam
	IMG        IMGRequest
	VCARD      VCardContactReq
	Padding    int      `json:"vcard_padding" query:"vcard_padding"`
	FontFamily string   `json:"vcard_fontfamily" query:"vcard_fontfamily"`
	FontSize   int      `json:"vcard_fontsize" query:"vcard_fontsize"`
	FontUnits  string   `json:"vcard_fontunits" query:"vcard_fontunits"`
	FontColor  string   `json:"vcard_fontcolor" query:"vcard_fontcolor"`
	Options    []string `json:"vcard_options" query:"vcard_options"`
}
