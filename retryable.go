package async

import "time"

type Retryable struct {
	Operation   IOperation
	Attempts    int
	MaxAttempts int // default MaxAttempts is 10
	Err         error
}

func (r *Retryable) Execute() (val interface{}, err error) {
	success := false

	if r.MaxAttempts == 0 {
		r.MaxAttempts = 10
	}

	for !success && r.Attempts < r.MaxAttempts-1 {
		val, err = r.Operation.Execute()
		if err == nil {
			success = true
			break
		} else {
			r.Attempts++
			// exponential backoff algorithm
			ms := GetExponentialBackoffMs(r.Attempts)

			time.Sleep(ms)
			continue
		}
	}
	return
}

// This algorithm is recommended by AWS for generating exponential backoff delay times.  It is
// transcribed from http://docs.aws.amazon.com/general/latest/gr/api-retries.html
func GetExponentialBackoffMs(attempts int) time.Duration {
	return time.Duration(((2 ^ attempts) * 50)) * time.Millisecond
}
