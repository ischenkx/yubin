package yubin

import "yubin/src/publication"

type PublishOption func(options *publication.Spec)

func ToTopics(topics ...string) PublishOption {
	return func(spec *publication.Spec) {
		spec.Destination.Topics = append(spec.Destination.Topics, topics...)
	}
}

func ToUsers(users ...string) PublishOption {
	return func(spec *publication.Spec) {
		spec.Destination.Users = append(spec.Destination.Users, users...)
	}
}

func WithTemplate(template string) PublishOption {
	return func(options *publication.Spec) {
		options.Template = template
	}
}

func WithSource(source string) PublishOption {
	return func(options *publication.Spec) {
		options.Source = source
	}
}

func WithProperty(key string, value any) PublishOption {
	return func(options *publication.Spec) {
		if options.Properties == nil {
			options.Properties = map[string]any{}
		}

		options.Properties[key] = value
	}
}
