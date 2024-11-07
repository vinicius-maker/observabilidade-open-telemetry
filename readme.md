# Desafio: Observabilidade & Open Telemetry

Tracing distribuído e span - Desenvolver um sistema distribuído em Go com dois serviços (Serviço A e Serviço B) que, ao receber um CEP, identifica a cidade e retorna o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin) junto com o nome da cidade. O sistema deverá implementar tracing distribuído utilizando OTEL (OpenTelemetry) e Zipkin para monitoramento e rastreamento.

## Arquitetura do Sistema

- **Serviço A (Input)**: Responsável por receber o input do CEP e encaminhar para o Serviço B.
- **Serviço B (Orquestração)**: Responsável por consultar a localização e a temperatura com base no CEP fornecido pelo Serviço A e retornar a resposta formatada.

## Requisitos

### Serviço A - Responsável pelo Input

1. **Endpoint**:
    - O serviço deve expor um endpoint via `POST` com o schema: `{ "cep": "29902555" }`.

2. **Validação**:
    - Verifica se o input contém exatamente 8 dígitos e está no formato de uma string.
    - Caso o CEP seja válido, ele é encaminhado para o Serviço B via HTTP.
    - Em caso de CEP inválido, retornar:
        - Código HTTP: `422`
        - Mensagem: `invalid zipcode`

### Serviço B - Responsável pela Orquestração

1. **Recebimento do CEP**:
    - Receber um CEP válido de 8 dígitos.

2. **Processamento**:
    - Realizar a busca pelo CEP para identificar a cidade.
    - Consultar a temperatura na localização encontrada e retornar as temperaturas em:
        - **Celsius**
        - **Fahrenheit** (usando a fórmula: `F = C * 1.8 + 32`)
        - **Kelvin** (usando a fórmula: `K = C + 273`)

3. **Respostas**:
    - **Sucesso**:
        - Código HTTP: `200`
        - Response Body: `{ "city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }`
    - **Falha (CEP inválido)**:
        - Código HTTP: `422`
        - Mensagem: `invalid zipcode`
    - **Falha (CEP não encontrado)**:
        - Código HTTP: `404`
        - Mensagem: `can not find zipcode`

## Implementação de OTEL e Zipkin

1. **Tracing Distribuído**:
    - Implementar OpenTelemetry (OTEL) para realizar tracing distribuído entre o Serviço A e o Serviço B.
    - Utilizar spans para medir o tempo de resposta nas operações de busca de CEP e de busca de temperatura.

2. **Observabilidade**:
    - Utilizar o Zipkin para coletar e exibir traces distribuídos dos serviços.

## APIs Externas Utilizadas

- **viaCEP**: Para busca de informações sobre o CEP - [https://viacep.com.br/](https://viacep.com.br/)
- **WeatherAPI**: Para consultar a temperatura com base na localização - [https://www.weatherapi.com/](https://www.weatherapi.com/)

## Entrega:
- O código-fonte completo da implementação.
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.

## Configuração do projeto

1. **Clone o Repositório:**

   ```bash
   git clone https://github.com/vinicius-maker/observabilidade-open-telemetry.git
   cd observabilidade-open-telemetry

2. **Configurar variáveis de ambiente:**
   - no caminho observabilidade-open-telemetry/orchestration/cmd
    ```bash
       cp .env.example .env
   ```
   - será necessário adicionar a API_KEY do WeatherAPI

3. **Configurar docker:**
   - no diretório raiz: observabilidade-open-telemetry/

    ```bash
        docker-compose build
        docker-compose up -d

4. **Executar requisição de cep:**
    - Acessar o arquivo em observabilidade-open-telemetry/cep_request.http e executar a requisição

5. **Visualizar tracing (zipkin):**
    - Acessar: http://localhost:9411/
    - Filtrar: serviceName=input