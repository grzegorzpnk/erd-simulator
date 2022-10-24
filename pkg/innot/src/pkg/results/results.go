package results

type Results struct {
	Failed     int `json:"relocation-failed"`
	Successful int `json:"relocation-successful"`
	Redundant  int `json:"relocation-redundant"`
	Skipped    int `json:"relocation-skipped"`
}

type Client struct {
	Results Results
}

func NewClient() *Client {
	return &Client{Results{
		Failed:     0,
		Successful: 0,
		Redundant:  0,
		Skipped:    0,
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

func (r *Results) IncSkipped() {
	r.Skipped++
}

func (r *Results) ResetCounter() {
	r.Redundant = 0
	r.Skipped = 0
	r.Successful = 0
	r.Failed = 0
}
