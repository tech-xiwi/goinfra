package exception

import (
	"runtime/debug"
	"time"
)

// Block struct to define try catch finally
// like java
type Block struct {
	Try           func(...interface{})
	Catch         func([]byte, ...interface{})
	Finally       func(...interface{})
	ReTry         func(...interface{})
	MaxRetryCount int
	RetryCount    int
	RetryTimeSpan time.Duration
}

// Do ,inner to exec try catch
func (thiz *Block) Do(params ...interface{}) {
	if thiz.Finally != nil {
		defer thiz.Finally(params...)
	}
	if thiz.Catch != nil {
		defer func(params ...interface{}) {
			if r := recover(); r != nil {
				thiz.Catch(debug.Stack(), r, params)
				if thiz.ReTry != nil && thiz.RetryCount < thiz.MaxRetryCount {
					thiz.RetryCount++
					if thiz.RetryTimeSpan > 0 {
						time.Sleep(thiz.RetryTimeSpan)
					}
					thiz.do(params...)
				}
			}
		}(params)
	}
	thiz.Try(params)
}

func (thiz *Block) do(params ...interface{}) {
	if thiz.Catch != nil {
		defer func(params ...interface{}) {
			if r := recover(); r != nil {
				thiz.RetryCount++
				thiz.Catch(debug.Stack(), r, params)
				if thiz.ReTry != nil && thiz.RetryCount < thiz.MaxRetryCount {
					if thiz.RetryTimeSpan > 0 {
						time.Sleep(thiz.RetryTimeSpan)
					}
					thiz.do(params...)
				}
			}
		}(params[0])
	}
	thiz.ReTry(params[0])
}
