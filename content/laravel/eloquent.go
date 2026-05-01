package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var EloquentSection = data.Section{
	ID:    "eloquent",
	Title: "Eloquent",
	Topics: []data.Topic{
		{
			ID:    "eloquent-insert-update",
			Title: "Inserindo e Atualizando",
			Description: "Eloquent oferece diferentes estratégias para persistir dados dependendo do fluxo de criação e atualização.\n" +
				"A documentação do Laravel 13 separa bem quando usar instância manual, mass assignment, métodos *OrCreate* e updates em massa.",
			Code: `use App\Models\Flight;

// --- INSERIR ---

// Forma 1: instanciar + salvar
$flight = new Flight;
$flight->name = 'London to Paris';
$flight->save();

// Forma 2: create() com mass assignment
$flight = Flight::create([
    'name'        => 'London to Paris',
    'airline'     => 'British Airways',
    'destination' => 'CDG',
]);

// firstOrCreate — busca ou cria
$flight = Flight::firstOrCreate(
    ['number' => 'BA114'],                    // condição de busca
    ['name' => 'London to Paris']            // atributos extras se criar
);

// --- ATUALIZAR ---

// Buscar e salvar
$flight = Flight::find(1);
$flight->name = 'Paris to London';
$flight->save();

// Update em massa (afeta múltiplos registros)
Flight::where('active', 1)
    ->where('destination', 'London')
    ->update(['delayed' => 1]);

// updateOrCreate — atualiza ou cria
$flight = Flight::updateOrCreate(
    ['departure' => 'Oakland', 'destination' => 'San Diego'],
    ['price' => 99, 'discounted' => 1]
);`,
			Explanation: "save() é explícito e funciona bem quando a entidade é montada passo a passo.\n" +
				"create() e updateOrCreate() reduzem boilerplate, mas dependem de uma política clara de mass assignment.\n" +
				"A documentação também lembra que updates em massa pulam eventos do model porque os registros não são hidratados individualmente.",
			Language: "php",
		},
		{
			ID:    "eloquent-deleting",
			Title: "Deletando Registros",
			Description: "Eloquent suporta remoção definitiva, remoção lógica e restauração de registros.\n" +
				"A escolha entre hard delete e soft delete depende do valor histórico do dado e das regras de auditoria do domínio.",
			Code: `use App\Models\Flight;
use Illuminate\Database\Eloquent\SoftDeletes;

// Hard delete — remove permanentemente
$flight = Flight::find(1);
$flight->delete();

// Delete por query
Flight::where('destination', 'cancelled')->delete();

// --- SOFT DELETE ---
// No model: use SoftDeletes;
// Na migration: $table->softDeletes();

// Ao chamar delete(), Eloquent preenche deleted_at em vez de remover
$flight = Flight::find(1);
$flight->delete(); // seta deleted_at = now()

// Queries automáticas excluem soft deleted (where deleted_at IS NULL)
$flights = Flight::all(); // NÃO inclui deletados

// Incluir soft deleted na query
$flights = Flight::withTrashed()->get();
$flight  = Flight::withTrashed()->find(1);

// Somente soft deleted
$flights = Flight::onlyTrashed()->get();

// Restaurar
$flight = Flight::withTrashed()->find(1);
$flight->restore();

// Force delete — remove do banco mesmo com soft delete ativado
$flight->forceDelete();`,
			Explanation: "Soft delete não remove a linha; ele apenas preenche deleted_at e ajusta as queries padrão.\n" +
				"A documentação mostra withTrashed e onlyTrashed como ferramentas para fluxos de admin, recuperação e histórico.\n" +
				"forceDelete deve ser reservado para cenários em que apagar de verdade é um requisito explícito.",
			Language: "php",
		},
		{
			ID:    "eloquent-relationships",
			Title: "Relacionamentos",
			Description: "Relacionamentos permitem modelar como registros se conectam e também consultar essas conexões de forma fluente.\n" +
				"A documentação trata cada relação como método do model e, ao mesmo tempo, como query builder reutilizável.",
			Code: `<?php
// --- HAS MANY (Um para Muitos) ---
class Post extends Model
{
    public function comments(): HasMany
    {
        return $this->hasMany(Comment::class);
    }
}

// --- BELONGS TO (Muitos para Um) ---
class Comment extends Model
{
    public function post(): BelongsTo
    {
        return $this->belongsTo(Post::class);
    }
}

// --- MANY TO MANY ---
class Movie extends Model
{
    public function genres(): BelongsToMany
    {
        return $this->belongsToMany(Genre::class);
    }
}

// --- USO ---
// Acessar relacionamento (lazy load — 1 query extra)
$comments = Post::find(1)->comments;

// Criar registro relacionado
$post = Post::find(1);
$post->comments()->create([
    'body' => 'Ótimo post!',
]);

// Many-to-Many: attach / detach / sync
$movie->genres()->attach([1, 2, 3]);       // adiciona
$movie->genres()->detach([1]);             // remove
$movie->genres()->sync([2, 3, 4]);         // sincroniza (remove antigos)`,
			Explanation: "hasMany, belongsTo e belongsToMany cobrem boa parte dos casos do dia a dia.\n" +
				"A documentação também reforça convenções de chaves, propriedades dinâmicas e customização quando o banco foge do padrão.\n" +
				"Definir relações corretamente simplifica criação de filhos, filtros por relação e carregamento antecipado.",
			Language: "php",
		},
		{
			ID:    "eloquent-eager",
			Title: "Eager Loading",
			Description: "Eager loading carrega relações de forma planejada para evitar consultas extras em loops.\n" +
				"Na prática, ele é uma das otimizações mais importantes ao sair de exemplos simples e entrar em listagens reais.",
			Code: `use App\Models\User;

// SEM eager loading — N+1 problem (1 query + N queries)
$users = User::all();
foreach ($users as $user) {
    echo $user->posts->count(); // 1 query por user!
}

// COM eager loading — apenas 2 queries no total
$users = User::with('posts')->get();
foreach ($users as $user) {
    echo $user->posts->count(); // sem query extra!
}

// Múltiplos relacionamentos
$users = User::with(['posts', 'comments'])->get();

// Relacionamento aninhado
$users = User::with('posts.comments')->get();

// Com constraints (subquery no eager load)
$users = User::with([
    'posts' => function ($query) {
        $query->where('active', 1)->orderBy('created_at', 'desc');
    },
])->get();

// Eager load em collection já carregada
$users = User::all();
$users->load('posts');
$users->load(['posts', 'comments.author']);

// Lazy eager loading (quando já tem o model)
$user = User::find(1);
if ($someCondition) {
    $user->load('posts');
}`,
			Explanation: "A documentação do Laravel 13 detalha with para eager load inicial e load para eager load tardio.\n" +
				"Também vale combinar eager loading com constraints para não trazer filhos desnecessários.\n" +
				"Se o projeto quiser endurecer isso, o framework ainda permite impedir lazy loading em desenvolvimento.",
			Language: "php",
		},
		{
			ID:    "eloquent-collections",
			Title: "Collections",
			Description: "Consultas que retornam vários models produzem Eloquent Collections, não arrays crus.\n" +
				"Isso permite tratar o resultado com uma API funcional rica antes de renderizar, agregar ou transformar os dados.",
			Code: `use App\Models\User;

// get() retorna uma Collection
$users = User::where('active', 1)->get();

// Iterar como array normal
foreach ($users as $user) {
    echo $user->name;
}

// map — transforma cada item
$names = $users->map(fn (User $user) => $user->name);

// filter — filtra por condição (preserva keys)
$admins = $users->filter(fn (User $user) => $user->is_admin);

// reject — inverso do filter
$regular = $users->reject(fn (User $user) => $user->is_admin);

// flatMap — expande coleções aninhadas
$allGenres = Movie::with('genres')->get()
    ->flatMap(fn ($movie) => $movie->genres->pluck('name'));

// countBy — conta por valor
$byRole = $users->countBy(fn ($user) => $user->role);
// ['admin' => 3, 'editor' => 10, 'user' => 87]

// groupBy — agrupa por campo
$byYear = Movie::all()->groupBy('release_year');

// Eager load em collection
$users->load(['comments', 'posts' => fn ($q) => $q->where('active', 1)]);

// find em collection (sem nova query)
$user = $users->find(5);

// Encadeamento
$top = Movie::with('reviews')
    ->get()
    ->sortByDesc(fn ($m) => $m->reviews->avg('rating'))
    ->take(5)
    ->values(); // values() reindexa o array`,
			Explanation: "A documentação lembra que Eloquent Collection estende a Collection base e ainda adiciona comportamentos próprios de models.\n" +
				"Métodos como map, filter e groupBy são ótimos depois que o conjunto já está em memória.\n" +
				"Mas, para volumes altos, continua sendo melhor empurrar ordenação, filtro e agregação para o SQL.",
			Language: "php",
		},
	},
}
