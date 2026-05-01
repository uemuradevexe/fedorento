package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var ControllersSection = data.Section{
	ID:    "controllers",
	Title: "Controllers",
	Topics: []data.Topic{
		{
			ID:    "controller-basic",
			Title: "Controller Básico",
			Description: "Agrupa handlers HTTP relacionados numa classe PHP. " +
				"Criado via `php artisan make:controller NomeController`.",
			Code: `<?php
// app/Http/Controllers/MovieController.php
namespace App\Http\Controllers;

use App\Models\Movie;
use Illuminate\Http\Request;
use Illuminate\View\View;
use Illuminate\Http\RedirectResponse;

class MovieController extends Controller
{
    public function index(): View
    {
        $movies = Movie::latest()->paginate(12);
        return view('movies.index', compact('movies'));
    }

    public function show(Movie $movie): View
    {
        // Route Model Binding — Laravel injeta o model automaticamente
        return view('movies.show', compact('movie'));
    }

    public function store(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'title'        => 'required|string|max:255',
            'description'  => 'nullable|string',
            'release_year' => 'required|integer|min:1888|max:2099',
        ]);

        Movie::create($validated);

        return redirect()->route('movies.index')
            ->with('success', 'Filme criado com sucesso!');
    }
}`,
			Explanation: "Type hints nos retornos (`View`, `RedirectResponse`) melhoram IDE e legibilidade. " +
				"Route Model Binding (`Movie $movie`) elimina o `Movie::findOrFail($id)` manual. " +
				"Sempre retorne `redirect()->with()` após POST para seguir o padrão PRG.",
			Language: "php",
		},
		{
			ID:    "controller-resource",
			Title: "Resource Controller",
			Description: "Implementa os 7 métodos CRUD convencionais. " +
				"Gerado com `php artisan make:controller MovieController --resource --model=Movie`.",
			Code: `<?php
namespace App\Http\Controllers;

use App\Models\Movie;
use App\Http\Requests\StoreMovieRequest;
use App\Http\Requests\UpdateMovieRequest;

class MovieController extends Controller
{
    public function index()    { /* lista todos */ }
    public function create()   { /* form de criação */ }
    public function store(StoreMovieRequest $request)
    {
        // Form Request cuida da validação e autorização
        $movie = Movie::create($request->validated());
        return redirect()->route('movies.show', $movie);
    }
    public function show(Movie $movie)  { /* detalhe */ }
    public function edit(Movie $movie)  { /* form de edição */ }
    public function update(UpdateMovieRequest $request, Movie $movie)
    {
        $movie->update($request->validated());
        return redirect()->route('movies.show', $movie);
    }
    public function destroy(Movie $movie)
    {
        $movie->delete();
        return redirect()->route('movies.index');
    }
}`,
			Explanation: "Form Requests (`make:request`) centralizam validação e autorização fora do controller. " +
				"Use `$request->validated()` — nunca `$request->all()` — para só passar campos validados ao model.",
			Language: "php",
		},
	},
}
