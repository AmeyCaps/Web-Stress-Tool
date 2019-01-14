package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var count float64
var length int
var flag int

type pts struct {
	x float64
	y float64
}

func stop() {
	consoleReader := bufio.NewReaderSize(os.Stdin, 1)
	input, _ := consoleReader.ReadByte()

	ascii := input

	// ESC = 27 and Ctrl-C = 3
	if ascii == 27 || ascii == 3 {
		flag = 1
		return
	}
}

func Requests(url string) {
	for i := 0; i < 10; i++ {
		resp, err := http.Get(url)

		if err != nil {
			panic(err)
		}
		count += 1.0
		body, err := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		length += len(bodyStr)
		resp.Body.Close()
	}

}

func handler(wg *sync.WaitGroup, url string, rout int) {

	//var pts []plotter.XYs
	var newpts []pts
	secs := 0.0
	flag = 0

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Stress Testing"
	p.X.Label.Text = "Seconds"
	p.Y.Label.Text = "No. of http requests sent"

	go stop()

	for {
		for i := 0; i < rout; i++ {
			go Requests(url)
		}

		time.Sleep(1000 * time.Millisecond)
		if flag == 1 {
			ppts := make(plotter.XYs, len(newpts))
			for i, xy := range newpts {
				ppts[i].X = xy.x
				ppts[i].Y = xy.y
				//fmt.Println(ppts)
			}
			err = plotutil.AddLinePoints(p, "Http requests per second", ppts)
			if err != nil {
				panic(err)
			}

			// Save the plot to a PNG file.
			if err := p.Save(10*vg.Inch, 5*vg.Inch, "points.png"); err != nil {
				panic(err)
			}

			wg.Done()
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), count)
		newpts = append(newpts, pts{secs, count})
		secs += 1.0
		//fmt.Println(newpts)
		count = 0.0

	}
	//fmt.Println(len(newpts))
}

func main() {

	//variables

	var nRoutines int
	count = 0
	length = 0

	var wg sync.WaitGroup

	wg.Add(1)

	//Take user inputs
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Please enter the website url :  ")
	inputUrl, _ := reader.ReadString('\n')
	inputUrl = strings.Replace(inputUrl, "\n", "", -1)

	fmt.Println("Enter number of co routines (1-500)")
	fmt.Scan(&nRoutines)

	fmt.Println("Press ESC followed by Enter to stop.")

	go handler(&wg, inputUrl, nRoutines)

	wg.Wait()
	//time.Sleep(15000 * time.Millisecond)
	fmt.Println("Total response size", length/1024, "MB")
	fmt.Println("Done")

}
