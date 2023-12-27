package model

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

const (
	htmlFilePath = "output/orm_chart.html"
	jpegFilePath = "output/orm_chart.jpeg"
)

func (rds RepositoryDetails) MakeHTMLChartFile() error {
	line := setUpChart()
	for _, d := range rds {
		line = d.makeLine(line, generateDates())
	}

	f, err := os.Create(htmlFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	return line.Render(f)
}

func ConvertHTMLToImage() error {
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

	return saveToFile(jpegFilePath, buf)
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
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	_, err = file.Write(data)

	return err
}

func setUpChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
	)
	return line
}
