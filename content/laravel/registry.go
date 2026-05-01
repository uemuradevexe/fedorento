package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var Laravel13 = data.Chapter{
	ID:    "laravel13",
	Title: "Laravel 13",
	Sections: []data.Section{
		RoutesSection,
		MiddlewareSection,
		ControllersSection,
		RequestsSection,
		ValidationSection,
		MigrationsSection,
		ModelsSection,
		EloquentSection,
	},
}
