package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	// "log"
)

//easyjson:json
type UsersSt struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	r := regexp.MustCompile("@")
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	i := -1

	var user UsersSt
	for scanner.Scan() {
		i++
		user.UnmarshalJSON([]byte(scanner.Text()))
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browsers {
			if ok := strings.Contains(browserRaw, "Android"); ok {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browserRaw {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}
		}

		for _, browserRaw := range user.Browsers {
			if ok := strings.Contains(browserRaw, "MSIE"); ok {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browserRaw {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browserRaw)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		//log.Println("Android and MSIE user:", user.name, user.email)
		email := r.ReplaceAllString(user.Email, " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
