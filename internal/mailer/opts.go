package mailer

import "time"

type PublishOption func(options *PublishOptions)

func ToTopics(topics ...string) PublishOption {
	return func(options *PublishOptions) {
		options.SendOptions.Topics = topics
	}
}

func ToUsers(users ...string) PublishOption {
	return func(options *PublishOptions) {
		options.SendOptions.Users = users
	}
}

func WithTemplate(template string) PublishOption {
	return func(options *PublishOptions) {
		options.SendOptions.Template = template
	}
}

func Use(opts PublishOptions) PublishOption {
	return func(options *PublishOptions) {
		*options = opts
	}
}

func At(timestamp time.Time) PublishOption {
	return func(options *PublishOptions) {
		options.At = &timestamp
	}
}
