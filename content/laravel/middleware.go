package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var MiddlewareSection = data.Section{
	ID:    "middleware",
	Title: "Middleware",
	Topics: []data.Topic{
		{
			ID:       "middleware-basic",
			Title:    "Middleware Básico",
			Audience: "shared",
			Description: "Middleware é a camada que inspeciona ou bloqueia a requisição antes que ela chegue na action final.\n" +
				"No Laravel 13, ele continua sendo a forma padrão de aplicar autenticação, logging, saneamento e regras transversais sem poluir controllers.",
			Code: `// Terminal
php artisan make:middleware EnsureTokenIsValid

// app/Http/Middleware/EnsureTokenIsValid.php
<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class EnsureTokenIsValid
{
    public function handle(Request $request, Closure $next): Response
    {
        if ($request->input('token') !== config('services.internal.token')) {
            return response()->json(['message' => 'Token inválido.'], 403);
        }

        return $next($request);
    }
}`,
			Explanation: "Pense no middleware como um pipeline: cada camada pode liberar, transformar ou encerrar a requisição.\n" +
				"Ele é ideal para regras que se repetem em várias rotas, como autenticação e normalização.\n" +
				"Quando a checagem é transversal, middleware costuma ser melhor escolha do que duplicar lógica em controllers.",
			Language: "php",
		},
		{
			ID:       "middleware-register",
			Title:    "Registro e Aplicação",
			Audience: "shared",
			Description: "No Laravel 13, a configuração principal de middleware mora em bootstrap/app.php.\n" +
				"Você pode registrar middleware global, de grupo ou por rota, dependendo do alcance da regra.",
			Code: `<?php

use App\Http\Middleware\EnsureTokenIsValid;
use App\Http\Middleware\EnsureUserIsSubscribed;
use Illuminate\Foundation\Configuration\Middleware;

return Application::configure(basePath: dirname(__DIR__))
    ->withMiddleware(function (Middleware $middleware): void {
        // Global: roda em toda request HTTP
        $middleware->append(EnsureTokenIsValid::class);

        // Grupo web: acrescenta ao grupo padrão
        $middleware->web(append: [
            EnsureUserIsSubscribed::class,
        ]);

        // Alias curto para uso em rota/controller
        $middleware->alias([
            'subscribed' => EnsureUserIsSubscribed::class,
        ]);
    });

// routes/web.php
Route::get('/billing', fn () => view('billing'))
    ->middleware(['auth', 'subscribed']);`,
			Explanation: "Use middleware global apenas quando a regra realmente fizer sentido para toda request.\n" +
				"Para regras de negócio, aliases e grupos costumam dar mais controle e legibilidade.\n" +
				"A documentação também destaca que web e api já possuem grupos padrão aplicados automaticamente aos seus arquivos de rota.",
			Language: "php",
		},
		{
			ID:       "middleware-params",
			Title:    "Parâmetros e Alias",
			Audience: "shared",
			Description: "Middleware pode receber parâmetros extras para variar o comportamento sem criar uma classe nova para cada caso.\n" +
				"Isso é útil para papéis, planos, features ou thresholds simples.",
			Code: `// app/Http/Middleware/EnsureUserHasRole.php
<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class EnsureUserHasRole
{
    public function handle(Request $request, Closure $next, string $role): Response
    {
        if (! $request->user() || ! $request->user()->hasRole($role)) {
            abort(403);
        }

        return $next($request);
    }
}

// bootstrap/app.php
$middleware->alias([
    'role' => EnsureUserHasRole::class,
]);

// routes/web.php
Route::get('/admin', fn () => view('admin.dashboard'))
    ->middleware('role:admin');`,
			Explanation: "Parâmetros vêm depois de : no middleware aplicado à rota.\n" +
				"Isso evita criar várias classes quase idênticas só para trocar um valor.\n" +
				"Se a regra ficar muito complexa ou depender de contexto de domínio, vale migrar para policy ou service dedicado.",
			Language: "php",
		},
		{
			ID:       "middleware-web-api",
			Title:    "Web vs API",
			Audience: "shared",
			Description: "Os grupos web e api têm responsabilidades diferentes e isso impacta sessão, CSRF e formato de resposta.\n" +
				"Saber essa diferença evita aplicar middleware errado no contexto errado.",
			Code: `// routes/web.php
Route::middleware(['web', 'auth'])->group(function () {
    Route::get('/dashboard', fn () => view('dashboard'));
    Route::post('/profile', [ProfileController::class, 'update']);
});

// routes/api.php
Route::middleware(['api', 'auth:sanctum'])->group(function () {
    Route::get('/me', fn (\Illuminate\Http\Request $request) => $request->user());
    Route::post('/posts', [Api\PostController::class, 'store']);
});`,
			Explanation: "Rotas web costumam depender de sessão, cookies, CSRF e redirect com flash de erro.\n" +
				"Rotas api tendem a ser stateless e devolver JSON, inclusive em falhas de autenticação ou validação.\n" +
				"A separação correta ajuda a manter cada fluxo com o contrato HTTP esperado pelo cliente.",
			Language: "php",
		},
	},
}
