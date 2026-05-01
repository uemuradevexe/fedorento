package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var RoutesSection = data.Section{
	ID:    "routes",
	Title: "Rotas",
	Topics: []data.Topic{
		{
			ID:    "routes-basic",
			Title: "Rota Básica (Closure)",
			Description: "A forma mais simples: URI + closure. Definida em routes/web.php.\n" +
				"Serve para prototipagem — em produção prefira Controllers.",
			Code: `use Illuminate\Support\Facades\Route;

Route::get('/greeting', function () {
    return 'Hello World';
});

// Outros verbos HTTP
Route::post('/users', function () { /* ... */ });
Route::put('/users/{id}', function (string $id) { /* ... */ });
Route::delete('/users/{id}', function (string $id) { /* ... */ });

// Responder a múltiplos verbos
Route::match(['get', 'post'], '/form', function () { /* ... */ });
Route::any('/anything', function () { /* ... */ });`,
			Explanation: "Cada método (get, post, put, delete) corresponde a um verbo HTTP.\n" +
				"O segundo argumento pode ser uma closure ou um array [Controller::class, 'método'].\n" +
				"Use Route::any() com cuidado — prefira ser explícito com o verbo.",
			Language: "php",
		},
		{
			ID:    "routes-params",
			Title: "Parâmetros de Rota",
			Description: "Capturam segmentos dinâmicos da URI. Parâmetros obrigatórios usam {nome},\n" +
				"opcionais usam {nome?} com valor padrão no método.",
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
			Explanation: "Parâmetros são injetados na closure/controller na mesma ordem que aparecem na URI.\n" +
				"whereNumber(), whereAlpha(), whereAlphaNumeric() são atalhos para regex comuns.\n" +
				"Restrições evitam que rotas erradas sejam acionadas (ex: /user/abc quando espera número).",
			Language: "php",
		},
		{
			ID:    "routes-named",
			Title: "Rotas Nomeadas",
			Description: "Atribui um nome à rota para gerar URLs e redirects sem hardcode de strings.\n" +
				"Indispensável em projetos reais — renomear uma URI não quebra os links.",
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
			Explanation: "A convenção de nomes usa ponto como separador: 'user.profile', 'admin.posts.index'.\n" +
				"Route::resource() gera nomes automáticos (photos.index, photos.show, etc).\n" +
				"Execute php artisan route:list para ver todos os nomes registrados.",
			Language: "php",
		},
		{
			ID:    "routes-groups",
			Title: "Grupos de Rotas",
			Description: "Agrupa rotas que compartilham middleware, prefix ou namespace.\n" +
				"Evita repetição e mantém routes/web.php organizado.",
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
			Explanation: "Grupos podem ser aninhados — o Laravel mescla os atributos de cada nível.\n" +
				"middleware(['auth', 'verified']) aplica ambos na ordem listada.\n" +
				"Route::resource() dentro de grupos herda prefix e name do grupo.",
			Language: "php",
		},
		{
			ID:    "routes-resource",
			Title: "Resource Routes",
			Description: "Uma linha registra 7 rotas RESTful para operações CRUD.\n" +
				"Segue as convenções REST do Laravel — verbos HTTP corretos para cada ação.",
			Code: `use App\Http\Controllers\PhotoController;

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

// Múltiplos resources
Route::resources([
    'photos' => PhotoController::class,
    'posts'  => PostController::class,
]);`,
			Explanation: "php artisan route:list mostra todas as rotas com nomes e middlewares.\n" +
				"Para APIs use Route::apiResource() — exclui create e edit (sem forms HTML).\n" +
				"O parâmetro gerado é o nome no singular: photos → {photo}.",
			Language: "php",
		},
	},
}
