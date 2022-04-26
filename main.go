package googleanalytics

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	ga "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
	"net/http"
	"time"
)

func Set(e *echo.Echo, clientJSONPath string) {
	e.GET("/ga", func(c echo.Context) error {

		ctx := context.Background()
		client, err := ga.NewService(ctx, option.WithCredentialsFile(clientJSONPath))

		if err != nil {
			panic(err)
		}

		now := time.Now()

		requests := []*ga.RunReportRequest{}
		today := &ga.RunReportRequest{
			DateRanges: []*ga.DateRange{
				{StartDate: "yesterday", EndDate: "today"},
			},
			Metrics: []*ga.Metric{
				{Name: "screenPageViews"},
			},
		}
		weekDay := now.Weekday().String()
		thisWeek := now

		if weekDay == "Tuesday" {
			thisWeek = now.AddDate(0, 0, -1)
		} else if weekDay == "Wednesday" {
			thisWeek = now.AddDate(0, 0, -2)
		} else if weekDay == "Thursday" {
			thisWeek = now.AddDate(0, 0, -3)
		} else if weekDay == "Friday" {
			thisWeek = now.AddDate(0, 0, -4)
		} else if weekDay == "Saturday" {
			thisWeek = now.AddDate(0, 0, -5)
		} else if weekDay == "Sunday" {
			thisWeek = now.AddDate(0, 0, -6)
		}

		week := &ga.RunReportRequest{

			DateRanges: []*ga.DateRange{

				{StartDate: thisWeek.Format("2006-01-02"), EndDate: "today"},
			},
			Metrics: []*ga.Metric{
				{Name: "screenPageViews"},
			},
		}

		month := &ga.RunReportRequest{
			DateRanges: []*ga.DateRange{
				{StartDate: fmt.Sprintf("%s-01", now.Format("2006-01")), EndDate: "today"},
			},

			Metrics: []*ga.Metric{
				{Name: "screenPageViews"},
			},
		}
		year := &ga.RunReportRequest{

			DateRanges: []*ga.DateRange{
				{StartDate: fmt.Sprintf("%s-01-01", now.Format("2006")), EndDate: "today"},
			},

			Metrics: []*ga.Metric{
				{Name: "screenPageViews"},
			},
		}
		allTime := &ga.RunReportRequest{

			DateRanges: []*ga.DateRange{
				{StartDate: "2022-04-01", EndDate: "today"},
			},

			Metrics: []*ga.Metric{
				{Name: "screenPageViews"},
			},
		}

		requests = append(requests, today)
		requests = append(requests, week)
		requests = append(requests, month)
		requests = append(requests, year)
		requests = append(requests, allTime)

		todayData := "0"
		weekData := "0"
		monthData := "0"
		yearData := "0"
		allTimeData := "0"

		r, _ := client.Properties.BatchRunReports("properties/311692779", &ga.BatchRunReportsRequest{
			Requests: requests,
		}).Do()

		if len(r.Reports) >= 5 {
			if len(r.Reports[0].Rows) >= 1 {
				todayData = r.Reports[0].Rows[0].MetricValues[0].Value
			}
			if len(r.Reports[1].Rows) >= 1 {
				weekData = r.Reports[1].Rows[0].MetricValues[0].Value
			}
			if len(r.Reports[2].Rows) >= 1 {
				monthData = r.Reports[2].Rows[0].MetricValues[0].Value
			}
			if len(r.Reports[3].Rows) >= 1 {
				yearData = r.Reports[3].Rows[0].MetricValues[0].Value
			}
			if len(r.Reports[4].Rows) >= 1 {
				allTimeData = r.Reports[4].Rows[0].MetricValues[0].Value
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"today":   todayData,
			"week":    weekData,
			"month":   monthData,
			"year":    yearData,
			"allTime": allTimeData,
		})
	})
}
