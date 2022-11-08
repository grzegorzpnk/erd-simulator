package results

type Results struct {
	Failed     ResultCounter   `json:"relocation-failed"`
	Successful ResultCounter   `json:"relocation-successful"`
	Redundant  ResultCounter   `json:"relocation-redundant"`
	Skipped    ResultCounter   `json:"relocation-skipped"`
	EvalTimes  EvaluationTimes `json:"evaluation-times"`
}

type ResultCounter map[string]int

// EvaluationTimes times should be in milliseconds
type EvaluationTimes struct {
	Failed     []int `json:"failed"`
	Successful []int `json:"successful"`
	Redundant  []int `json:"redundant"`
	Skipped    []int `json:"skipped"`
}

type Client struct {
	Results Results
}

func NewClient() *Client {
	return &Client{Results{
		Failed:     ResultCounter{},
		Successful: ResultCounter{},
		Redundant:  ResultCounter{},
		Skipped:    ResultCounter{},
		EvalTimes: EvaluationTimes{
			Failed:     []int{},
			Successful: []int{},
			Redundant:  []int{},
			Skipped:    []int{},
		},
	}}
}

func (r *Results) IncFailed(t string) {
	r.Failed[t]++
}

func (r *Results) IncSuccessful(t string) {
	r.Successful[t]++
}

func (r *Results) IncRedundant(t string) {
	r.Redundant[t]++
}

func (r *Results) IncSkipped(t string) {
	r.Skipped[t]++
}

func (r *Results) AddFailedTime(t int) {
	r.EvalTimes.Failed = append(r.EvalTimes.Failed, t)
}

func (r *Results) AddSuccessfulTime(t int) {
	r.EvalTimes.Successful = append(r.EvalTimes.Successful, t)
}

func (r *Results) AddRedundantTime(t int) {
	r.EvalTimes.Redundant = append(r.EvalTimes.Redundant, t)
}

func (r *Results) AddSkippedTime(t int) {
	r.EvalTimes.Skipped = append(r.EvalTimes.Skipped, t)
}

func (r *Results) Reset() {
	r.Redundant = ResultCounter{}
	r.Skipped = ResultCounter{}
	r.Successful = ResultCounter{}
	r.Failed = ResultCounter{}
	r.EvalTimes = EvaluationTimes{
		Failed:     []int{},
		Successful: []int{},
		Redundant:  []int{},
		Skipped:    []int{},
	}
}
