package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var ControllersSection = data.Section{
	ID:    "controllers",
	Title: "Controllers",
	Topics: []data.Topic{
		{
			ID:       "controller-basic",
			Title:    "Controller Básico",
			Audience: "web",
			Description: "Organiza ações HTTP relacionadas dentro de uma classe dedicada, em vez de espalhar lógica nas rotas.\n" +
				"No Laravel 13, esse continua sendo o ponto de partida mais comum para estruturar endpoints web e API de forma limpa.",
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
			Explanation: "Use controllers quando a ação já merece nome, teste e evolução própria.\n" +
				"A documentação também reforça que métodos públicos podem receber dependências por injeção do container.\n" +
				"O controller base não é obrigatório, mas costuma ser útil para concentrar comportamento compartilhado do projeto.",
			Language: "php",
		},
		{
			ID:       "controller-invokable",
			Title:    "Single Action Controller",
			Audience: "shared",
			Description: "Use esse formato quando uma classe representar exatamente uma ação HTTP ou um caso de uso.\n" +
				"Ao implementar __invoke(), a rota aponta direto para a classe e o arquivo fica enxuto e focado.",
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
			Explanation: "Esse estilo combina bem com ações orientadas a intenção, como processar pagamento ou provisionar recurso.\n" +
				"Ele evita controllers gigantes e deixa a rota mais legível.\n" +
				"Também facilita testes isolados porque cada classe tende a ter uma única responsabilidade.",
			Language: "php",
		},
		{
			ID:       "controller-middleware",
			Title:    "Middleware no Controller",
			Audience: "shared",
			Description: "Permite declarar regras transversais do próprio controller, como autenticação, logging e assinatura ativa.\n" +
				"No Laravel 13, isso pode ser feito com HasMiddleware ou com atributos PHP, dependendo do estilo adotado pelo projeto.",
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
			Explanation: "A principal vantagem é manter as regras de acesso perto das ações que elas protegem.\n" +
				"A docs também mostra atributos como #[Middleware] quando a declaração precisa ficar ainda mais local.\n" +
				"Use only e except para evitar duplicação e deixar explícito quais métodos participam de cada regra.",
			Language: "php",
		},
		{
			ID:       "controller-resource",
			Title:    "Resource Controller",
			Audience: "shared",
			Description: "Empacota o ciclo CRUD completo de um recurso em um controller com convenções conhecidas pelo framework.\n" +
				"No Laravel 13, ele trabalha junto com Route::resource, implicit model binding e form requests para reduzir bastante repetição.",
			Code: `// Terminal
php artisan make:controller PhotoController --resource --model=Photo --requests

// app/Http/Controllers/PhotoController.php
<?php

namespace App\Http\Controllers;

use App\Models\Photo;
use App\Http\Requests\StorePhotoRequest;
use App\Http\Requests\UpdatePhotoRequest;

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

    public function store(StorePhotoRequest $request)
    {
        $photo = Photo::create($request->validated());

        return redirect()->route('photos.index')
            ->with('success', "{$photo->title} criada com sucesso.");
    }

    public function show(Photo $photo)        // Route Model Binding
    {
        return view('photos.show', compact('photo'));
    }

    public function edit(Photo $photo)
    {
        return view('photos.edit', compact('photo'));
    }

    public function update(UpdatePhotoRequest $request, Photo $photo)
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
			Explanation: "A força desse padrão está em padronizar nomes, URIs, redirects e binding do recurso inteiro.\n" +
				"Com --model e --requests, o Artisan já aproxima o código do formato recomendado para produção.\n" +
				"Quando o recurso é API-only, prefira apiResource ou controllers gerados com --api.",
			Language: "php",
		},
	},
}
