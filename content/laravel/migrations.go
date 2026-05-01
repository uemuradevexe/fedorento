package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var MigrationsSection = data.Section{
	ID:    "migrations",
	Title: "Migrations",
	Topics: []data.Topic{
		{
			ID:    "migration-structure",
			Title: "Estrutura de uma Migration",
			Description: "Migration é o histórico versionado da estrutura do banco.\n" +
				"No Laravel 13, ela continua sendo a forma principal de compartilhar schema entre ambientes e manter deploys reproduzíveis.",
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
			Explanation: "A documentação reforça que up aplica a mudança e down descreve a reversão correspondente.\n" +
				"O nome do arquivo com timestamp define a ordem de execução e permite replay consistente do schema.\n" +
				"Depois que algo foi para produção, o caminho seguro quase sempre é criar nova migration em vez de reescrever a antiga.",
			Language: "php",
		},
		{
			ID:    "migration-columns",
			Title: "Tipos de Colunas",
			Description: "O Schema Blueprint expõe métodos semânticos para declarar colunas de forma agnóstica ao banco.\n" +
				"A ideia é escrever a intenção do domínio e deixar o Laravel traduzir isso para o dialeto suportado pelo driver.",
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
			Explanation: "Além do tipo em si, modifiers como nullable, default e unique definem o comportamento real da coluna.\n" +
				"A documentação também destaca tipos modernos como json, uuid, ulid e métodos específicos para relacionamento.\n" +
				"Escolher bem o tipo desde o início reduz casts desnecessários e migrações corretivas depois.",
			Language: "php",
		},
		{
			ID:    "migration-foreign",
			Title: "Chaves Estrangeiras",
			Description: "Definem integridade referencial direto no schema, não apenas por convenção no código.\n" +
				"No Laravel 13, foreignId e constrained continuam sendo o caminho mais expressivo para criar relações típicas com pouco boilerplate.",
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
			Explanation: "A documentação trata constraint como regra de verdade do banco, não como detalhe opcional.\n" +
				"cascadeOnDelete, restrictOnDelete e nullOnDelete devem refletir o comportamento real do negócio.\n" +
				"Em tabelas pivot, índices e chave composta ajudam a impedir relações duplicadas e manter consultas previsíveis.",
			Language: "php",
		},
		{
			ID:    "migration-alter",
			Title: "Alterando Tabelas Existentes",
			Description: "Nem toda migration cria tabela; muitas existem para evoluir estruturas já em produção.\n" +
				"Esse tipo de mudança exige mais cuidado porque pode afetar dados, locks, deploys paralelos e compatibilidade entre versões do código.",
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
			Explanation: "A documentação mostra rename, change, drop e criação de índices como operações comuns de manutenção.\n" +
				"Em produção, prefira mudanças graduais: expandir primeiro, migrar dados, só depois apertar constraints e remover legado.\n" +
				"Esse fluxo reduz risco de downtime e evita quebrar versões antigas do app durante o deploy.",
			Language: "php",
		},
	},
}
