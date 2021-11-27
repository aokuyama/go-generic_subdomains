package api

type StartApiResponse struct {
	ExecutionArn string `json:"executionArn"`
}

type DescribeRequest struct {
	ExecutionArn string `json:"executionArn"`
}

func NewDescribeRequest(a string) *DescribeRequest {
	d := DescribeRequest{
		ExecutionArn: a,
	}
	return &d
}

type DescribeResponse struct {
	Status string      `json:"status"`
	Output interface{} `json:"output"`
}

func (r *DescribeResponse) isCompleted() bool {
	return r.Status == "SUCCEEDED" || r.Status == "FAILED"
}
