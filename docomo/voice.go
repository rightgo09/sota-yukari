package docomo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const url = "https://api.apigw.smt.docomo.ne.jp/aiTalk/v1/textToSpeech?APIKEY="
const xml = `<?xml version="1.0" encoding="utf-8" ?>
<speak version="1.1">
	<voice name="%s">
		<prosody rate="%.2f" pitch="%.2f" volume="%.2f" range="%.2f">%s</prosody>
	</voice>
</speak>`
const (
	yukari = "sumire"
	maki   = "maki"
)

type Speaker struct {
	Name    string
	Rate    float32
	Pitch   float32
	Volume  float32
	Range   float32
	Prosody string
}

type Client struct {
	apikey     string
	httpClient *http.Client
	Q          chan string
}

func NewClient(apikey string) *Client {
	return &Client{
		apikey:     apikey,
		httpClient: &http.Client{Timeout: time.Duration(10) * time.Second},
		Q:          make(chan string, 10),
	}
}

func Yukari(prosody string) *Speaker {
	return &Speaker{
		Name:    yukari,
		Rate:    1.2,
		Pitch:   1.1,
		Volume:  2.0,
		Range:   1.2,
		Prosody: prosody,
	}
}

func Maki(prosody string) *Speaker {
	return &Speaker{
		Name:    maki,
		Rate:    1.2,
		Pitch:   1.0,
		Volume:  1.0,
		Range:   1.0,
		Prosody: prosody,
	}
}

func newXML(s *Speaker) string {
	return fmt.Sprintf(xml, s.Name, s.Rate, s.Pitch, s.Volume, s.Range, s.Prosody)
}

func (c *Client) Synthesize(speaker *Speaker, filepath string) error {
	xml := newXML(speaker)
	req, err := http.NewRequest("POST", url+c.apikey, strings.NewReader(xml))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/ssml+xml")
	req.Header.Add("Accept", "audio/L16")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, res.Body); err != nil {
		return err
	}

	return nil
}
