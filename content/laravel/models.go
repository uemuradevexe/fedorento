package laravel

import "github.com/uemuradevexe/fedorento/internal/data"

var ModelsSection = data.Section{
	ID:    "models",
	Title: "Models",
	Topics: []data.Topic{
		{
			ID:    "model-conventions",
			Title: "Convenções e Configuração",
			Description: "Eloquent infere tabela, chave primária e timestamps automaticamente.\n" +
				"Conheça as convenções para saber quando e como sobrescrevê-las.",
			Code: `// Terminal
php artisan make:model Flight              // só o model
php artisan make:model Flight --migration  // model + migration
php artisan make:model Flight -mfsc        // model + migration + factory + seeder + controller

// app/Models/Flight.php
<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Flight extends Model
{
    // Convenções automáticas (não precisa declarar se seguir):
    // tabela     → 'flights'   (plural snake_case do nome da classe)
    // chave PK   → 'id'
    // timestamps → created_at, updated_at

    // Sobrescrever quando necessário:
    protected $table = 'my_flights';       // nome de tabela diferente
    protected $primaryKey = 'flight_id';   // PK diferente
    public $timestamps = false;            // sem timestamps
    protected $dateFormat = 'U';           // Unix timestamp
    protected $connection = 'sqlite';      // conexão específica
}`,
			Explanation: "Seguir as convenções economiza código de configuração.\n" +
				"Se herdar um banco legado com nomes diferentes, use $table e $primaryKey.\n" +
				"$dateFormat = 'U' armazena datas como Unix timestamp (inteiro) em vez de datetime.",
			Language: "php",
		},
		{
			ID:    "model-fillable",
			Title: "Mass Assignment",
			Description: "Protege o model de receber campos não autorizados via create() ou fill().\n" +
				"Laravel 13 suporta tanto $fillable quanto o atributo PHP #[Fillable].",
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
			Explanation: "MassAssignmentException é lançada se tentar preencher campo não listado em $fillable.\n" +
				"$guarded = [] desabilita a proteção — só use se tiver validação rigorosa antes.\n" +
				"Nunca passe $request->all() direto para create() — use $request->validated().",
			Language: "php",
		},
		{
			ID:    "model-retrieving",
			Title: "Buscando Registros",
			Description: "Eloquent retorna Collections ou instâncias do model.\n" +
				"Conheça os métodos mais usados para buscar dados.",
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
			Explanation: "all() carrega TUDO da tabela — evite em tabelas grandes, use paginate().\n" +
				"findOrFail() e firstOrFail() são preferidos em controllers: retornam 404 automaticamente.\n" +
				"get() executa a query; métodos como where() constroem o query builder sem executar.",
			Language: "php",
		},
		{
			ID:    "model-casts",
			Title: "Casts e Mutators",
			Description: "Casts convertem atributos automaticamente ao ler/escrever no banco.\n" +
				"Accessors e Mutators (Laravel 9+) transformam valores com Attribute::make().",
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
			Explanation: "cast 'array' serializa/deserializa JSON automaticamente — ideal para configurações.\n" +
				"Accessors não alteram o banco, só o valor retornado pelo model.\n" +
				"Mutators transformam o valor ANTES de salvar no banco.",
			Language: "php",
		},
	},
}
