package dates

import (
	"strconv"
	"time"
)

const (
	FmtDate              = "2006-01-02"
	FmtTime              = "15:04:05"
	FmtDateTime          = "2006-01-02 15:04:05"
	FmtDateTimeNoSeconds = "2006-01-02 15:04"
)

// NowUnix 秒时间戳
func NowUnix() int64 {
	return time.Now().Unix()
}

// FromUnix 秒时间戳转时间
func FromUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// NowTimestamp 当前毫秒时间戳
func NowTimestamp() int64 {
	return Timestamp(time.Now())
}

func NowTime() time.Time {
	return time.Unix(0, NowTimestamp()*int64(time.Millisecond))
}

// Timestamp 毫秒时间戳
func Timestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// FromTimestamp 毫秒时间戳转时间
func FromTimestamp(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}

// Format 时间格式化
func Format(time time.Time, layout string) string {
	return time.Format(layout)
}

// Parse 字符串时间转时间类型
func Parse(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// GetDay return yyyyMMdd
func GetDay(time time.Time) int {
	ret, _ := strconv.Atoi(time.Format("20060102"))
	return ret
}

// WithTimeAsStartOfDay
// 返回指定时间当天的开始时间
func WithTimeAsStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

/**
 * 将时间格式换成 xx秒前，xx分钟前...
 * 规则：
 * 59秒--->刚刚
 * 1-59分钟--->x分钟前（23分钟前）
 * 1-24小时--->x小时前（5小时前）
 * 昨天--->昨天 hh:mm（昨天 16:15）
 * 前天--->前天 hh:mm（前天 16:15）
 * 前天以后--->mm-dd（2月18日）
 */
func PrettyTime(milliseconds int64) string {
	t := FromTimestamp(milliseconds)
	duration := (NowTimestamp() - milliseconds) / 1000
	if duration < 60 {
		return "刚刚"
	} else if duration < 3600 {
		return strconv.FormatInt(duration/60, 10) + "分钟前"
	} else if duration < 86400 {
		return strconv.FormatInt(duration/3600, 10) + "小时前"
	} else if Timestamp(WithTimeAsStartOfDay(time.Now().Add(-time.Hour*24))) <= milliseconds {
		return "昨天 " + Format(t, FmtTime)
	} else if Timestamp(WithTimeAsStartOfDay(time.Now().Add(-time.Hour*24*2))) <= milliseconds {
		return "前天 " + Format(t, FmtTime)
	} else {
		return Format(t, FmtDate)
	}
}

func GetEndDayTime() time.Time {
	currentTime := time.Now()
	endTime := time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		23, 59, 59, 0,
		currentTime.Location())
	return endTime
}

func GetSurplusSec() time.Duration {
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)
	remainSecond := time.Duration(todayLastTime.Unix()-time.Now().Local().Unix()) * time.Second
	return remainSecond
}

func AddDay(time int64, num int64) int64 {
	if time == 0 {
		time = NowTimestamp()
	}
	return time + (86400000 * num)
}

func ReduceDay(time int64, num int64) int64 {
	if time == 0 {
		time = NowTimestamp()
	}
	return time - (86400000 * num)
}

// 获取当天起始和结束时间
func GetNowStartAndEndTime() (int64, int64) {
	//1.获取当前时区
	loc, _ := time.LoadLocation("Local")

	//2.今日日期字符串
	date := time.Now().Format("2006-01-02")

	//3.拼接成当天0点时间字符串
	startDate := date + " 00:00:00"
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate, loc)

	//4.拼接成当天23点时间字符串
	endDate := date + " 23:59:59"
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate, loc)

	//5.返回当天0点和23点59分的时间戳
	return Timestamp(startTime), Timestamp(endTime)
}

// 获取昨天起始和结束时间
func GetYesterdayStartAndEndTime() (int64, int64) {
	//1.获取当前时区
	loc, _ := time.LoadLocation("Local")

	//2.昨日日期字符串
	date := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	//3.拼接成当天0点时间字符串
	startDate := date + " 00:00:00"
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate, loc)

	//4.拼接成当天23点时间字符串
	endDate := date + " 23:59:59"
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate, loc)

	//5.返回当天0点和23点59分的时间戳
	return Timestamp(startTime), Timestamp(endTime)
}

func GetMonthStartEnd(t time.Time) (int64, int64) {
	monthStartDay := t.AddDate(0, 0, -t.Day()+1)
	monthStartTime := time.Date(monthStartDay.Year(), monthStartDay.Month(), monthStartDay.Day(), 0, 0, 0, 0, t.Location())
	monthEndDay := monthStartTime.AddDate(0, 1, -1)
	monthEndTime := time.Date(monthEndDay.Year(), monthEndDay.Month(), monthEndDay.Day(), 23, 59, 59, 0, t.Location())
	return Timestamp(monthStartTime), Timestamp(monthEndTime)
}
