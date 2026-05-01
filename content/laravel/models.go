package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var ModelsSection = data.Section{
	ID:    "models",
	Title: "Models",
	Topics: []data.Topic{
		{
			ID:          "model-basic",
			Title:       "Model Básico",
			Description: "Representa uma tabela do banco. Criado com `php artisan make:model Movie`.",
			Code: `<?php
namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;

class Movie extends Model
{
    use HasFactory, SoftDeletes;

    // Campos permitidos para mass assignment
    protected $fillable = [
        'title', 'description', 'release_year', 'poster_url', 'tmdb_id',
    ];

    // Casts automáticos de tipo
    protected $casts = [
        'release_year' => 'integer',
        'vote_average' => 'float',
        'is_featured'  => 'boolean',
    ];

    // Relacionamentos
    public function genres()
    {
        return $this->belongsToMany(Genre::class);
    }

    public function reviews()
    {
        return $this->hasMany(Review::class);
    }
}`,
			Explanation: "`$fillable` previne mass assignment vulnerability. " +
				"`$casts` converte automaticamente tipos ao ler do banco. " +
				"`SoftDeletes` adiciona `deleted_at` — registros são ocultados, não apagados.",
			Language: "php",
		},
		{
			ID:          "model-scopes",
			Title:       "Scopes e Acessores",
			Description: "Scopes encapsulam queries reutilizáveis. Acessores/Mutators transformam atributos.",
			Code: `<?php
class Movie extends Model
{
    // Local Scope — use: Movie::featured()->get()
    public function scopeFeatured($query)
    {
        return $query->where('is_featured', true);
    }

    // Scope com parâmetro — use: Movie::byGenre('acao')->get()
    public function scopeByGenre($query, string $genre)
    {
        return $query->whereHas('genres', fn($q) => $q->where('slug', $genre));
    }

    // Accessor (Laravel 9+) — $movie->poster_full_url
    public function getPosterFullUrlAttribute(): string
    {
        return $this->poster_url
            ? asset('storage/' . $this->poster_url)
            : asset('images/no-poster.jpg');
    }
}`,
			Explanation: "Scopes mantêm controllers limpos — lógica de query fica no model. " +
				"Laravel 9+ suporta accessor via `Attribute::make()` também (mais explícito).",
			Language: "php",
		},
	},
}
