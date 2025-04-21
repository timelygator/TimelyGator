package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

type QueryParams struct {
	BidBrowsers    []string    // e.g. ["tg-observer-web_chrome_laptop", ...]
	Classes        []ClassItem // from the server or fallback
	FilterClasses  [][]string  // e.g. [ ["Work","Programming"], ... ]
	FilterAfk      bool        // whether to filter out AFK times
	IncludeAudible bool        // whether to include audible browser events
}

// DesktopQueryParams merges QueryParams with `_DesktopQueryParamsBase` (bid_window, bid_afk).
type DesktopQueryParams struct {
	QueryParams
	BidWindow string
	BidAfk    string
}

// AndroidQueryParams merges QueryParams with `_AndroidQueryParamsBase` (bid_android).
type AndroidQueryParams struct {
	QueryParams
	BidAndroid string
}

// It builds partial query code for either Desktop or Android based on `params`.
func CanonicalEvents(clientObj *TimelyGatorClient, params interface{}) string {
	switch p := params.(type) {

	case *DesktopQueryParams:
		// If no classes are provided, fetch from server
		if len(p.Classes) == 0 {
			p.Classes = GetClasses(clientObj)
		}

		// Convert classes and filterClasses to JSON strings
		classesJSON := encodeToJSONString(p.Classes)
		classesJSON = fixDoubleBackslashes(classesJSON)

		catFilterStr := encodeToJSONString(p.FilterClasses)

		queryLines := []string{}
		// Flood from bid_window
		queryLines = append(queryLines,
			fmt.Sprintf(`events = flood(query_bucket(find_bucket("%s")));`, p.BidWindow),
		)

		// If we have an AFK bucket, gather "not_afk"
		if p.BidAfk != "" {
			notAfkSnippet := fmt.Sprintf(`
not_afk = flood(query_bucket(find_bucket("%s")));
not_afk = filter_keyvals(not_afk, "status", ["not-afk"]);`, p.BidAfk)
			queryLines = append(queryLines, notAfkSnippet)
		}

		// If we have browsers
		if len(p.BidBrowsers) > 0 {
			// Produce the snippet for all browser events
			browserSnippet := browserEvents(p)
			queryLines = append(queryLines, browserSnippet)
			// If including audible => union that with not_afk
			if p.IncludeAudible {
				queryLines = append(queryLines, `
audible_events = filter_keyvals(browser_events, "audible", [true]);
not_afk = period_union(not_afk, audible_events);
`)
			}
			// if filterAfk => filter out times when user was AFK
			if p.FilterAfk && p.BidAfk != "" {
				queryLines = append(queryLines, `events = filter_period_intersect(events, not_afk);`)
			}
		} else {
			// no browsers => do nothing
		}

		// If classes => categorize
		if len(p.Classes) > 0 {
			queryLines = append(queryLines,
				fmt.Sprintf(`events = categorize(events, %s);`, classesJSON),
			)
		}
		// If filterClasses => filter them out
		if len(p.FilterClasses) > 0 {
			queryLines = append(queryLines,
				fmt.Sprintf(`events = filter_keyvals(events, "$category", %s);`, catFilterStr),
			)
		}

		return strings.Join(queryLines, "\n")

	case *AndroidQueryParams:
		// If no classes => fetch from server
		if len(p.Classes) == 0 {
			p.Classes = GetClasses(clientObj)
		}

		classesJSON := encodeToJSONString(p.Classes)
		classesJSON = fixDoubleBackslashes(classesJSON)
		catFilterStr := encodeToJSONString(p.FilterClasses)

		queryLines := []string{
			fmt.Sprintf(`events = flood(query_bucket(find_bucket("%s")));`, p.BidAndroid),
			`events = merge_events_by_keys(events, ["app"]);`,
		}
		if len(p.Classes) > 0 {
			queryLines = append(queryLines,
				fmt.Sprintf(`events = categorize(events, %s);`, classesJSON),
			)
		}
		if len(p.FilterClasses) > 0 {
			queryLines = append(queryLines,
				fmt.Sprintf(`events = filter_keyvals(events, "$category", %s);`, catFilterStr),
			)
		}
		return strings.Join(queryLines, "\n")

	default:
		// unknown => return empty
		return ""
	}
}

func browserEvents(p *DesktopQueryParams) string {
	code := "browser_events = [];\n"
	// We gather each recognized browser from p.BidBrowsers
	// then produce a code snippet. We do a simpler approach here:
	bpairs := browsersWithBuckets(p.BidBrowsers)
	for _, pair := range bpairs {
		browserName := pair[0]
		bucketID := pair[1]
		appsJSON := encodeToJSONString(browser_appnames[browserName])
		code += fmt.Sprintf(`
events_%[1]s = flood(query_bucket("%[2]s"));
window_%[1]s = filter_keyvals(events, "app", %s);
events_%[1]s = filter_period_intersect(events_%[1]s, window_%[1]s);
events_%[1]s = split_url_events(events_%[1]s);
browser_events = concat(browser_events, events_%[1]s);
browser_events = sort_by_timestamp(browser_events);
`, browserName, bucketID, appsJSON)
	}
	return code
}

func FullDesktopQuery(clientObj *TimelyGatorClient, params *DesktopQueryParams) string {
	// Escape double quotes
	params.BidWindow = escapeDoubleQuotes(params.BidWindow)
	params.BidAfk = escapeDoubleQuotes(params.BidAfk)
	for i, b := range params.BidBrowsers {
		params.BidBrowsers[i] = escapeDoubleQuotes(b)
	}

	// Build partial query from canonicalEvents
	cEvents := CanonicalEvents(clientObj, params)

	// Then add the lines for merging by (app,title), etc.
	query := fmt.Sprintf(`
%s
title_events = sort_by_duration(merge_events_by_keys(events, ["app", "title"]));
app_events   = sort_by_duration(merge_events_by_keys(title_events, ["app"]));
cat_events   = sort_by_duration(merge_events_by_keys(events, ["$category"]));
app_events   = limit_events(app_events, %d);
title_events = limit_events(title_events, %d);
duration     = sum_durations(events);
`, cEvents, defaultLimit, defaultLimit)

	// If we have browser buckets, produce extra snippet
	if len(params.BidBrowsers) > 0 {
		query += fmt.Sprintf(`
browser_events = split_url_events(browser_events);
browser_urls = merge_events_by_keys(browser_events, ["url"]);
browser_urls = sort_by_duration(browser_urls);
browser_urls = limit_events(browser_urls, %d);
browser_domains = merge_events_by_keys(browser_events, ["$domain"]);
browser_domains = sort_by_duration(browser_domains);
browser_domains = limit_events(browser_domains, %d);
browser_duration = sum_durations(browser_events);
`, defaultLimit, defaultLimit)
	} else {
		// no browser => set them empty
		query += `
browser_events = [];
browser_urls = [];
browser_domains = [];
browser_duration = 0;
`
	}

	// Return statement
	query += `
RETURN = {
    "events": events,
    "window": {
        "app_events": app_events,
        "title_events": title_events,
        "cat_events": cat_events,
        "active_events": not_afk,
        "duration": duration
    },
    "browser": {
        "domains": browser_domains,
        "urls": browser_urls,
        "duration": browser_duration
    }
};
`
	return query
}

// -- UTILITIES --

func fixDoubleBackslashes(s string) string {
	return strings.ReplaceAll(s, `\\\\`, `\\`)
}

// encodeToJSONString => JSON-encode ignoring errors
func encodeToJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// browsersWithBuckets => returns pairs (browserName, bucketID) if a bucket is found that references that browser
func browsersWithBuckets(bidBrowsers []string) [][2]string {
	var results [][2]string
	for browserName := range browser_appnames {
		bucketID := findBucketID(browserName, bidBrowsers)
		if bucketID != "" {
			results = append(results, [2]string{browserName, bucketID})
		}
	}
	return results
}

// findBucketID => tries to locate a user-provided bucket containing the browserName
func findBucketID(browserName string, buckets []string) string {
	for _, b := range buckets {
		if strings.Contains(b, browserName) {
			return b
		}
	}
	return ""
}

// browser_appnames => dictionary that maps "chrome" => array of possible app names
var browser_appnames = map[string][]string{
	"chrome": {
		"Google Chrome", "Google-chrome", "chrome.exe", "google-chrome-stable",
		"Chromium", "Chromium-browser", "chromium.exe",
		"Google-chrome-beta", "Google-chrome-unstable", "Brave-browser",
	},
	"firefox": {
		"Firefox", "Firefox.exe", "firefox", "firefox.exe",
		"Firefox Developer Edition", "firefoxdeveloperedition",
		"Firefox-esr", "Firefox Beta", "Nightly", "org.mozilla.firefox",
	},
	"opera":   {"opera.exe", "Opera"},
	"brave":   {"brave.exe"},
	"edge":    {"msedge.exe", "Microsoft Edge"},
	"vivaldi": {"Vivaldi-stable", "Vivaldi-snapshot", "vivaldi.exe"},
}

// defaultLimit => used in the final query for limit_events
const defaultLimit = 100

// escapeDoubleQuotes => replaces `"` with `\"`
func escapeDoubleQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
