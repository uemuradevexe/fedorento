package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var RequestsSection = data.Section{
	ID:    "requests",
	Title: "Requests",
	Topics: []data.Topic{
		{
			ID:       "request-access",
			Title:    "Acessando a Request",
			Audience: "shared",
			Description: "A classe Illuminate\\Http\\Request representa a requisição HTTP atual de forma orientada a objeto.\n" +
				"No Laravel 13, ela é injetada automaticamente pelo container em closures, controllers e outros pontos compatíveis.",
			Code: `<?php

namespace App\Http\Controllers;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;

class UserController extends Controller
{
    // Exemplo web
    public function store(Request $request): RedirectResponse
    {
        $name = $request->input('name');

        // salvar usuário...

        return redirect('/users');
    }

    // Exemplo API
    public function show(Request $request, string $id): JsonResponse
    {
        return response()->json([
            'id' => $id,
            'requested_by' => $request->user()?->id,
        ]);
    }
}`,
			Explanation: "A request deve transportar contexto HTTP, não concentrar regra de negócio.\n" +
				"Ela é ideal para ler input, headers, arquivos, usuário autenticado e metadados da chamada.\n" +
				"Quando a assinatura do método mistura dependências e parâmetros de rota, deixe a Request primeiro e os parâmetros depois.",
			Language: "php",
		},
		{
			ID:       "request-input",
			Title:    "Lendo Input",
			Audience: "shared",
			Description: "A Request oferece helpers para ler query string, payload JSON, formulários e tipos comuns sem parsing manual.\n" +
				"Isso reduz conversão repetida e deixa a intenção do código mais clara.",
			Code: `public function index(Request $request): JsonResponse
{
    $filters = [
        'search'   => $request->string('search')->trim()->toString(),
        'page'     => $request->integer('page', 1),
        'featured' => $request->boolean('featured'),
        'tags'     => $request->array('tags'),
        'status'   => $request->enum('status', Status::class, Status::Draft),
    ];

    return response()->json($filters);
}

// Query string apenas
$sort = $request->query('sort', 'created_at');

// JSON aninhado
$author = $request->input('author.name');`,
			Explanation: "Prefira helpers tipados como integer, boolean, array e enum quando eles fizerem sentido.\n" +
				"Isso evita casts espalhados e ambiguidades comuns de input vindo do navegador ou do cliente HTTP.\n" +
				"Mesmo com esses helpers, dados externos continuam não confiáveis até passarem por validação.",
			Language: "php",
		},
		{
			ID:       "request-files",
			Title:    "Uploads e Arquivos",
			Audience: "web",
			Description: "Uploads chegam na Request como UploadedFile e já vêm com helpers para presença, validade e armazenamento.\n" +
				"Esse fluxo é o mesmo para formulários web tradicionais e endpoints multipart de API.",
			Code: `public function storeAvatar(Request $request): JsonResponse
{
    if (! $request->hasFile('avatar') || ! $request->file('avatar')->isValid()) {
        return response()->json(['message' => 'Upload inválido.'], 422);
    }

    $path = $request->file('avatar')->store('avatars', 'public');

    return response()->json([
        'path' => $path,
        'extension' => $request->file('avatar')->extension(),
    ]);
}`,
			Explanation: "Use hasFile e isValid antes de depender do upload em fluxos manuais.\n" +
				"store gera nome único automaticamente e delega o destino ao filesystem configurado.\n" +
				"Para produção, normalmente o próximo passo é validar tipo e tamanho do arquivo antes de persistir.",
			Language: "php",
		},
		{
			ID:       "request-content-negotiation",
			Title:    "Web vs API na Request",
			Audience: "api",
			Description: "A própria Request ajuda a distinguir o tipo de cliente que está chamando a aplicação.\n" +
				"Isso é útil para adaptar resposta, paginação, erros e comportamento entre interfaces web e consumidores de API.",
			Code: `public function show(Request $request)
{
    $post = Post::published()->latest()->firstOrFail();

    if ($request->expectsJson()) {
        return response()->json($post);
    }

    return view('posts.show', compact('post'));
}

if ($request->is('admin/*')) {
    // rota da área administrativa
}

if ($request->routeIs('api.*')) {
    // rota nomeada da API
}`,
			Explanation: "expectsJson, routeIs e is ajudam a evitar condicionais frágeis baseadas em URL hardcoded.\n" +
				"Na web, o fluxo padrão tende a redirect e view; na API, a resposta costuma ser JSON com status explícito.\n" +
				"Esses helpers ajudam a respeitar o contrato HTTP esperado por cada consumidor.",
			Language: "php",
		},
	},
}
