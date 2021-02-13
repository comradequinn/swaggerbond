package index

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"
)

//Start populates the index from swagger files found in the specified swaggerFilesDir and then commences polling for at the specified interval
func Start(interval time.Duration, swaggerFilesDir string) {
	type (
		swaggerInfoDTO struct {
			Description string `json:"description"`
			Version     string `json:"version"`
			Title       string `json:"title"`
		}
		swaggerTagDTO struct {
			Name string `json:"name"`
		}
		swaggerFileDTO struct {
			Info swaggerInfoDTO  `json:"info"`
			Tags []swaggerTagDTO `json:"tags"`
		}
	)

	refreshIndex := func(lastRefreshed time.Time) {
		files := modifiedFiles(lastRefreshed, swaggerFilesDir)

		if len(files) > 0 {
			c, wg := make(chan *Summary, 1), sync.WaitGroup{}
			wg.Add(len(files))

			for _, f := range files {
				go func(f os.FileInfo) {
					defer wg.Done()

					var b []byte
					var e error

					p := path.Join(swaggerFilesDir, f.Name())

					if b, e = ioutil.ReadFile(p); e != nil {
						log.Printf("unable to parse %q. %v", p, e)
					}

					dto := swaggerFileDTO{}

					if e = json.NewDecoder(bytes.NewReader(b)).Decode(&dto); e != nil {
						log.Printf("the file %q was not a valid swagger file. %v", p, e)
						return
					}

					summary := Summary{Slug: f.Name(), Title: dto.Info.Title, Description: dto.Info.Description, Version: dto.Info.Version, Tags: make([]string, len(dto.Tags))}
					summary.Slug = strings.Replace(summary.Slug, path.Ext(f.Name()), "", 1)

					for i, tag := range dto.Tags {
						summary.Tags[i] = tag.Name
					}

					c <- &summary
				}(f)
			}

			go func() { wg.Wait(); close(c) }()

			summaries := []*Summary{}

			for summary := range c {
				summaries = append(summaries, summary)
			}

			log.Printf("updating the index with %v swagger file summaries", len(summaries))

			updates <- summaries
		}
	}

	lastRefreshed := time.Time{}

	refreshIndex(lastRefreshed)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic while maintaining index: %s", err)
			}
		}()

		for {
			select {
			case _ = <-time.After(interval):
				refreshStarted := time.Time{}
				refreshIndex(lastRefreshed)
				lastRefreshed = refreshStarted
			}
		}
	}()
}

func modifiedFiles(lastChangeTime time.Time, swaggerFilesDir string) []os.FileInfo {
	log.Printf("checking for swagger file changes since %v", lastChangeTime)

	var err error
	dir := []os.FileInfo{}

	dir, err = ioutil.ReadDir(swaggerFilesDir)

	if err != nil {
		log.Printf("index sync process was unable to ascertain last modification date of files. %s\n", err)
		return nil
	}

	files, changed := []os.FileInfo{}, false

	for _, file := range dir {
		if path.Ext(file.Name()) == ".json" {
			files = append(files, file)

			ctim := file.Sys().(*syscall.Stat_t).Ctim //ModTime() isn't appropriate as it isn't updated when a file is copied from one place to another
			ct := time.Unix(ctim.Sec, ctim.Nano())

			if ct.Sub(lastChangeTime) > 0 {
				changed = true
				log.Printf("swagger file %q was changed at %v", file.Name(), lastChangeTime)
			}
		}
	}

	if changed {
		log.Printf("swagger file directory contents last modified at %v. directory contains %v json file(s)", lastChangeTime, len(files))
		return files
	}

	return nil
}
