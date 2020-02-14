package rss

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/feeds"
)

// Data stores the live rss and the
// stored rss
type Data struct {
	live   string
	stored string
}

// Generator handles the main business logic
type Generator struct {
	URL      string
	Feed     feeds.Feed
	document *goquery.Document
	Data     Data
}

type video struct {
	GUID        string    `json:"guid"`
	NaturalKey  string    `json:"languageAgnosticNaturalKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"firstPublished"`
}

type videos []video

func (vs videos) Each(f func(int, video)) {
	for i, v := range vs {
		f(i, v)
	}
}

// GetDocument downloands the document and sets it
func (g *Generator) GetDocument(resp *http.Response) {
	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	g.document = document
}

// ParseDocument populates the feed fetching data from the document
func (g *Generator) ParseDocument() {
	g.setFeedCreationTime()

	// Find all items and process them with the functions
	// in charge to process elements
	videosDone := make(chan bool)
	newsDone := make(chan bool)

	go func(done chan bool) {
		videosFeed := g.getVideos()
		videosFeed.Each(g.processVideoElement)
		done <- true
	}(videosDone)

	go func(done chan bool) {
		newsFeed := g.document.Find("div.whatsNewItems div.synopsis")
		newsFeed.Each(g.processNewsFeedElement)
		done <- true
	}(newsDone)

	_, _ = <-videosDone, <-newsDone
}

func (g Generator) getVideos() videos {
	type category struct {
		Media []video `json:"media"`
	}

	type videos struct {
		Category category `json:"category"`
	}

	start := time.Now()
	fmt.Printf("start fetching %v\n", os.Getenv("VIDEOS_URL"))

	response, err := http.Get(os.Getenv("VIDEOS_URL"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	v := videos{}
	jsonErr := json.Unmarshal(body, &v)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	secs := time.Since(start).Seconds()
	fmt.Printf("%s request fulfilled, %.2fs elapsed\n", os.Getenv("VIDEOS_URL"), secs)

	return v.Category.Media
}

// GetRssData generates the live rss feed and sets it
// in the Data struct
func (g *Generator) GetRssData() {
	g.Feed.Sort(func(a, b *feeds.Item) bool {
		return a.Created.After(b.Created)
	})

	data, err := g.Feed.ToRss()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	g.Data.live = data
}

func (g *Generator) setFeedCreationTime() {
	g.Feed.Created = time.Now()
}

func (g *Generator) processNewsFeedElement(index int, element *goquery.Selection) {
	category := sanitize(element.Find(".contextTitle").Text())
	link := element.Find("h3 a")
	title := sanitize(link.Text())
	text := sanitize(element.Find(".desc").Text())
	pubDate := sanitize(element.Find(".pubDate").Text())
	created, err := time.Parse("2006-01-02", pubDate)
	if err != nil {
		fmt.Println(err)
	}

	// See if the href attribute exists on the element
	href, exists := link.Attr("href")
	if exists {
		link := os.Getenv("HREF_BASE_URL") + href
		g.Feed.Add(&feeds.Item{
			Title: title,
			Id:    link,
			Link: &feeds.Link{
				Href: link,
			},
			Description: category, // Descriptions is treated as a category, should not be empty
			Content:     text,
			Created:     created,
		})
	}
}

func (g *Generator) processVideoElement(index int, v video) {
	link := os.Getenv("VIDEO_BASE_PATH") + v.NaturalKey
	g.Feed.Add(&feeds.Item{
		Title: v.Title,
		Id:    v.GUID,
		Link: &feeds.Link{
			Href: link,
		},
		Description: "VIDEO", // Descriptions is treated as a category, should not be empty
		Content:     v.Description,
		Created:     v.PubDate,
	})
}

// WriteRss writes the rss xml file
func (g *Generator) WriteRss(w io.Writer) error {
	return g.Feed.WriteRss(w)
}

// GetStored fetches s3 rss xml
func (g *Generator) GetStored(done chan bool) {
	start := time.Now()
	fmt.Println("start fetching stored S3 rss data")

	svc := s3.New(session.New())
	input := &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(os.Getenv("RSS_FILENAME")),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)
	s3Rss := buf.String()
	g.Data.stored = s3Rss

	secs := time.Since(start).Seconds()
	fmt.Printf("S3 object retrieved, %.2fs elapsed\n", secs)

	done <- true
}

// Compare examines the downloaded rss and s3 stored one
// and returns a bool indicating if they are different
func (g *Generator) Compare() bool {
	return (getFirstItemTitle(g.Data.live) != getFirstItemTitle(g.Data.stored))
}

func sanitize(entry string) string {
	return html.UnescapeString(strings.TrimSpace(entry))
}

func getFirstItemTitle(rssContents string) string {
	type item struct {
		XMLName     xml.Name `xml:"item"`
		Title       string   `xml:"title"`
		Link        string   `xml:"link"`
		Description string   `xml:"description"`
		PubDate     string   `xml:"pubDate"`
	}

	type channel struct {
		XMLName     xml.Name `xml:"channel"`
		Title       string   `xml:"title"`
		Link        string   `xml:"link"`
		Description string   `xml:"description"`
		PubDate     string   `xml:"pubDate"`
		Items       []item   `xml:"item"`
	}

	type rssStruct struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}

	var rss rssStruct
	xml.Unmarshal([]byte(rssContents), &rss)

	return rss.Channel.Items[0].Title
}
