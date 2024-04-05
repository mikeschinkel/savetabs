package augment

import (
	"regexp"
)

// privateIPRegexp matches a URL to a private network IP address
// See: https://en.wikipedia.org/wiki/Private_network
var privateIPRegexp = `^https?://(172\.(1[6-9]|2[0-9]|3[0-1])\.|192\.168\.|10\.)`

var keywordPatterns = []struct {
	keyword string
	regexp  string
}{
	{"private-ip", privateIPRegexp},
	{"video", `(/videos?[/.])`},
	{"localhost", `^https?://(localhost|127\.0\.0\.1)`},
	{"web-app", `^https?://app\.`},
	{"q&a", `^https://(stackoverflow|superuser|quora|reddit)\.com/`},
	{"q&a", `^https://\w+\.stackexchange\.com/`},
	{"admin", `^https?://(app|(my)?accounts?)\.`},
	{"blog", `^https?://(blog\.|medium\.com)`},
	{"news", `^https?://(news\.)`},
	{"support", `^https?://(support\.)`},
	{"social-media", `^https?://(twitter|x|facebook|linkedin|reddit)\.com`},
	{"app-store", `^https?://app\.`},
	{"knowledge-base", `^https?://kb\.`},
	{"appstore", `^https?://apps\.apple\.com`},
	{"appstore", `^https?://chromewebstore\.google\.com`},
	{"auth", `^https?://appleid\.apple\.com`},
	{"auth", `/(signin|login|o?auth)/?`},
	{"messaging", `/messaging/`},
	{"repo", `^https://github.com/([^/]+)/([^/]+?)([?#].+)?$`},
	{"org", `^https://github.com/([^/]+?)([?#].+)?$`},
	{"javascript", `^https://unpkg\.com/`},
	{"shopping", `^https://www\.(amazon|bestbuy|newegg)\.com`},
	{"shopping", `^https://www\.(aliexpress)\.us`},
	{"search-results", `^https://www\.(bing|google)\.com/search`},
	{"search-results", `^https://www\.youtube\.com/results?search_query=`},
	{"financial", `^https://www\.(chase|capitalone)\.com`},
	{"events", `^https://www\.(eventbrite)\.com`},
	{"career", `^https://www\.(dice)\.com`},
	{"user-profile", `^https://www\.linkedin\.com/in/`},
	{"shipping", `^https://www\.(fedex|ups|usps)\.com`},
	{"shipment-tracking", `^https://www\.fedex\.com/fedextrack/?trknbr=`},
	{"shipment-tracking", `^https://www\.ups\.com/track?track=`},
	{"web-hosting", `^https://www\.(digitalocean)\.com`},
	{"web-maps", `^https://www\.google\.com/maps/@`},
	{"computer-hardware", `^https://www\.(dell)\.com`},
	{"computer-software", `^https://www\.(docker|microsoft)\.com`},
}
var keywordRegexps = make([]*regexp.Regexp, len(keywordPatterns))

func init() {
	for i, kp := range keywordPatterns {
		keywordRegexps[i] = regexp.MustCompile(kp.regexp)
	}
}
func ParseKeywords(u string) []string {
	keywords := make([]string, 0)
	for i, r := range keywordRegexps {
		if !r.MatchString(u) {
			continue
		}
		keywords = append(keywords, keywordPatterns[i].keyword)
	}
	return keywords
}
