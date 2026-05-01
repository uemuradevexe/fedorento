package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var EloquentSection = data.Section{
	ID:    "eloquent",
	Title: "Eloquent",
	Topics: []data.Topic{
		{
			ID:    "eloquent-insert-update",
			Title: "Inserindo e Atualizando",
			Description: "Eloquent oferece múltiplas formas de persistir dados.\n" +
				"Conheça as diferenças entre new + save(), create(), e updateOrCreate().",
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
			Explanation: "create() exige que os campos estejam em $fillable no model.\n" +
				"firstOrCreate() evita duplicatas — ideal para seeds e imports.\n" +
				"updateOrCreate() é o upsert do Eloquent — busca pela condição, atualiza ou cria.\n" +
				"update() em massa não dispara eventos do model (creating, updating, etc).",
			Language: "php",
		},
		{
			ID:    "eloquent-deleting",
			Title: "Deletando Registros",
			Description: "Delete, soft delete e restore. Soft delete é preferível em dados sensíveis\n" +
				"— preserva o histórico e permite recuperação.",
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
			Explanation: "SoftDeletes adiciona o trait no model E softDeletes() na migration.\n" +
				"withTrashed() é útil em painéis admin para ver histórico completo.\n" +
				"forceDelete() contorna o soft delete — use com cautela em produção.",
			Language: "php",
		},
		{
			ID:    "eloquent-relationships",
			Title: "Relacionamentos",
			Description: "Eloquent suporta todos os tipos de relacionamento SQL.\n" +
				"Defina métodos no model para navegar entre registros relacionados.",
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
			Explanation: "Sempre type-hint o retorno dos métodos de relacionamento — ajuda a IDE.\n" +
				"Use with() para eager load e evitar N+1 queries.\n" +
				"sync() é preferível a attach() em updates — evita duplicatas automaticamente.",
			Language: "php",
		},
		{
			ID:    "eloquent-eager",
			Title: "Eager Loading",
			Description: "Carrega relacionamentos junto com a query principal.\n" +
				"Elimina o problema N+1 — essencial para performance.",
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
			Explanation: "N+1 problem: 1 query para os users + 1 query POR user para os posts.\n" +
				"with() resolve isso com 2 queries: busca users, depois posts WHERE user_id IN (1,2,3...).\n" +
				"load() carrega relacionamentos de uma collection já existente — evita refetch.",
			Language: "php",
		},
		{
			ID:    "eloquent-collections",
			Title: "Collections",
			Description: "Eloquent retorna Collections — arrays com superpoderes.\n" +
				"Operações funcionais encadeáveis: map, filter, groupBy, flatMap e mais.",
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
			Explanation: "filter() mantém as keys originais — use values() para reiniciar índices.\n" +
				"flatMap() é útil para achatar coleções de coleções (ex: genres de múltiplos movies).\n" +
				"Para grandes datasets, prefira orderBy/groupBy na query SQL em vez de na collection.",
			Language: "php",
		},
	},
}
