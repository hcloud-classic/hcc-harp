package vnstat

type vnStat struct {
	Vnstatversion string `json:"vnstatversion"`
	Jsonversion   string `json:"jsonversion"`
	Interfaces    []struct {
		ID      string `json:"id"`
		Nick    string `json:"nick"`
		Created struct {
			Date struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"date"`
		} `json:"created"`
		Updated struct {
			Date struct {
				Year  int `json:"year"`
				Month int `json:"month"`
				Day   int `json:"day"`
			} `json:"date"`
			Time struct {
				Hour    int `json:"hour"`
				Minutes int `json:"minutes"`
			} `json:"time"`
		} `json:"updated"`
		Traffic struct {
			Total struct {
				Rx int `json:"rx"`
				Tx int `json:"tx"`
			} `json:"total"`
			Days []struct {
				ID   int `json:"id"`
				Date struct {
					Year  int `json:"year"`
					Month int `json:"month"`
					Day   int `json:"day"`
				} `json:"date"`
				Rx int64 `json:"rx"`
				Tx int64 `json:"tx"`
			} `json:"days"`
		} `json:"traffic"`
	} `json:"interfaces"`
}
