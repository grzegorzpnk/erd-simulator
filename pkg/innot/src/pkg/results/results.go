package results

type Results struct {
	Failed     int `json:"relocation-failed"`
	Successful int `json:"relocation-successful"`
	Redundant  int `json:"relocation-redundant"`
}

type Client struct {
	Results Results
}

func NewClient() *Client {
	return &Client{Results{
		Failed:     0,
		Successful: 0,
		Redundant:  0,
	}}
}

func (r *Results) IncFailed() {
	r.Failed++
}

func (r *Results) IncSuccessful() {
	r.Successful++
}

func (r *Results) IncRedundant() {
	r.Redundant++
}
