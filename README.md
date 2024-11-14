# Chat Server Golang

Um projeto de faculdade, 6º semestre de Ciência da Computação, da disciplina de Paradigmas de Porgramação.

## Instalação e execução:

1. Clonar o repositório: 'git clone https://github.com/ArticN/chat_server_go.git

2. Após abrir o projeto no editor, instalar as dependências: 'go mod tidy'

3. O projeto foi construido com uma estrutura de cmd | internal, onde o ponto inicial da aplicação se encontra no cmd/diretório-desejado, onde cada diretório possui seu 'main.go' próprio e apartado.

Para se iniciar a aplicação, deve rodar os comandos em terminais separados, sendo necessário iniciar o server primeiro.

'go run ./cmd/server/main.go'

Após isso, fica a critério do usuário a ordem de inicio entre client e bot.

'go run ./cmd/client/main.go' e 'go run ./cmd/bot/main.go'

### Funcionamento

O servidor aceita conexões de múltiplos clientes, permitindo que eles enviem mensagens uns aos outros em um chat coletivo. O projeto inclui bots que respondem automaticamente a certos comandos ou mensagens específicas, melhorando a interação.

1. Os comandos funcionam da seguinte forma:

Para enviar uma mensagem pública, \msg {mensagem} -> exemplo: \msg olá, pessoal!
Para enviar uma mensagem privada, \msg @{user} {mensagem} -> exemplo: \msg @BOT olá, bot!

Para sair, \exit

Para alterar o nickname, \changename {novoNome} -> exemplo: \changename eu
