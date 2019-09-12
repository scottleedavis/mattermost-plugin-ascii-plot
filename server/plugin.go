package main

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/guptarohit/asciigraph"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

//Plugin strucutre
type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex

	configuration *configuration
}

//
const (
	BASE    = 10
	BITSIZE = 8
	HEIGHT  = 15
	WIDTH   = 60
	PATTERN = `(asciiplot|asciigraph)\s(-?\d+)((,\s*-?\d+)|(\s*,\s*-?\d+))*`
)

//MessageWillBePosted runs before the post is saved to store
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	regx := regexp.MustCompile(PATTERN)
	matches := regx.FindAllString(post.Message, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			pointsString := strings.TrimPrefix(m, "asciiplot ")
			pointsString = strings.TrimPrefix(pointsString, "asciigraph ")
			pointsStringArray := strings.Split(pointsString, ",")

			var numbers []float64
			if len(pointsStringArray) <= 1 {
				return nil, ""
			}
			for _, arg := range pointsStringArray {
				arg = strings.Trim(arg, " ")
				if n, err := strconv.ParseFloat(arg, 64); err == nil {
					numbers = append(numbers, n)
				}
			}
			configuration := p.getConfiguration()
			height, err := strconv.ParseInt(configuration.Height, BASE, BITSIZE)
			if err != nil {
				height = HEIGHT
			}
			width, err := strconv.ParseInt(configuration.Width, BASE, BITSIZE)
			if err != nil {
				width = WIDTH
			}
			graph := asciigraph.Plot(numbers, asciigraph.Height(int(height)), asciigraph.Width(int(width)))
			post.Message = strings.Replace(post.Message, m, "\n```\n"+graph+"\n```\n", 1)
		}

		return post, ""
	}
	return nil, ""
}
