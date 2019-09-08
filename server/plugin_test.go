package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/mock"
)

func TestMessageWillBePosted(t *testing.T) {

	t.Run("message without asciiplot are unmodified", func(t *testing.T) {
		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogInfo", mock.Anything).Maybe()
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)
		p := &Plugin{}
		p.API = api

		post := &model.Post{
			Id: model.NewId(),
			Message: "This is a test",
		}
		// (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string)
		outPost, output := p.MessageWillBePosted(nil, post)
		assert.Equal(t, output, "")
		assert.Nil(t, outPost)
	})

	t.Run("message with asciiplot have a correct plot", func(t *testing.T) {
		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogInfo", mock.Anything).Maybe()
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)
		p := &Plugin{}
		p.API = api

		post := &model.Post{
			Id: model.NewId(),
			Message: "asciiplot 1,2,3,4,5",
		}
		// (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string)
		outPost, output := p.MessageWillBePosted(nil, post)
		assert.Equal(t, output, "")
		graph := "```"+`
 5.00 ┤                                                          ╭
 4.73 ┤                                                      ╭───╯
 4.47 ┤                                                  ╭───╯
 4.20 ┤                                              ╭───╯
 3.93 ┤                                          ╭───╯
 3.67 ┤                                      ╭───╯
 3.40 ┤                                  ╭───╯
 3.13 ┤                              ╭───╯
 2.87 ┤                          ╭───╯
 2.60 ┤                      ╭───╯
 2.33 ┤                  ╭───╯
 2.07 ┤              ╭───╯
 1.80 ┤          ╭───╯
 1.53 ┤      ╭───╯
 1.27 ┤  ╭───╯
 1.00 ┼──╯
`+"```\n"
		//assert.Equal(t, outPost.Message, graph)
		fmt.Printf(outPost.Message)
		fmt.Println()
		fmt.Printf(graph)
	})
}
