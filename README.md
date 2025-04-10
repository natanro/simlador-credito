# Simulador de Crédito

Este projeto implementa um simulador de crédito que utiliza uma arquitetura hexagonal, permitindo uma separação clara entre a lógica de negócios e as interações externas, como bancos de dados e filas de mensagens.

## Executando o Projeto com Docker Compose

Para rodar o projeto utilizando `docker-compose`, siga os passos abaixo:

1. **Certifique-se de ter o Docker e o Docker Compose instalados** em sua máquina.

2. **Clone o repositório**:

   ```bash
   git clone https://github.com/natanro/simlador-credito.git
   cd simlador-credito
   ```

3. **Crie um arquivo `.env`** na raiz do projeto para definir as variáveis de ambiente necessárias. Um exemplo de conteúdo pode ser:

   ```env
   MONGO_URI=mongodb://mongo:27017/simlador_credito
   ```

4. **Inicie os serviços** com o comando:

   ```bash
   docker-compose up --build
   ```

5. **Acesse a aplicação** em `http://localhost:8080` (ou a porta que você configurou).

6. **Para parar os serviços**, use:

   ```bash
   docker-compose down
   ```

## Endpoint da API

### Criar Simulação

- **URL**: `/simulacao`
- **Método HTTP**: `POST`
- **Descrição**: Este endpoint cria uma nova simulação de crédito.

#### Parâmetros do Corpo da Requisição

O corpo da requisição deve ser um JSON com os seguintes parâmetros:

```json
{
  "requested_amount": 10000.00,
  "installments": 12,
  "age": 30
}
```

- `requested_amount` (float): O valor solicitado para a simulação.
- `installments` (int): O número de parcelas desejadas.
- `age` (int): A idade do solicitante.

#### Exemplo de Chamada

Você pode usar `curl` para fazer uma chamada ao endpoint:

```bash
curl -X POST http://localhost:8080/motor-simulacao/simulacao \
-H "Content-Type: application/json" \
-d '{
  "requested_amount": 10000.00,
  "installments": 12,
  "age": 30
}'
```

#### Resposta

- **Status 202 Accepted**: A simulação foi aceita e está sendo processada.
- **Status 400 Bad Request**: Se os parâmetros fornecidos estiverem incorretos ou faltando.

## Padrões de Projeto

### Padrão Observer

O padrão Observer foi implementado na fila de mensagens localizada em `motor-simulacao/infra/queue.go`. Esse padrão permite que múltiplos observadores se inscrevam para receber notificações sobre eventos específicos. Quando uma simulação é publicada na fila, todos os observadores registrados são notificados.

```go:motor-simulacao/infra/queue.go
func (q *queue) runListener() {
	for message := range q.queue {
		for _, observer := range q.observers {
			if err := observer.Notify(message); err != nil {
				q.Publish(message)
			}
		}
	}
}
```

### Padrão Strategy

O padrão Strategy foi utilizado para a obtenção da taxa de juros com base na idade do usuário, conforme implementado em `motor-simulacao/interactor/rate_strategy.go`. Essa abordagem permite que diferentes estratégias de cálculo de taxa sejam facilmente intercambiáveis.

```go:motor-simulacao/interactor/rate_strategy.go
func (t *rateStrategy) GetRateByAge(age int) (float64, error) {
	params, err := t.paramRepository.FindParams()
	if err != nil {
		return 0, err
	}

	for _, param := range params {
		if age < 25 && param.Class == "25-" {
			return param.Rate, nil
		}
		// ... outras condições
	}
	return 0, errors.New("classe não encontrada para a idade fornecida")
}
```

### Contract-First

O desenvolvimento do endpoint foi realizado seguindo a abordagem Contract-First, utilizando Swagger para definir a API antes da implementação. Isso garante que a documentação da API esteja sempre atualizada e que o contrato entre o cliente e o servidor seja claro e bem definido.

```yaml
# Exemplo de definição Swagger (não incluído no código, mas deve ser parte da documentação)
paths:
  /simulacao:
    post:
      summary: "Cria uma nova simulação"
      responses:
        '202':
          description: "Simulação aceita"
```

## Arquitetura Hexagonal

A arquitetura hexagonal foi aplicada para separar a lógica de negócios da infraestrutura. As interações com o banco de dados e outras dependências externas são tratadas por meio de interfaces, permitindo que a lógica de negócios permaneça independente de detalhes de implementação.

```go:motor-simulacao/interactor/process_simulation.go
type SimulationProcessor interface {
	Notify(simulation *entity.Simulation) error
}
```

## Escalabilidade e Concurrency

O código foi projetado para lidar com um grande volume de requisições. O endpoint exposto retorna um status `202 ACCEPTED` após a simulação ser salva no banco de dados com o status `CREATED`. Em seguida, a simulação é encaminhada para uma fila, onde será processada em paralelo.

```go:motor-simulacao/interactor/process_simulation.go
func (s *simulationProcessor) Notify(simulation *entity.Simulation) error {
	// ... lógica para salvar a simulação
	simulation.Status = entity.SimulationStatusProcessed
	return s.simulationRepository.Update(simulation)
}
```

### Conclusão

Este simulador de crédito é um exemplo de como aplicar boas práticas de desenvolvimento e padrões de projeto para criar um sistema robusto e escalável. A arquitetura hexagonal, combinada com os padrões Observer e Strategy, proporciona uma base sólida para futuras extensões e manutenções.