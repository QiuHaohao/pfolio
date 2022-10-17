package clock

import "time"

type NowFn func() time.Time

var nowFn NowFn

func init() {
	nowFn = time.Now
}

func Now() time.Time {
	return nowFn()
}

func SetNowFn(fn NowFn) {
	nowFn = fn
}
