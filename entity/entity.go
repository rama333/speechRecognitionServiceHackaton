package entity

import "encoding/xml"

type Recording struct {
	XMLName     xml.Name `xml:"recording"`
	Text        string   `xml:",chardata"`
	Version     string   `xml:"version,attr"`
	Inum        string   `xml:"inum,attr"`
	Audiochans  string   `xml:"audiochans,attr"`
	Screens     string   `xml:"screens,attr"`
	Audioformat string   `xml:"audioformat,attr"`
	Starttime   string   `xml:"starttime,attr"`
	Endtime     string   `xml:"endtime,attr"`
	Filename    string   `xml:"filename,attr"`
	Moved       string   `xml:"moved,attr"`
	Cti         struct {
		Text          string `xml:",chardata"`
		Internal      string `xml:"internal,attr"`
		Incoming      string `xml:"incoming,attr"`
		Nostart       string `xml:"nostart,attr"`
		Noend         string `xml:"noend,attr"`
		Externalbreak string `xml:"externalbreak,attr"`
		Callid        struct {
			Text   string `xml:",chardata"`
			ID     string `xml:"id,attr"`
			Typeid string `xml:"typeid,attr"`
			Native string `xml:"native,attr"`
		} `xml:"callid"`
		Parties struct {
			Text  string `xml:",chardata"`
			Count string `xml:"count,attr"`
			Party []struct {
				Text        string `xml:",chardata"`
				Callid      string `xml:"callid,attr"`
				Segmentnum  string `xml:"segmentnum,attr"`
				ID          string `xml:"id,attr"`
				State       string `xml:"state,attr"`
				Partytypeid string `xml:"partytypeid,attr"`
				Name        string `xml:"name,attr"`
				Address     string `xml:"address,attr"`
				Dirn        string `xml:"dirn,attr"`
				Owner       string `xml:"owner,attr"`
				Desc        string `xml:"desc,attr"`
			} `xml:"party"`
		} `xml:"parties"`
		Tags struct {
			Text    string `xml:",chardata"`
			Current string `xml:"current,attr"`
		} `xml:"tags"`
	} `xml:"cti"`
}

type Data struct {
	Data string `json:"text"`
	Recording `json:"info"`
}
