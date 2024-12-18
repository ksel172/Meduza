package utils

type StatusS struct {
	SUCCESS string
	ERROR   string
	FAILED  string
}

var Status = StatusS{
	SUCCESS: "success",
	ERROR:   "error",
	FAILED:  "failed",
}
