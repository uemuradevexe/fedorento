package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var MigrationsSection = data.Section{
	ID:    "migrations",
	Title: "Migrations",
	Topics: []data.Topic{
		{
			ID:    "migration-structure",
			Title: "Estrutura de uma Migration",
			Description: "Toda migration tem dois métodos: up() adiciona, down() reverte.\n" +
				"São o controle de versão do banco de dados do seu projeto.",
			Code: `// Terminal
php artisan make:migration create_flights_table

// database/migrations/2024_01_01_000000_create_flights_table.php
<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Executa a migration.
     */
    public function up(): void
    {
        Schema::create('flights', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->string('airline');
            $table->timestamps();
        });
    }

    /**
     * Reverte a migration.
     */
    public function down(): void
    {
        Schema::dropIfExists('flights');
    }
};

// Comandos Artisan
// php artisan migrate           — executa migrations pendentes
// php artisan migrate:rollback  — reverte o último lote
// php artisan migrate:fresh     — drop all + migrate (DEV ONLY)
// php artisan migrate:status    — lista status de cada migration`,
			Explanation: "down() deve ser o inverso exato de up().\n" +
				"migrate:fresh apaga tudo e recria — nunca use em produção.\n" +
				"Nunca edite uma migration já executada em produção — crie uma nova migration.",
			Language: "php",
		},
		{
			ID:    "migration-columns",
			Title: "Tipos de Colunas",
			Description: "Blueprint oferece métodos para todos os tipos SQL mais comuns.\n" +
				"Modificadores complementam o tipo (nullable, default, unsigned, etc).",
			Code: `Schema::create('movies', function (Blueprint $table) {
    // Identificação
    $table->id();                          // BIGINT UNSIGNED AUTO_INCREMENT PK
    $table->uuid('uuid')->unique();        // UUID string

    // Strings
    $table->string('title');               // VARCHAR(255)
    $table->string('slug', 100)->unique(); // VARCHAR(100) UNIQUE
    $table->text('description')->nullable();
    $table->longText('body')->nullable();

    // Números
    $table->integer('release_year');
    $table->unsignedBigInteger('tmdb_id')->nullable();
    $table->float('vote_average')->default(0.0);
    $table->decimal('budget', total: 15, places: 2)->nullable();

    // Booleano / Enum
    $table->boolean('is_featured')->default(false);
    $table->enum('status', ['draft', 'published', 'archived'])->default('draft');

    // Datas
    $table->date('release_date')->nullable();
    $table->timestamp('published_at')->nullable();

    // Soft Delete + Timestamps automáticos
    $table->softDeletes();   // coluna deleted_at
    $table->timestamps();    // created_at + updated_at
});`,
			Explanation: "id() é equivalente a bigIncrements('id') — BIGINT UNSIGNED PK.\n" +
				"nullable() permite NULL — campos opcionais sempre devem ser nullable.\n" +
				"softDeletes() registra quando foi deletado em vez de remover o row.",
			Language: "php",
		},
		{
			ID:    "migration-foreign",
			Title: "Chaves Estrangeiras",
			Description: "foreignId() com constrained() cria FK com índice automaticamente.\n" +
				"onDelete e onUpdate controlam o comportamento em cascata.",
			Code: `// FK simples — infere tabela pelo nome da coluna
Schema::table('posts', function (Blueprint $table) {
    $table->foreignId('user_id')->constrained();
    // equivale a: REFERENCES users(id)
});

// FK com cascade
$table->foreignId('user_id')
    ->constrained()
    ->onUpdate('cascade')
    ->onDelete('cascade');

// Fluent: método mais curto
$table->foreignId('user_id')->constrained()->cascadeOnDelete();

// FK para tabela com nome diferente
$table->foreignId('author_id')
    ->constrained('users')   // especifica a tabela
    ->nullOnDelete();         // SET NULL quando user apagado

// Tabela pivot (Many-to-Many)
Schema::create('genre_movie', function (Blueprint $table) {
    // Convenção: nomes em ordem alfabética, singular
    $table->foreignId('genre_id')->constrained()->cascadeOnDelete();
    $table->foreignId('movie_id')->constrained()->cascadeOnDelete();
    $table->primary(['genre_id', 'movie_id']); // PK composta
});`,
			Explanation: "constrained() infere a tabela: user_id → users, movie_id → movies.\n" +
				"cascadeOnDelete() remove os filhos quando o pai é deletado.\n" +
				"nullOnDelete() seta NULL nos filhos — útil para relações opcionais.\n" +
				"A PK composta na pivot evita duplicatas na relação.",
			Language: "php",
		},
		{
			ID:    "migration-alter",
			Title: "Alterando Tabelas Existentes",
			Description: "Migrations para modificar colunas existentes em produção.\n" +
				"Siga o padrão expand/contract para zero downtime.",
			Code: `// Adicionar coluna nova
// php artisan make:migration add_poster_url_to_movies_table

Schema::table('movies', function (Blueprint $table) {
    $table->string('poster_url')->nullable()->after('description');
    $table->integer('runtime_minutes')->unsigned()->nullable()->after('release_year');
});

// Modificar coluna existente (requer doctrine/dbal)
Schema::table('movies', function (Blueprint $table) {
    $table->string('title', 500)->change(); // aumenta o tamanho
    $table->text('description')->nullable()->change(); // torna nullable
});

// Renomear coluna
Schema::table('movies', function (Blueprint $table) {
    $table->renameColumn('vote_count', 'votes_total');
});

// Remover coluna
Schema::table('movies', function (Blueprint $table) {
    $table->dropColumn('old_column');
    $table->dropColumn(['col1', 'col2']); // múltiplas
});

// Adicionar índice em coluna existente
Schema::table('movies', function (Blueprint $table) {
    $table->index('release_year');
    $table->unique(['title', 'release_year']); // índice composto
});`,
			Explanation: "after() posiciona a coluna logo após a especificada (MySQL only).\n" +
				"change() requer o pacote doctrine/dbal no Laravel < 11.\n" +
				"Em produção: sempre adicione colunas nullable, depois popule, depois adicione constraint.\n" +
				"Isso é o padrão expand/contract — evita downtime em deploys.",
			Language: "php",
		},
	},
}
