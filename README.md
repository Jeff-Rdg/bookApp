# BookApp

Api REST criada para realizar o controle de uma biblioteca, com controle de livros e seus respectivos autores.

## 🚀 Começando
Essas instruções permitirão que você obtenha uma cópia do projeto em operação na sua máquina local para fins de desenvolvimento e teste.

## Índice

- [Instalação](#instalação)
- [Configuração](#configuração)
- [Uso](#uso)
- [Endpoints](#endpoints)
- [Contribuição](#contribuição)

## Instalação

Certifique-se de ter o Go instalado em seu sistema antes de prosseguir.

**1. Clone o repositório:**
   ```
   git clone https://github.com/Jeff-Rdg/bookApp.git
  ```

**2. Navegue até o diretório da API:**
   ```
   cd BookApp
  ```
**3. Baixe e as dependências:**
   ```
   go mod download
  ```
**4. Compile e execute a API:**
   ```
   go run cmd/api/main.go
  ```
## ⚙️ Configuração
Para a criação do banco de dados mysql, Crie um arquivo .env na raiz do projeto, contendo as seguintes chaves:

```
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=
DB_PASSWORD=
DB_NAME=bookApi
WEB_SERVER_PORT=8000

```
Preencha conforme necessário.

Para executar o projeto, se faz necessário a criação de um banco de dados primeiro, com o nome contido na chave DB_NAME.

**Observação: Caso não queira realizar o uso do mysql, também é possivel utilizar o sqlite, só necessita tê-lo instalado na máquina, que é criado automaticamente a migração do banco de dados.**

## Uso

Para realizar o controle de livros, é necessário ser feito um cadastro de autores, e para isso a API disponibiliza um endpoint para carregar um CSV com os nomes, fazendo com que seja insertado primeiramente.
Após isso, é possível realizar as operações CRUD (Create, Read, Update and Delete) para o controle dos livros.

## Endpoints
A API possui endpoints para autores e livros, posteriormente será descrito cada funcionalidade.

### Autores

**1. GET/author**
 - Descrição: retorna informações paginadas sobre os autores, podendo ser filtrada pelo nome, limite, pagina e ordenação
 - Parâmetros: limit, page, sort, name
 - Exemplo:
   ```
   http://localhost:8080/author?limit=5&page=1
   ```
 - Respostas previstas:
  - Success:
```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 17 Jan 2024 19:07:17 GMT
Content-Length: 641
```

   ```json
{
    "response": {
        "limit": 1,
        "page": 1,
        "sort": "Id desc",
        "total_rows": 42,
        "total_pages": 42,
        "rows": [
            {
                "ID": 43,
                "CreatedAt": "2024-01-17T14:14:43.984-03:00",
                "UpdatedAt": "2024-01-17T14:14:43.984-03:00",
                "DeletedAt": null,
                "name": "Socrates"
            }
        ]
    }
}
   ```
- BadRequest:
```http
HTTP/1.1 400 Bad Request
Content-Type: application/problem+json
Date: Wed, 17 Jan 2024 19:13:39 GMT
Content-Length: 184
```
```json
{
    "title": "error to list authors",
    "detail": "please, refer to the errors property for additional details",
    "status": 400,
    "error": [
        "page must contain only numbers"
    ],
    "instance": "/author/"
}
```
**2. GET/author/:id**
 - Descrição: retorna informações especificas de um autor, buscando por seu id.
 - Parâmetros: id
 - Exemplo:
   ```
   http://localhost:8080/author/1
   ```
 - Respostas previstas:
- Success:
```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 17 Jan 2024 19:24:52 GMT
Content-Length: 150
```
```json
{
  "response": {
    "ID": 1,
    "CreatedAt": "2024-01-17T14:14:43.912-03:00",
    "UpdatedAt": "2024-01-17T14:14:43.912-03:00",
    "DeletedAt": null,
    "name": "J. K. Rowling"
  }
}
```
- Not Found:
```http
HTTP/1.1 404 Not Found
Content-Type: application/problem+json
Date: Wed, 17 Jan 2024 19:25:38 GMT
Content-Length: 177
```
```json
{
  "title": "error to find author by id",
  "detail": "please, refer to the errors property for additional details",
  "status": 404,
  "error": [
    "record not found"
  ],
  "instance": "/author/500"
}
```

**3. POST/author/upload_csv**
 - Descrição: carrega um arquivo .csv realizando a inserção dos autores informados.
 - Exemplo:
   ```
   http://localhost:8080/author/upload_csv
   ```
- Modelo arquivo csv:
```
| Nome   |
|--------|
| João   |
| Maria  |
| Carlos |
| Ana    |
```
 - Respostas previstas:

- Success:
```http
HTTP/1.1 201 Created
Content-Type: application/json
Date: Wed, 17 Jan 2024 19:24:52 GMT
Content-Length: 150
```
```json
{
"message":"csv uploaded successfully"
}
```
