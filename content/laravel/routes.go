package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var RoutesSection = data.Section{
	ID:    "routes",
	Title: "Rotas",
	Topics: []data.Topic{
		{
			ID:    "routes-basic",
			Title: "Rotas Tradicionais",
			Description: "Rotas definidas diretamente no arquivo `routes/web.php`. " +
				"Cada rota mapeia um método HTTP + URI para uma closure ou controller.",
			Code: `// routes/web.php
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return view('welcome');
});

Route::get('/filmes', function () {
    $filmes = \App\Models\Movie::latest()->paginate(12);
    return view('filmes.index', compact('filmes'));
});

Route::post('/filmes', function (\Illuminate\Http\Request $request) {
    $validated = $request->validate([
        'title'       => 'required|string|max:255',
        'description' => 'nullable|string',
    ]);
    \App\Models\Movie::create($validated);
    return redirect()->route('filmes.index');
});`,
			Explanation: "Closures são boas para prototipagem rápida. " +
				"Para produção prefira Controllers — são testáveis e reutilizáveis. " +
				"Use `Route::name()` para nomear rotas e gerar URLs com `route('nome')`.",
			Language: "php",
		},
		{
			ID:    "routes-controller",
			Title: "Rotas com Controller",
			Description: "Aponta a rota para um método de Controller em vez de uma closure. " +
				"Mantém o arquivo de rotas limpo e a lógica testável.",
			Code: `// routes/web.php
use App\Http\Controllers\MovieController;

// Rota única
Route::get('/filmes', [MovieController::class, 'index'])->name('filmes.index');

// Resource controller (7 rotas de uma vez)
Route::resource('filmes', MovieController::class);

// Resource parcial
Route::resource('filmes', MovieController::class)
    ->only(['index', 'show', 'store', 'destroy']);

// Rotas agrupadas com middleware
Route::middleware(['auth'])->group(function () {
    Route::resource('filmes', MovieController::class)
        ->except(['index', 'show']);
});`,
			Explanation: "`Route::resource()` gera automaticamente: index, create, store, " +
				"show, edit, update, destroy. Execute `php artisan route:list` para ver todas as rotas.",
			Language: "php",
		},
	},
}
