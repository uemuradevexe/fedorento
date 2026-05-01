package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var ControllersSection = data.Section{
	ID:    "controllers",
	Title: "Controllers",
	Topics: []data.Topic{
		{
			ID:    "controller-basic",
			Title: "Controller Básico",
			Description: "Agrupa a lógica de tratamento de requisições em uma classe.\n" +
				"Criado em app/Http/Controllers. Gerado via Artisan.",
			Code: `// Terminal
php artisan make:controller UserController

// app/Http/Controllers/UserController.php
<?php

namespace App\Http\Controllers;

use App\Models\User;
use Illuminate\View\View;

class UserController extends Controller
{
    /**
     * Exibe o perfil de um usuário específico.
     */
    public function show(string $id): View
    {
        return view('user.profile', [
            'user' => User::findOrFail($id),
        ]);
    }
}

// routes/web.php
use App\Http\Controllers\UserController;

Route::get('/user/{id}', [UserController::class, 'show']);`,
			Explanation: "findOrFail() lança ModelNotFoundException (HTTP 404) se não encontrar.\n" +
				"Type hint View no retorno ajuda a IDE e documenta a intenção do método.\n" +
				"O Controller base não é obrigatório no Laravel 11+ mas é boa prática manter.",
			Language: "php",
		},
		{
			ID:    "controller-invokable",
			Title: "Single Action Controller",
			Description: "Quando um controller tem apenas uma ação, implemente __invoke().\n" +
				"Mantém o arquivo focado e a rota mais legível.",
			Code: `// Terminal
php artisan make:controller ProvisionServer --invokable

// app/Http/Controllers/ProvisionServer.php
<?php

namespace App\Http\Controllers;

class ProvisionServer extends Controller
{
    /**
     * Provisiona um novo servidor web.
     */
    public function __invoke()
    {
        // lógica aqui...
    }
}

// routes/web.php — sem especificar método
use App\Http\Controllers\ProvisionServer;

Route::post('/server', ProvisionServer::class);`,
			Explanation: "Invokable controllers seguem o princípio de responsabilidade única.\n" +
				"Ideal para actions complexas: ProcessPayment, SendWelcomeEmail, GenerateReport.\n" +
				"Na rota, passe a classe diretamente sem array — Laravel chama __invoke() automaticamente.",
			Language: "php",
		},
		{
			ID:    "controller-middleware",
			Title: "Middleware no Controller",
			Description: "Aplica middleware a métodos específicos diretamente no controller.\n" +
				"Laravel 11+ usa a interface HasMiddleware em vez do construtor.",
			Code: `<?php

namespace App\Http\Controllers;

use Illuminate\Routing\Controllers\HasMiddleware;
use Illuminate\Routing\Controllers\Middleware;

class UserController extends Controller implements HasMiddleware
{
    /**
     * Middleware aplicado ao controller.
     */
    public static function middleware(): array
    {
        return [
            'auth',                                          // todos os métodos
            new Middleware('log', only: ['index']),          // só index
            new Middleware('subscribed', except: ['store']), // todos exceto store
        ];
    }

    public function index() { /* ... */ }
    public function store() { /* ... */ }
}`,
			Explanation: "HasMiddleware substitui o $this->middleware() do construtor do Laravel 10.\n" +
				"only e except aceitam arrays: only: ['index', 'show'].\n" +
				"Middleware aplicado na rota tem precedência sobre o do controller.",
			Language: "php",
		},
		{
			ID:    "controller-resource",
			Title: "Resource Controller",
			Description: "Implementa os 7 métodos CRUD padrão. O Artisan gera o esqueleto completo.\n" +
				"Com --model, já injeta o model via Route Model Binding.",
			Code: `// Terminal
php artisan make:controller PhotoController --resource --model=Photo

// app/Http/Controllers/PhotoController.php
<?php

namespace App\Http\Controllers;

use App\Models\Photo;
use Illuminate\Http\Request;

class PhotoController extends Controller
{
    public function index()
    {
        $photos = Photo::latest()->paginate(20);
        return view('photos.index', compact('photos'));
    }

    public function create()
    {
        return view('photos.create');
    }

    public function store(Request $request)
    {
        $validated = $request->validate([
            'title' => 'required|string|max:255',
            'image' => 'required|image|max:2048',
        ]);

        Photo::create($validated);

        return redirect()->route('photos.index')
            ->with('success', 'Foto criada!');
    }

    public function show(Photo $photo)        // Route Model Binding
    {
        return view('photos.show', compact('photo'));
    }

    public function edit(Photo $photo)
    {
        return view('photos.edit', compact('photo'));
    }

    public function update(Request $request, Photo $photo)
    {
        $photo->update($request->validated());
        return redirect()->route('photos.show', $photo);
    }

    public function destroy(Photo $photo)
    {
        $photo->delete();
        return redirect()->route('photos.index');
    }
}`,
			Explanation: "Route Model Binding: Laravel busca Photo::find($id) automaticamente quando\n" +
				"o tipo-hint do parâmetro é o model. Se não encontrar, retorna 404.\n" +
				"Sempre use $request->validated() — nunca $request->all() — para evitar mass assignment.",
			Language: "php",
		},
	},
}
