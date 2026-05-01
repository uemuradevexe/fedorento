package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var EloquentSection = data.Section{
	ID:    "eloquent",
	Title: "Eloquent",
	Topics: []data.Topic{
		{
			ID:          "eloquent-queries",
			Title:       "Queries Avançadas",
			Description: "Eloquent ORM: queries fluentes com eager loading, agregações e subqueries.",
			Code: `<?php
// Eager loading — evita N+1 query
$movies = Movie::with(['genres', 'reviews'])->latest()->get();

// Eager loading condicional
$movies = Movie::with([
    'reviews' => fn($q) => $q->latest()->limit(3),
])->get();

// Agregações
$movies = Movie::withAvg('reviews', 'rating')
    ->withCount('reviews')
    ->having('reviews_count', '>=', 5)
    ->orderByDesc('reviews_avg_rating')
    ->take(10)
    ->get();

// each() → $movie->reviews_avg_rating, $movie->reviews_count

// whereHas — filtra por relacionamento
$movies = Movie::whereHas('genres', function ($q) {
    $q->where('slug', 'acao');
})->get();

// Raw expressions
$movies = Movie::selectRaw('*, (vote_average * 0.3 + COALESCE(reviews_avg_rating,0) * 0.7) as score')
    ->orderByDesc('score')
    ->get();`,
			Explanation: "Sempre use `with()` quando vai acessar relacionamentos em loop — " +
				"sem isso cada iteração dispara uma query extra (N+1). " +
				"`withAvg/withCount` adicionam colunas virtuais ao resultado sem subquery manual.",
			Language: "php",
		},
		{
			ID:          "eloquent-collections",
			Title:       "Collections",
			Description: "Laravel Collections: pipelines funcionais sobre arrays de models.",
			Code: `<?php
$movies = Movie::with('genres')->get(); // Retorna Collection

// map — transforma cada item
$titles = $movies->map(fn($m) => $m->title);

// filter — filtra por condição
$featured = $movies->filter(fn($m) => $m->is_featured);

// flatMap — expande relacionamentos
$allGenres = $movies->flatMap(fn($m) => $m->genres->pluck('name'));

// countBy — agrupa e conta
$byGenre = $movies->flatMap(fn($m) => $m->genres->pluck('slug'))
    ->countBy()
    ->sortDesc();

// groupBy
$byYear = $movies->groupBy('release_year');

// first / firstWhere
$movie = $movies->firstWhere('title', 'Inception');

// Encadeamento
$top = Movie::with('reviews')
    ->get()
    ->sortByDesc(fn($m) => $m->reviews->avg('rating'))
    ->take(5)
    ->values(); // reindexa o array`,
			Explanation: "Collections são lazy por padrão — encadeie métodos sem custo extra. " +
				"`values()` reindexa após `filter`/`sortBy` para evitar keys não sequenciais. " +
				"Para datasets grandes prefira fazer a ordenação no banco (orderBy) em vez de na collection.",
			Language: "php",
		},
	},
}
