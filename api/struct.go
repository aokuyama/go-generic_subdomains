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
	Status string `json:"status"`
	Output string `json:"output"`
}

func (r *DescribeResponse) IsCompleted() bool {
	return r.Status == "SUCCEEDED" || r.Status == "FAILED"
}

func (r *DescribeResponse) IsSucceeded() bool {
	return r.Status == "SUCCEEDED"
}

func (r *DescribeResponse) IsDescribeResponse() bool {
	return len(r.Status) > 0
}
