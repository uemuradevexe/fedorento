package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var ModelsSection = data.Section{
	ID:    "models",
	Title: "Models",
	Topics: []data.Topic{
		{
			ID:    "model-conventions",
			Title: "Convenções e Configuração",
			Description: "Todo model Eloquent representa uma tabela e já nasce com várias convenções configuradas automaticamente.\n" +
				"No Laravel 13, vale entender primeiro esse caminho feliz antes de começar a sobrescrever nomes, chaves e conexões.",
			Code: `// Terminal
php artisan make:model Flight              // só o model
php artisan make:model Flight --migration  // model + migration
php artisan make:model Flight -mfsc        // model + migration + factory + seeder + controller

// app/Models/Flight.php
<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Connection;
use Illuminate\Database\Eloquent\Attributes\DateFormat;
use Illuminate\Database\Eloquent\Attributes\Table;
use Illuminate\Database\Eloquent\Attributes\WithoutTimestamps;
use Illuminate\Database\Eloquent\Model;

#[Table('my_flights', key: 'flight_id')]
#[Connection('sqlite')]
#[DateFormat('U')]
#[WithoutTimestamps]
class Flight extends Model
{
    // Se seguir as convenções, nem precisa declarar nada:
    // tabela     → flights
    // chave PK   → id
    // timestamps → created_at e updated_at
}`,
			Explanation: "Tabela plural em snake_case, chave id e timestamps são o caminho feliz do framework.\n" +
				"No Laravel 13, parte dessa configuração também pode ser expressa com attributes como #[Table], #[Connection] e #[DateFormat].\n" +
				"Quanto menos exceções você cria, mais natural fica integrar factories, bindings e relacionamentos.",
			Language: "php",
		},
		{
			ID:    "model-fillable",
			Title: "Mass Assignment",
			Description: "Controla quais atributos podem ser preenchidos em massa por create, update e fill.\n" +
				"Essa proteção evita que campos sensíveis sejam alterados silenciosamente a partir de input externo.",
			Code: `<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;  // Laravel 13+
use Illuminate\Database\Eloquent\Model;

// Forma moderna (Laravel 13) — PHP attribute
#[Fillable(['name', 'email', 'password'])]
class User extends Model {}

// ---

// Forma clássica (ainda válida)
class Flight extends Model
{
    // Campos que PODEM ser preenchidos em massa
    protected $fillable = [
        'name',
        'airline',
        'destination',
    ];

    // Alternativa: bloquear campos específicos (permite todo o resto)
    // protected $guarded = ['id', 'is_admin'];

    // Sem proteção alguma (não recomendado)
    // protected $guarded = [];
}

// Uso
$flight = Flight::create([
    'name'        => 'London to Paris',
    'airline'     => 'British Airways',
    'destination' => 'CDG',
]);`,
			Explanation: "A documentação do Laravel 13 mostra duas formas principais: a clássica com $fillable e a moderna com #[Fillable].\n" +
				"Se o app estiver em modo estrito, preencher atributo não permitido pode falhar cedo e expor problema de modelagem.\n" +
				"O fluxo seguro continua sendo validar antes e repassar apenas dados já autorizados ao model.",
			Language: "php",
		},
		{
			ID:    "model-retrieving",
			Title: "Buscando Registros",
			Description: "Eloquent funciona como um query builder orientado a models.\n" +
				"A partir dele você busca um registro, uma coleção ou agregados com uma API fluente e consistente com o restante do framework.",
			Code: `use App\Models\Flight;

// Todos os registros (retorna Collection)
$flights = Flight::all();

// Com constraints — retorna Collection
$flights = Flight::where('active', 1)
    ->orderBy('name')
    ->limit(10)
    ->get();

// Buscar por PK — retorna model ou null
$flight = Flight::find(1);

// Buscar por PK — lança 404 se não encontrar
$flight = Flight::findOrFail(1);

// Primeira ocorrência por coluna
$flight = Flight::firstWhere('active', 1);

// Primeiro resultado
$flight = Flight::where('active', 1)->first();

// Primeiro ou lança 404
$flight = Flight::where('legs', '>', 100)->firstOrFail();

// Buscar múltiplos por PK (retorna Collection)
$flights = Flight::find([1, 2, 3]);

// Contar registros
$count = Flight::where('active', 1)->count();

// Verificar existência
$exists = Flight::where('destination', 'Sydney')->exists();`,
			Explanation: "find, first, firstWhere, get, count e exists cobrem a maior parte do trabalho cotidiano.\n" +
				"A docs também alerta para o custo de all() em tabelas grandes e recomenda chunk, lazy ou cursor para grandes volumes.\n" +
				"Em controllers, os métodos *OrFail* combinam bem com a resposta 404 automática do Laravel.",
			Language: "php",
		},
		{
			ID:    "model-casts",
			Title: "Casts e Mutators",
			Description: "Casts, accessors e mutators moldam a forma como o dado entra e sai do model.\n" +
				"Eles ajudam a manter conversões e regras de apresentação perto da camada de domínio, sem espalhar transformação por controllers e views.",
			Code: `<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Casts\Attribute;
use Illuminate\Database\Eloquent\Model;

class User extends Model
{
    // Casts automáticos de tipo
    protected $casts = [
        'email_verified_at' => 'datetime',
        'is_admin'          => 'boolean',
        'options'           => 'array',    // JSON ↔ array PHP
        'price'             => 'decimal:2',
    ];

    // Accessor — $user->first_name (leitura)
    protected function firstName(): Attribute
    {
        return Attribute::make(
            get: fn (string $value) => ucfirst($value),
        );
    }

    // Accessor + Mutator — $user->name = 'JOHN' → salva 'john'
    protected function name(): Attribute
    {
        return Attribute::make(
            get: fn (string $value) => ucfirst($value),
            set: fn (string $value) => strtolower($value),
        );
    }
}

// Uso
$user = User::find(1);
echo $user->first_name;  // Chama o accessor
$user->name = 'JOHN';    // Chama o mutator (salva 'john')
$user->save();`,
			Explanation: "A docs separa bem três responsabilidades: casts tipam valores, accessors modelam leitura e mutators modelam escrita.\n" +
				"Isso é útil para datas, JSON, booleans, valores monetários e normalização de texto.\n" +
				"Quando bem usados, esses recursos deixam o resto da aplicação consumir o model em formato mais previsível.",
			Language: "php",
		},
	},
}
