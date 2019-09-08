package main

import (
	"strconv"
	"strings"
	"sync"

	"github.com/guptarohit/asciigraph"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex

	configuration *configuration
}

//MessageWillBePosted hook
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	if strings.HasPrefix(post.Message, "asciiplot ") {
		pointsString := strings.TrimPrefix(post.Message, "asciiplot ")
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
		graph := asciigraph.Plot(numbers, asciigraph.Height(15), asciigraph.Width(60))
		post.Message = "```\n" + graph + "\n```"
		return post, ""

	}
	return nil, ""
}
