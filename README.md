**Content Declaration PDF Generator for Brazilian Post Office**
===========================================================

**Overview**
------------

This is a Go-based AWS Lambda function that generates a Content Declaration PDF for the Brazilian Post Office. The function is designed to be used with API Gateway and accepts requests in JSON format.

**Functionality**
----------------

The Lambda function takes a JSON payload containing the necessary information to generate a Content Declaration PDF. The payload is expected to be in the format defined by the `SolicitarDeclaracaoConteudo` struct in the `types` package.

The function uses the `gofpdf` library to generate the PDF document, which includes the following sections:

* Header with sender and receiver information
* Declaration of contents
* List of items being shipped
* Signature block

**API Gateway Integration**
---------------------------

The Lambda function is designed to be used with API Gateway. The API Gateway endpoint expects a JSON payload in the format defined by the `SolicitarDeclaracaoConteudo` struct.

**Request Payload**
-------------------

The request payload should be in the following format:
```json
{
  "httpMethod": "POST",
  "body": "{\"remetente\":{\"nome_remetente\":\"João da Silva\",\"logradouro_remetente\":\"Rua das Flores\",\"numero_remetente\":\"123\",\"complemento_remetente\":\"Apto 45\",\"bairro_remetente\":\"Jardim Primavera\",\"cep_remetente\":\"12345-678\",\"cidade_remetente\":\"São Paulo\",\"uf_remetente\":\"SP\",\"telefone_remetente\":\"11987654321\",\"cpf_cnpj_remetente\":\"123.456.789-00\"},\"objetosPostais\":[{\"peso\":500.5,\"destinatario\":{\"nome_destinatario\":\"Maria Oliveira\",\"telefone_destinatario\":\"21987654321\",\"logradouro_destinatario\":\"Avenida Central\",\"complemento_destinatario\":\"Bloco B\",\"numero_destinatario\":\"789\",\"cpf_cnpj_destinatario\":\"987.654.321-00\",\"bairro_destinatario\":\"Centro\",\"cidade_destinatario\":\"Rio de Janeiro\",\"uf_destinatario\":\"RJ\",\"cep_destinatario\":\"98765-432\",\"numero_nota_fiscal\":123456},\"declaracoesConteudo\":[{\"conteudo\":\"Livro\",\"quantidade\":2,\"valor_unitario\":29.90},{\"conteudo\":\"CD\",\"quantidade\":1,\"valor_unitario\":15.50}]},{\"peso\":1000.0,\"destinatario\":{\"nome_destinatario\":\"Carlos Pereira\",\"telefone_destinatario\":\"31987654321\",\"logradouro_destinatario\":\"Rua da Paz\",\"complemento_destinatario\":\"Casa\",\"numero_destinatario\":\"456\",\"cpf_cnpj_destinatario\":\"111.222.333-44\",\"bairro_destinatario\":\"Bela Vista\",\"cidade_destinatario\":\"Belo Horizonte\",\"uf_destinatario\":\"MG\",\"cep_destinatario\":\"54321-876\",\"numero_nota_fiscal\":789012},\"declaracoesConteudo\":[{\"conteudo\":\"Camiseta\",\"quantidade\":3,\"valor_unitario\":39.99}]}]}"
}
```
**Response**
-------------

The Lambda function returns a JSON response with a base64-encoded PDF document:
```json
{
    "statusCode": 200,
    "headers": {
        "Content-Type": "application/json"
    },
    "multiValueHeaders": null,
    "body": "{\"stringBase64\":\"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\"}"
}
```
