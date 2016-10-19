package schedule

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strings"
	"strconv"
)

const scheduleURL = "http://web.mh-xr.jp/schedule/index"
const eventScheduleBlockClassName = "label-main"
const dateBlockClassName = "label-wood"
const eventBlockClassName = "bg-paper text-center"
const eventNowBlockClassName = "bg-paper text-center now"

type EventSchedule struct {
	// 日付文字列
	// ex)  "2016/10/13 (木)"
	dateString string

	// イベントリスト
	events []Event
}

type Event struct {
	// イベントバナーURL
	// ex)  "//dl.mh-xr.jp/web/image/banner/asev/32ed71d517ff03f067cd2c8300ae7324632b948bc2b8c801dc3882917de8e5f3.png?1475041224"
	imageURL string

	// 開催時間リスト
	// 各時間は「00:00 〜 24:00」の形式
	openTimeRanges []string

	// 24時間タイムテーブル
	// 開催時間ならtrue
	timeTable [24]bool
}

func GetEventSchedule() (eventSchedules []EventSchedule, e error) {
	// 戻り値の初期化
	e = nil

	// ページの取得
	doc, err := goquery.NewDocument(scheduleURL)
	if err != nil {
		e = fmt.Errorf("%s is not found.", scheduleURL)
	}

	// #schedule要素の取得。この下に目的のデータがある。
	scheduleBlock := doc.Find("#schedule")
	if scheduleBlock.Length() == 0 {
		// #schedule ブロックが見つからない
		e = fmt.Errorf("#schedule not found.")
	}

	// 各日付とイベント群のスクレイピング
	// 日付のdivは class="label-wood"、イベントのdivは class="bg-paper text-center" が設定されているので、これを利用。
	// ※ class="bg-paper text-center now" は、当日のイベント（開催時間以外のイベントも含む）。
	// 日付と各イベントは同じ階層のブロック。
	// 日付ブロックの後ろにイベントブロックがn個続く形。新たに日付ブロックが出現したら、その下は次の日のイベントとなる。
	eventSchedule := EventSchedule{}
	scheduleBlock.Children().Each(func(_ int, s *goquery.Selection) {
		class, _ := s.Attr("class")

		// 日付取得
		switch class {
			case eventScheduleBlockClassName: // 「イベントスケジュール」バナー
				// 何もしない
			case dateBlockClassName: // 日付ブロック
				// 新しい日付が出現したらEventScheduleを新しく作成
				if (eventSchedule.dateString != "") {
					eventSchedules = append(eventSchedules, eventSchedule)
				}
				eventSchedule = EventSchedule{}

				// 日付のセット
				eventSchedule.dateString = s.Text()
			case eventBlockClassName: // イベントブロック
				fallthrough
			case eventNowBlockClassName: // 当日のイベントブロック
				event := Event{}

				// 開催時間の各ブロックは<div class="inner">の中にある
				innerBlock := s.Children()
				if !innerBlock.HasClass("inner") {
					// <div class="inner">が無い
					e = fmt.Errorf("No event block inner.")
				}

				// イベントバナーの取得
				img := innerBlock.ChildrenFiltered("img")
				if img.Length() == 0 {
					// イベントバナーが無い
					e = fmt.Errorf("No event image.")
				}
				event.imageURL, _ = img.Attr("src")

				// イベント開催時間の取得
				// 各開催時間ブロックは、class="bg-text-time margin-m font-red relative"が設定されているので、これを利用。
				timeBlocks := innerBlock.Find("div[class='bg-text-time margin-m font-red relative']")
				if timeBlocks.Length() == 0 {
					// 開催時間が１つもない
					e = fmt.Errorf("holding hours dose not found.")
				}
				openTimeRanges := []string{}
				timeBlocks.Each(func(_ int, s *goquery.Selection) {
					openTimeRange := s.Text()
					openTimeRanges = append(openTimeRanges, openTimeRange)
				})
				event.openTimeRanges = openTimeRanges
				// 開催時間から24時間タイムテーブルを作成
				event.timeTable = setTimeTable(event.openTimeRanges)

				eventSchedule.events = append(eventSchedule.events, event)
			default:
				// 想定外のclassが出現した
				e = fmt.Errorf("appeared class unexpected : %s", class)
		}
	})

	return
}

func setTimeTable(openTimeRanges []string) (timeTable [24]bool) {
	for _, openTimeRange := range openTimeRanges {
		times := strings.Split(openTimeRange, " 〜 ")
		startTime := times[0] // ex) "00:00"
		endTime := times[1] // ex) "24:00"
		startHours, _ := strconv.Atoi(strings.Split(startTime, ":")[0])
		endHours, _ := strconv.Atoi(strings.Split(endTime, ":")[0])
		for i := startHours; i < endHours; i++ {
			timeTable[i] = true
		}
	}
	return
}

func CreateHTML() (string, error) {
	eventSchedules, err := GetEventSchedule()
	if err != nil {
		return "", err
	}

	timeTables := [][24][]string{}
	for _, eventSchedule := range eventSchedules {
		timeTable := [24][]string{}
		//for i := 0; i < 24; i++ {
			for _, event := range eventSchedule.events {
				for i, flag := range event.timeTable {
					if (flag) {
						timeTable[i] = append(timeTable[i], event.imageURL)
					}
				}
			}
		//}
		timeTables = append(timeTables, timeTable)
	}

	html := ""
	for j, timeTable := range timeTables {
		html += "<h1>" + eventSchedules[j].dateString + "</h1>"
		html += "<table>"
		for i, times := range timeTable {
			html += "<tr>"
			html += "<td>" + strconv.Itoa(i) + "</td>"
			for _, imgOrEmpty := range times {
				html += "<td>"
				if (imgOrEmpty != "") {
					html += "<img src=http://" + imgOrEmpty + " width=100>"
				}
				html += "</td>"
			}
			html += "</tr>"
		}
		html += "</table>"
	}

	return html, nil
}

func CreateHorizonHTML() (string, error) {
	eventSchedules, err := GetEventSchedule()
	if err != nil {
		return "", err
	}

	html := ""
	indent := "    "
	for _, eventSchedule := range eventSchedules {
		html += "<table>\n"
		html += indent + "<tr>\n"
		for i := 0; i < 24; i++ {
			html += indent + indent + "<td>" + strconv.Itoa(i) + "</td>\n"
		}
		html += indent + "</tr>\n"
		for _, event := range eventSchedule.events {
			html += indent + "<tr>\n"
			for _, flag := range event.timeTable {
				html += indent + indent + "<td>\n"
				if (flag) {
					html += indent + indent + indent +"<img src=http:" + event.imageURL + " width=100>\n"
				}
				html += indent + indent + "</td>\n"
			}
			html += indent + "</tr>\n"
		}
		html += "</table>\n"
	}

	return html, nil
}
