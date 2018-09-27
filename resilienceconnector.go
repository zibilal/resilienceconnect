// HttpCoonnect wraps the process to connect with an external api or resources
// Resources mentioned should be on http protocol

package resilienceconnect

import (
	"errors"
	"time"
)

const (
	IsBackingOff   = "is backingoff"
	IsRestrying    = "is restarting"
	Retry          = "restart"
	Wait           = "wait"
	ConnectorFunc  = "connector"
	DefaultRestart = 3
	DefaultWait    = 1
)

// ResilienceConnect actually a type wrapper for httpClient.
// This type is compiled to encapsulate
// the ability to choose between fix retry  connection and back off retry connection
type ResilienceConnect struct {
}

// NewResilienceConnector defined the new ResilienceConnect with particular timeout connection
// timeout is the time out value on second, default is 5 seconds
func NewResilienceConnector() *ResilienceConnect {
	return new(ResilienceConnect)
}

// Connect is method that connect with external resource defined in url
// url contain the valid path to the external resource
// request is a Requestor object
func (h *ResilienceConnect) Connect(url string, request Requestor, options ConnectorOption, output interface{}) error {
	var (
		isBackoff   bool
		isRetry     bool
		retry       int
		wait        int
		connectFunc ConnectionFunc
	)

	isBackoff, _ = options.Get(IsBackingOff).(bool)
	isRetry, _ = options.Get(IsRestrying).(bool)
	retry, _ = options.Get(Retry).(int)
	wait, _ = options.Get(Wait).(int)
	connectFunc, _ = options.Get(ConnectorFunc).(ConnectionFunc)

	if connectFunc == nil {
		return errors.New(ConnectorFunc + " is required")
	}

	if isBackoff || isRetry {
		if retry == 0 {
			retry = DefaultRestart
		}
		if wait == 0 {
			wait = DefaultWait
		}

		errChannel := make(chan error)
		go func(url string, request Requestor, output interface{}, errChannel chan<- error) {
			var err error
			for i := 0; ((i < retry) && isRetry) || isBackoff; i++ {
				err = connectFunc(request, output)
				if err != nil {
					time.Sleep(time.Duration(wait) * time.Second)
					if isBackoff {
						wait = backingOff(i, wait)
					}
				} else {
					break
				}
			}
			errChannel <- err
		}(url, request, output, errChannel)

		select {
		case err := <-errChannel:
			if err != nil {
				return err
			}
			close(errChannel)
		}

		return nil
	} else {
		return connectFunc(request, output)
	}

	return nil
}

// backingOff counting the backing off value on wait parameter
// the simplest way to backing off is by mult
func backingOff(count int, wait int) int {
	if count%wait == 0 {
		return wait * 2
	} else {
		return wait
	}
}
