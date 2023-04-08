# url-fmt

Create constant URL strings that can be matched, filled, extracted, fetched, and souped all by utilising the power of Go's string interpolation verbs!

## Why

Scraping from multiple different URLs can become messy. Defining string constants with string interpolation verbs embedded into them so you can insert parameters and path segments. `url-fmt` makes it easy to define enumerations of constant URL formats that can be:

- `Fill`-ed: fill a URL format with the given arguments. Acts the same as `fmt.Sprintf`.
- `Regex`-ed: generates a regular expression pattern for the URL format by converting each string interpolation verb into its corresponding regex pattern.
- `Match`-ed: match a URL format to an already filled URL string.
- `ExtractArgs`-ed: extract the corresponding string interpolation verbs from a filled URL string.
- `Standardise`-d: extract the arguments from a filled URL string and fill the URL format with the extracted args.
- `Request`-ed: generate a `http.Request` for the given URL format.
- `Soup`-ed: make a request to the given URL format and parse the returned HTML content into a searchable BeautifulSoup-like object that can be searched. The BeautifulSoup implementation comes from Anas Khan's [soup](https://github.com/anaskhan96/soup) library.
- `JSON`-ed: make a request to the given URL format and parse the returned JSON content into a `map[string]any`.

## How

Define your URL formats for the URLs that you will access/scrape. We will be using some examples for Steam and Itch.IO scraping:

```go
package main

import "github.com/andygello555/url-fmt"

const (
	SteamAppPage urlfmt.URL = "%s://store.steampowered.com/app/%d"
	SteamAppReviews urlfmt.URL = "https://store.steampowered.com/appreviews/%d?json=1&cursor=%s&language=%s&day_range=9223372036854775807&num_per_page=%d&review_type=all&purchase_type=%s&filter=%s&start_date=%d&end_date=%d&date_range_type=%s"
	ItchIOGamePage urlfmt.URL = "http://%s.itch.io/%s"
)
```

Notice how we can provide the protocol (`https://` or `http://`), or not (`%s://`). `url-fmt` will automatically add the HTTPS protocol when filling (this won't interfere with the arguments that you provide), and generate the following regex when `Regex` is called: `https?`.

Then you can use these however you require:

```go
package main

import (
	"fmt"
	"github.com/andygello555/url-fmt"
)

const (
	SteamAppPage    urlfmt.URL = "%s://store.steampowered.com/app/%d"
	SteamAppReviews urlfmt.URL = "https://store.steampowered.com/appreviews/%d?json=1&cursor=%s&language=%s&day_range=9223372036854775807&num_per_page=%d&review_type=all&purchase_type=%s&filter=%s&start_date=%d&end_date=%d&date_range_type=%s"
	ItchIOGamePage  urlfmt.URL = "http://%s.itch.io/%s"
)

func main() {
	steamPageURLs := []string{
		"http://store.steampowered.com/app/477160",
		"http://store.steampowered.com/app/477160/Human_Fall_Flat/",
		"https://store.steampowered.com/app/477160",
		"https://store.steampowered.com/app/Human_Fall_Flat/",
	}

	// Check to see which URLs above match the schema of the SteamAppPage URL format
	for _, steamPageURL := range steamPageURLs {
		fmt.Printf("%s matches %s = %t\n", steamPageURL, SteamAppPage.Regex(), SteamAppPage.Match(steamPageURL))
	}

	itchIOPageURLs := []string{
		"https://hempuli.itch.io/baba-files-taxes",
		"https://sokpop.itch.io/ballspell",
    }

	// Extract the arguments from the above Itch.IO game page URLs
	for _, itchIOPageURL := range itchIOPageURLs {
		fmt.Printf("extracted %v from %s\n", ItchIOGamePage.ExtractArgs(itchIOPageURL), itchIOPageURL)
    }

	// Fetch the name for Steam App 477160 via its HTML Steam store page
	if soup, _, err := SteamAppPage.Soup(nil, 477160); err != nil {
		fmt.Printf("Could not get soup for %s, because %s\n", SteamAppPage.Fill(477160), err.Error())
	} else {
		fmt.Println("Name of Steam", soup.Find("div", "id", "appHubAppName").Text())
	}

	// Fetch the number of positive and negative reviews for Steam app 477160 via the Steam store's JSON API
	args := []any{477160, "*", "all", 20, "all", "all", -1, -1, "all"}
	if j, _, err := SteamAppReviews.JSON(nil, args...); err != nil {
		fmt.Printf("Could not get reviews for %s, because %s\n", SteamAppReviews.Fill(args...), err.Error())
	} else {
		fmt.Printf("Query summary for Steam app 477160: %v\n", j["query_summary"])
	}
}
```


