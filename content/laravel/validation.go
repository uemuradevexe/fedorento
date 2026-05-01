package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var ValidationSection = data.Section{
	ID:    "validation",
	Title: "Validation",
	Topics: []data.Topic{
		{
			ID:       "validation-quick",
			Title:    "Validação Rápida",
			Audience: "web",
			Description: "A forma mais direta de validar input é usar validate() na própria Request.\n" +
				"No Laravel 13, esse fluxo ainda é ótimo para casos simples e já devolve o comportamento correto para web e API automaticamente.",
			Code: `public function store(Request $request): RedirectResponse
{
    $validated = $request->validate([
        'title' => ['required', 'string', 'max:255'],
        'body' => ['required', 'string'],
        'published_at' => ['nullable', 'date'],
    ]);

    Post::create($validated);

    return to_route('posts.index');
}

// Em requests XHR / API, falhas retornam JSON 422 automaticamente.
// Em requests web, falhas redirecionam de volta com errors + old input.`,
			Explanation: "Esse é o melhor ponto de partida quando a validação é curta e vive em um único endpoint.\n" +
				"A grande vantagem é o comportamento automático: redirect com sessão na web e JSON 422 em clientes API.\n" +
				"Quando a lista de regras crescer ou precisar de autorização própria, form requests tendem a escalar melhor.",
			Language: "php",
		},
		{
			ID:       "validation-form-request",
			Title:    "Form Requests",
			Audience: "api",
			Description: "Form Request encapsula validação e autorização em uma classe própria.\n" +
				"No Laravel 13, ele continua sendo a abordagem recomendada para endpoints reais com regras mais ricas.",
			Code: `// Terminal
php artisan make:request StorePostRequest

// app/Http/Requests/StorePostRequest.php
<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\Attributes\FailOnUnknownFields;
use Illuminate\Foundation\Http\FormRequest;

#[FailOnUnknownFields]
class StorePostRequest extends FormRequest
{
    public function authorize(): bool
    {
        return $this->user()?->can('create', Post::class) ?? false;
    }

    public function rules(): array
    {
        return [
            'title' => ['required', 'string', 'max:255'],
            'body' => ['required', 'string'],
            'tags' => ['array'],
            'tags.*' => ['string', 'distinct'],
        ];
    }
}

// Controller
public function store(StorePostRequest $request): RedirectResponse
{
    Post::create($request->validated());

    return to_route('posts.index');
}`,
			Explanation: "Form Request limpa o controller e concentra regras, mensagens, autorização e preparação de input.\n" +
				"O atributo FailOnUnknownFields é um detalhe moderno útil para rejeitar chaves inesperadas cedo.\n" +
				"Esse padrão combina especialmente bem com CRUDs, APIs públicas e fluxos com policies.",
			Language: "php",
		},
		{
			ID:       "validation-prepare",
			Title:    "Preparando e Sanitizando Input",
			Audience: "shared",
			Description: "Antes de validar, você pode normalizar valores para que a regra trabalhe sobre um formato previsível.\n" +
				"Isso evita espalhar trim, slugify e defaults por controllers e services.",
			Code: `<?php

namespace App\Http\Requests;

use Illuminate\Support\Str;

class StorePostRequest extends FormRequest
{
    protected function prepareForValidation(): void
    {
        $this->merge([
            'slug' => Str::slug($this->input('slug', $this->input('title'))),
        ]);
    }

    protected function passedValidation(): void
    {
        $this->replace([
            ...$this->validated(),
            'title' => trim($this->validated('title')),
        ]);
    }
}`,
			Explanation: "prepareForValidation roda antes das regras e é ideal para normalização de entrada.\n" +
				"passedValidation é útil quando você quer reescrever o payload já aceito para o restante do fluxo.\n" +
				"Esse cuidado reduz inconsistência entre dados enviados pela web, API e integrações externas.",
			Language: "php",
		},
		{
			ID:       "validation-web-api",
			Title:    "Web vs API na Validação",
			Audience: "shared",
			Description: "A mesma regra de validação gera respostas diferentes conforme o tipo de cliente que fez a chamada.\n" +
				"Entender isso é essencial para desenhar formulários web e APIs sem duplicar tratamento de erro.",
			Code: `// Web: redirect + sessão + old input
public function storeWeb(StorePostRequest $request): RedirectResponse
{
    Post::create($request->validated());

    return to_route('posts.index');
}

// API: JSON + status code explícito
public function storeApi(StorePostRequest $request): JsonResponse
{
    $post = Post::create($request->validated());

    return response()->json($post, 201);
}

// Exemplo de resposta 422 esperada em API:
// {
//   "message": "The title field is required.",
//   "errors": {
//     "title": ["The title field is required."]
//   }
// }`,
			Explanation: "Na web, o framework preserva input e redireciona para a tela anterior.\n" +
				"Na API, o contrato padrão é JSON com status 422 e mapa de erros por campo.\n" +
				"Esse comportamento automático é uma das razões para manter a validação perto da camada HTTP.",
			Language: "php",
		},
	},
}
