package playlist

import (
	"fmt"
	"strconv"
	"os/exec"
	"log"
	"regexp"
	"strings"
	"math"
	"math/rand"
	"os"
	"path"
	yaml "github.com/kylelemons/go-gypsy/yaml"
)

func (b *AvailableBlock) Init(t string, s string, filepaths yaml.List, u bool) *AvailableBlock {
	b.Title = t
	b.Series = s
	b.DoPublish = u
	count := filepaths.Len()
	if (count <= 0) {
		log.Fatalf("No filepaths for %s.", t)
	}	
	b.Items = make([]*os.File, count)
	for i, e := range filepaths {
		f, err := os.Open(e.(yaml.Scalar).String())
		if err != nil {
			log.Fatalf("Missing file for %s.", t)
		}
		b.Items[i] = f
	}
	b.Duration = b.GetDuration()
	return b
}

func (b *AvailableBlock) GetDuration() int {
	total := 0
	var itemPath, tmpPath string
	exp, err := regexp.Compile("Duration: ([0-9]{2}:[0-9]{2}:[0-9]{2}).[0-9]{2}")
	if err != nil {
		log.Fatalf("regexp.Compile: %v", err)
	}
	for _, f := range b.Items {
		itemPath = f.Name()
		tmpPath = path.Join( os.TempDir(), "ffmpeg" + strconv.Itoa(rand.Intn(5000000 + 1) + 1000000) + path.Ext(itemPath))
		fmt.Printf("\n%v\n", tmpPath)
		cmd := exec.Command("ffmpeg", "-i", itemPath, "-c", "copy", "-t", "1", tmpPath) // hack to get zero exit code
		stdout, er := cmd.CombinedOutput()
		if er != nil {
			log.Fatalf("cmd.CombinedOutput: %v", er)
		}
		result := exp.FindSubmatch(stdout)
		if result == nil {
			log.Fatalf("Could not determine duration")
		}
		parts := strings.Split(string(result[1]), ":")
		for i, part := range parts {
			x, _ := strconv.Atoi(part)
			val := int( math.Pow(60, float64(len(parts[i+1:]))) ) * x
			total = total + val
		}
		cleanup(tmpPath)
	}
	fmt.Printf("Using %s [%ds]\n", b.Title, total)
	return total
}

func cleanup(path string) {
	if f, err := os.Open(path); err == nil {
		if err := os.Remove(f.Name()); err != nil {
			log.Fatal(err)
		}
	}
}