package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// MessageWillBePosted is invoked when a message is posted by a user before it is committed to the
// database. If you also want to act on edited posts, see MessageWillBeUpdated. Return values
// should be the modified post or nil if rejected and an explanation for the user.
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	p.API.LogError("accepted to query user", "user_id", post.UserId)

	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError("could not find channel", "channel_id", post.ChannelId)
		return nil, ""
	}

	if channel.Name != "random_spam" {
		return nil, ""
	}

	if post.Message[0] == '~' {
		// Remove first char, let pass
		post.Message = post.Message[1:]
		post.Props["randomSpamLetPass"] = true
		return post, ""
	}

	// TODO: Check team as well
	post.Message = "I have successfully edited this message"

	// Otherwise, allow the post through.
	return post, ""
}

// MessageWillBeUpdated is invoked when a message is updated by a user before it is committed to
// the database. If you also want to act on new posts, see MessageWillBePosted. Return values
// should be the modified post or nil if rejected and an explanation for the user. On rejection,
// the post will be kept in its previous state.
//
// If you don't need to modify or rejected updated posts, use MessageHasBeenUpdated instead.
//
// Note that this method will be called for posts updated by plugins, including the plugin that
// updated the post.
//
// This demo implementation rejects posts that @-mention the demo plugin user.
func (p *Plugin) MessageWillBeUpdated(c *plugin.Context, newPost, oldPost *model.Post) (*model.Post, string) {
	channel, err := p.API.GetChannel(oldPost.ChannelId)
	if err != nil {
		p.API.LogError("could not find channel", "channel_id", oldPost.ChannelId)
		return newPost, ""
	}

	if channel.Name != "random_spam" {
		return newPost, ""
	}

	if newPost.Props["randomSpamLetPass"] == true {
		return newPost, ""
	}

	// Don't allow users to edit their posts
	return oldPost, "Error 500: Internal server error. Please contact your website administrator with detais on this message."
}

// MessageHasBeenPosted is invoked after the message has been committed to the database. If you
// need to modify or reject the post, see MessageWillBePosted Note that this method will be called
// for posts created by plugins, including the plugin that created the post.
//
// This demo implementation logs a message to the demo channel whenever a message is posted,
// unless by the demo plugin user itself.
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		p.API.LogError("could not find channel", "channel_id", post.ChannelId)
		return
	}

	if channel.Name != "random_spam" {
		return
	}

	if post.Props["randomSpamLetPass"] == true {
		return
	}

	// This edit is fake, it will be reverted later by the WillUpdate method
	post.Message = "Lorem ipsum dolor sit amet"

	p.API.UpdatePost(post)
}

func main() {
	plugin.ClientMain(&Plugin{})
}
