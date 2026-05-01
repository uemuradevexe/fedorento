package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var RoutesSection = data.Section{
	ID:    "routes",
	Title: "Rotas",
	Topics: []data.Topic{
		{
			ID:       "routes-basic",
			Title:    "Rota Básica (Closure)",
			Audience: "shared",
			Description: "A porta de entrada mais direta do roteador: um verbo HTTP, uma URI e um callback.\n" +
				"No Laravel 13, rotas web vivem em routes/web.php e já recebem o grupo web com sessão, cookies e proteção CSRF.",
			Code: `use Illuminate\Support\Facades\Route;

// routes/web.php
Route::get('/greeting', function () {
    return 'Hello World';
});

// Outros verbos HTTP
Route::post('/users', function () { /* ... */ });
Route::put('/users/{id}', function (string $id) { /* ... */ });
Route::delete('/users/{id}', function (string $id) { /* ... */ });

// Responder a múltiplos verbos
Route::match(['get', 'post'], '/form', function () { /* ... */ });
Route::any('/anything', function () { /* ... */ });

// routes/api.php
Route::get('/user', function (\Illuminate\Http\Request $request) {
    return $request->user();
})->middleware('auth:sanctum');`,
			Explanation: "Closures são ótimas para exemplos rápidos, health checks e endpoints pequenos.\n" +
				"Quando a lógica cresce, mova para controllers para manter rotas curtas e organizadas.\n" +
				"A documentação também reforça ser explícito com verbos e evitar Route::any() quando o comportamento real é restrito.",
			Language: "php",
		},
		{
			ID:       "routes-params",
			Title:    "Parâmetros de Rota",
			Audience: "shared",
			Description: "Capturam trechos dinâmicos da URI para dentro da closure ou do controller.\n" +
				"A documentação destaca parâmetros obrigatórios, opcionais e constraints para garantir que a rota aceite apenas formatos válidos.",
			Code: `// Parâmetro obrigatório
Route::get('/user/{id}', function (string $id) {
    return 'User ' . $id;
});

// Múltiplos parâmetros
Route::get('/posts/{post}/comments/{comment}', function (
    string $postId,
    string $commentId,
) {
    // ...
});

// Parâmetro opcional (note o ? e o default)
Route::get('/user/{name?}', function (?string $name = 'Guest') {
    return $name;
});

// Restrição com expressão regular
Route::get('/user/{id}', function (string $id) {
    // ...
})->where('id', '[0-9]+');

// Restrições helpers
Route::get('/user/{id}', function (string $id) {
    // ...
})->whereNumber('id');`,
			Explanation: "Os nomes dos placeholders importam para model binding, mas a injeção posicional continua valendo.\n" +
				"Use whereNumber, whereUuid, whereUlid e whereIn para documentar a intenção sem regex verbosa.\n" +
				"Se a constraint falhar, o Laravel responde com 404 em vez de cair numa action errada.",
			Language: "php",
		},
		{
			ID:       "routes-named",
			Title:    "Rotas Nomeadas",
			Audience: "shared",
			Description: "Dá uma identidade estável para a rota além da URI física.\n" +
				"Com isso, links, redirects e verificações da rota atual passam a depender do nome, não de strings hardcoded.",
			Code: `// Definir nome
Route::get('/user/profile', function () {
    // ...
})->name('profile');

// Com controller
Route::get('/user/{id}/profile', [UserController::class, 'show'])
    ->name('user.profile');

// --- Usando o nome ---

// Gerar URL
$url = route('profile');

// Com parâmetros
$url = route('user.profile', ['id' => 1]);
// resultado: /user/1/profile

// Redirect para rota nomeada
return redirect()->route('profile');

// Verificar rota atual
if ($request->routeIs('profile')) {
    // ...
}`,
			Explanation: "A documentação recomenda nomes únicos e previsíveis, normalmente com notação por pontos.\n" +
				"route(), redirect()->route() e to_route() passam a funcionar sem acoplamento à URL final.\n" +
				"Isso facilita refactors, internacionalização de paths e reorganização de grupos sem quebrar chamadas.",
			Language: "php",
		},
		{
			ID:       "routes-groups",
			Title:    "Grupos de Rotas",
			Audience: "shared",
			Description: "Permitem compartilhar atributos entre várias rotas de uma vez.\n" +
				"São a base para organizar áreas como admin, APIs internas, webhooks e módulos protegidos por middleware.",
			Code: `// Grupo com middleware
Route::middleware(['auth'])->group(function () {
    Route::get('/dashboard', function () { /* ... */ });
    Route::get('/account', function () { /* ... */ });
});

// Grupo com prefix de URI
Route::prefix('admin')->group(function () {
    Route::get('/users', function () {
        // URI: /admin/users
    });
});

// Grupo com prefix de nome
Route::name('admin.')->group(function () {
    Route::get('/users', function () {
        // Nome: admin.users
    })->name('users');
});

// Combinando tudo
Route::middleware(['auth', 'admin'])
    ->prefix('admin')
    ->name('admin.')
    ->group(function () {
        Route::resource('users', AdminUserController::class);
        Route::resource('movies', AdminMovieController::class);
    });`,
			Explanation: "No Laravel 13, middlewares e constraints são mesclados; prefixes e names são concatenados.\n" +
				"A ordem do middleware continua importando, porque ela define o pipeline da requisição.\n" +
				"Agrupar bem as rotas reduz repetição e deixa regras transversais explícitas em um único lugar.",
			Language: "php",
		},
		{
			ID:       "routes-resource",
			Title:    "Resource Routes",
			Audience: "shared",
			Description: "Transforma um controller de recurso em um conjunto completo de rotas CRUD seguindo convenções REST.\n" +
				"A documentação usa esse padrão como caminho principal para recursos tradicionais com index, show, create, store, edit, update e destroy.",
			Code: `use App\Http\Controllers\PhotoController;

// Web: CRUD completo com telas HTML
// Registra 7 rotas de uma vez
Route::resource('photos', PhotoController::class);

// As 7 rotas geradas:
// GET    /photos           → index
// GET    /photos/create    → create
// POST   /photos           → store
// GET    /photos/{photo}   → show
// GET    /photos/{photo}/edit → edit
// PUT    /photos/{photo}   → update
// DELETE /photos/{photo}   → destroy

// Somente algumas rotas
Route::resource('photos', PhotoController::class)
    ->only(['index', 'show']);

// Excluir algumas rotas
Route::resource('photos', PhotoController::class)
    ->except(['create', 'edit']);

// API: remove create e edit automaticamente
Route::apiResource('photos', PhotoController::class);

// Múltiplos resources
Route::resources([
    'photos' => PhotoController::class,
    'posts'  => PostController::class,
]);`,
			Explanation: "Além de gerar as 7 rotas, o Laravel também gera nomes e parâmetros coerentes automaticamente.\n" +
				"only, except e apiResource permitem expor apenas o recorte necessário do recurso.\n" +
				"Quando combinado com route model binding, esse padrão reduz bastante o boilerplate de CRUD.",
			Language: "php",
		},
	},
}
