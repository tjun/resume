package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	output := flag.String("o", "resume.pdf", "Output PDF file path")
	dir := flag.String("d", "public", "Directory to serve")
	port := flag.String("p", "8080", "Port to serve on")
	flag.Parse()

	// Start file server
	go func() {
		fs := http.FileServer(http.Dir(*dir))
		http.Handle("/", fs)
		log.Printf("Serving %s on http://localhost:%s\n", *dir, *port)
		if err := http.ListenAndServe(":"+*port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Generate PDF
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel2 := chromedp.NewContext(allocCtx)
	defer cancel2()

	var buf []byte
	url := fmt.Sprintf("http://localhost:%s/", *port)

	log.Printf("Generating PDF from %s...\n", url)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				Do(ctx)
			return err
		}),
	); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(*output, buf, 0644); err != nil {
		log.Fatal(err)
	}

	log.Printf("PDF saved to %s\n", *output)
}
