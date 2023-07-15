package pkg

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"github.com/chromedp/chromedp"
)

const (
	htmlFilePath = "output/orm_chart.html"
	jpegFilePath = "output/orm_chart.jpeg"
)

func (gh GitHub) MakeChart() error {
	if err := gh.makeHTMLChartFile(); err != nil {
		return err
	}

	if err := convertHTMLToImage(); err != nil {
		return err
	}
	return nil
}

func (gh GitHub) makeHTMLChartFile() error {
	line := setUpChart()
	dates := generateDates()
	for _, d := range gh.ReadmeDetailsRepositorys {
		line = d.makeLine(line, dates)
	}
	f, err := os.Create(htmlFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	line.Render(f)
	return nil
}

func setUpChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
	)
	return line
}

func (d ReadmeDetailsRepository) makeLine(line *charts.Line, dates []string) *charts.Line {
	starHistorys := d.generateStarHistorys()
	line.SetXAxis(dates).AddSeries(d.RepoName, starHistorys)
	line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	return line
}

func (d ReadmeDetailsRepository) generateStarHistorys() []opts.LineData {
	starHistorys := make([]opts.LineData, 0)

	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount72MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount60MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount48MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount36MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount24MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCount12MouthAgo"]})
	starHistorys = append(starHistorys, opts.LineData{Value: d.StarCounts["StarCountNow"]})

	return starHistorys
}

func convertHTMLToImage() error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)
	execAllocatorOpts := append(chromedp.DefaultExecAllocatorOptions[:], options...)
	ctx, cancel = chromedp.NewExecAllocator(ctx, execAllocatorOpts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	htmlContent, err := os.Open(htmlFilePath)
	if err != nil {
		return err
	}
	defer htmlContent.Close()

	data, err := io.ReadAll(htmlContent)
	if err != nil {
		return err
	}

	var buf []byte
	if err := chromedp.Run(ctx, screenshotTask(string(data), &buf)); err != nil {
		return err
	}

	if err := saveToFile(jpegFilePath, buf); err != nil {
		return err
	}

	fmt.Println("HTML converted to JPEG successfully.")
	return nil
}

func screenshotTask(htmlContent string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("data:text/html," + htmlContent),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Screenshot("body", res, chromedp.NodeVisible, chromedp.ByQuery),
	}
}

func saveToFile(filepath string, data []byte) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
