package urlfmt

import (
	"fmt"
	"sort"
)

func ExampleURL_Regex() {
	const (
		SteamAppPage   URL = "%s://store.steampowered.com/app/%d"
		ItchIOGamePage URL = "%s://%s.itch.io/%s"
	)

	fmt.Println(SteamAppPage.Regex())
	fmt.Println(ItchIOGamePage.Regex())
	// Output:
	// https?://store.steampowered.com/app/(\d+)
	// https?://([a-zA-Z0-9-._~]+).itch.io/([a-zA-Z0-9-._~]+)
}

func ExampleURL_Match() {
	const (
		SteamAppPage   URL = "%s://store.steampowered.com/app/%d"
		ItchIOGamePage URL = "%s://%s.itch.io/%s"
	)

	fmt.Println(SteamAppPage.Match("http://store.steampowered.com/app/477160"))
	fmt.Println(SteamAppPage.Match("http://store.steampowered.com/app/477160/Human_Fall_Flat/"))
	fmt.Println(SteamAppPage.Match("https://store.steampowered.com/app/477160"))
	fmt.Println(SteamAppPage.Match("https://store.steampowered.com/app/Human_Fall_Flat/"))
	fmt.Println(ItchIOGamePage.Match("https://hempuli.itch.io/baba-files-taxes"))
	fmt.Println(ItchIOGamePage.Match("https://sokpop.itch.io/ballspell"))
	// Output:
	// true
	// true
	// true
	// false
	// true
	// true
}

func ExampleURL_ExtractArgs() {
	const (
		SteamAppPage   URL = "%s://store.steampowered.com/app/%d"
		ItchIOGamePage URL = "%s://%s.itch.io/%s"
	)

	fmt.Println(SteamAppPage.ExtractArgs("http://store.steampowered.com/app/477160"))
	fmt.Println(SteamAppPage.ExtractArgs("https://store.steampowered.com/app/477160/Human_Fall_Flat/"))
	fmt.Println(ItchIOGamePage.ExtractArgs("https://hempuli.itch.io/baba-files-taxes"))
	fmt.Println(ItchIOGamePage.ExtractArgs("https://sokpop.itch.io/ballspell"))
	// Output:
	// [477160]
	// [477160]
	// [hempuli baba-files-taxes]
	// [sokpop ballspell]
}

func ExampleURL_Soup() {
	const SteamAppPage URL = "%s://store.steampowered.com/app/%d"
	fmt.Printf("Getting name of app 477160 from %s:\n", SteamAppPage.Fill(477160))
	if soup, _, err := SteamAppPage.Soup(nil, 477160); err != nil {
		fmt.Printf("Could not get soup for %s, because %s", SteamAppPage.Fill(477160), err.Error())
	} else {
		fmt.Println(soup.Find("div", "id", "appHubAppName").Text())
	}
	// Output:
	// Getting name of app 477160 from https://store.steampowered.com/app/477160:
	// Human: Fall Flat
}

func ExampleURL_JSON() {
	const SteamAppReviews URL = "%s://store.steampowered.com/appreviews/%d?json=1&cursor=%s&language=%s&day_range=9223372036854775807&num_per_page=%d&review_type=all&purchase_type=%s&filter=%s&start_date=%d&end_date=%d&date_range_type=%s"
	args := []any{477160, "*", "all", 20, "all", "all", -1, -1, "all"}
	fmt.Printf("Getting review stats for 477160 from %s:\n", SteamAppReviews.Fill(args...))
	if j, _, err := SteamAppReviews.JSON(nil, args...); err != nil {
		fmt.Printf("Could not get reviews for %s, because %s", SteamAppReviews.Fill(args...), err.Error())
	} else {
		var i int
		keys := make([]string, len(j["query_summary"].(map[string]any)))
		for key := range j["query_summary"].(map[string]any) {
			keys[i] = key
			i++
		}
		sort.Strings(keys)
		fmt.Println(keys)
	}
	// Output:
	// Getting review stats for 477160 from https://store.steampowered.com/appreviews/477160?json=1&cursor=*&language=all&day_range=9223372036854775807&num_per_page=20&review_type=all&purchase_type=all&filter=all&start_date=-1&end_date=-1&date_range_type=all:
	// [num_reviews review_score review_score_desc total_negative total_positive total_reviews]
}
