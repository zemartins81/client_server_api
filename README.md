# Desafio Pós Go Expert 2024

## Descrição do Projeto

Este projeto foi desenvolvido como parte do desafio proposto na pós-graduação em GoLang da FullCycle. O objetivo deste desafio é aplicar conhecimentos sobre webserver HTTP, contextos, banco de dados e manipulação de arquivos com Go.

### Objetivos do Desafio

Você precisará nos entregar dois sistemas em Go:
- `client.go`
- `server.go`

Os requisitos para cumprir este desafio são:

1. **Client**:
    - O `client.go` deverá realizar uma requisição HTTP no `server.go` solicitando a cotação do dólar.
    - O `client.go` precisará receber do `server.go` apenas o valor atual do câmbio (campo "bid" do JSON).
    - Utilizando o package `context`, o `client.go` terá um timeout máximo de 300ms para receber o resultado do `server.go`.
    - O `client.go` terá que salvar a cotação atual em um arquivo `cotacao.txt` no formato: `Dólar: {valor}`.

2. **Server**:
    - O `server.go` deverá consumir a API contendo o câmbio de Dólar e Real no endereço: `https://economia.awesomeapi.com.br/json/last/USD-BRL` e em seguida deverá retornar no formato JSON o resultado para o cliente.
    - Usando o package `context`, o `server.go` deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
    - Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
    - O endpoint necessário gerado pelo `server.go` para este desafio será: `/cotacao` e a porta a ser utilizada pelo servidor HTTP será a 8080.


## Instruções de Execução

### Passos para Executar

1. Clone o repositório do projeto.
   ```bash
   git clone <link-do-repositorio>
   cd <nome-do-repositorio>
   ```

2. Execute o servidor.
   ```bash
   go run server.go
   ```

3. Em outro terminal, execute o cliente.
   ```bash
   go run client.go
   ```

4. Verifique o arquivo `cotacao.txt` para o valor atual do câmbio.

## Implementação

### Server

O `server.go` faz uma chamada à API de câmbio, registra a cotação no banco de dados SQLite e responde ao cliente com a cotação atual.

### Client

O `client.go` faz uma requisição ao `server.go`, recebe o valor do câmbio e salva em um arquivo de texto.

## Observações

- Os contextos são utilizados para controlar o tempo de execução das operações de chamada à API, gravação no banco de dados e recebimento de resposta do servidor.
- Em caso de timeout, erros são registrados nos logs.

