package testdata

import (
	"bufio"
	"fmt"
	"github.com/lcserny/goutils"
	"math/rand"
	"os"
	"strings"
)

var (
	pattern      = "%s;%s%d;%d###%s;%s\n"
	normalized   = "A text here of normalized name"
	typesSlice   = []string{"TV", "MOVIE"}
	randOutNames = []string{"Something Somewhere (2018)", "Hello There", "A Movie Title (2015)", "Aliens (1997)"}
	randYears    = []int{1997, 2008, 2009, 2010, 2018}
)

func main() {
	// generateCacheFiles()

	bench("testdata/bench10.cache", "MOVIE;A text here of normalized name6;2009")
	bench("testdata/bench100.cache", "TV;A text here of normalized name79;2018")
	bench("testdata/bench1000.cache", "MOVIE;A text here of normalized name950;2008")
}

func bench(fileName, prefix string) {
	startNano := goutils.MakeTimestamp()
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	scanner := bufio.NewScanner(file)
	counter := 0
	for scanner.Scan() {
		counter++
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			println(line[len(prefix)+3:])
		}
	}
	println("Lines:", counter)
	_ = file.Close()

	endNano := goutils.MakeTimestamp() - startNano
	println("It took (in ms):", endNano)
}

func generateCacheFiles() {
	file1, _ := os.Create("testdata/bench10.cache")
	defer file1.Close()
	file2, _ := os.Create("testdata/bench100.cache")
	defer file2.Close()
	file3, _ := os.Create("testdata/bench1000.cache")
	defer file3.Close()
	for i := 0; i < 10; i++ {
		_, _ = file1.WriteString(generateRandomCacheKey(i))
	}
	for i := 0; i < 100; i++ {
		_, _ = file2.WriteString(generateRandomCacheKey(i))
	}
	for i := 0; i < 1000; i++ {
		_, _ = file3.WriteString(generateRandomCacheKey(i))
	}
}

func generateRandomCacheKey(i int) string {
	return fmt.Sprintf(pattern,
		typesSlice[rand.Intn(len(typesSlice))],
		normalized,
		i,
		randYears[rand.Intn(len(randYears))],
		randOutNames[rand.Intn(len(randOutNames))],
		randOutNames[rand.Intn(len(randOutNames))],
	)
}
