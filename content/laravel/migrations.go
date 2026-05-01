package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var MigrationsSection = data.Section{
	ID:    "migrations",
	Title: "Migrations",
	Topics: []data.Topic{
		{
			ID:          "migration-basic",
			Title:       "Migration Básica",
			Description: "Controle de versão do banco. Criada com `php artisan make:migration create_movies_table`.",
			Code: `<?php
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('movies', function (Blueprint $table) {
            $table->id();
            $table->string('title');
            $table->text('description')->nullable();
            $table->integer('release_year');
            $table->string('poster_url')->nullable();
            $table->unsignedInteger('tmdb_id')->nullable()->unique();
            $table->float('vote_average')->default(0);
            $table->boolean('is_featured')->default(false);
            $table->softDeletes(); // Adiciona coluna deleted_at
            $table->timestamps();  // created_at + updated_at
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('movies');
    }
};`,
			Explanation: "`down()` deve reverter exatamente o que `up()` fez. " +
				"`nullable()` em colunas opcionais evita erros ao inserir sem o campo. " +
				"Nunca edite uma migration já executada em produção — crie uma nova.",
			Language: "php",
		},
		{
			ID:          "migration-pivot",
			Title:       "Tabela Pivot (Many-to-Many)",
			Description: "Tabela intermediária para relacionamentos N:N entre models.",
			Code: `<?php
// php artisan make:migration create_genre_movie_table
return new class extends Migration
{
    public function up(): void
    {
        Schema::create('genre_movie', function (Blueprint $table) {
            // Convenção: nomes em ordem alfabética, singular
            $table->foreignId('genre_id')->constrained()->cascadeOnDelete();
            $table->foreignId('movie_id')->constrained()->cascadeOnDelete();
            $table->primary(['genre_id', 'movie_id']);
        });
    }
};

// No Model Movie:
// public function genres() {
//     return $this->belongsToMany(Genre::class);
// }
//
// Uso:
// $movie->genres()->attach([1, 2, 3]);
// $movie->genres()->sync([1, 4]);   // remove antigos, adiciona novos`,
			Explanation: "`constrained()` infere a tabela pelo nome da coluna (genre_id → genres). " +
				"`cascadeOnDelete()` apaga os pivots quando o parent é apagado. " +
				"`sync()` é mais útil que `attach()` em updates — evita duplicatas.",
			Language: "php",
		},
	},
}
