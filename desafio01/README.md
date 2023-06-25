# Desafio do projeto GO com Docker

Este é um projeto simples em Go que exibe uma mensagem de teste. O projeto usa um Dockerfile para facilitar a criação e execução do aplicativo em um contêiner Docker.

## Pré-requisitos

Certifique-se de ter o Docker instalado em sua máquina antes de prosseguir.

- Docker: [Instalação do Docker](https://docs.docker.com/get-docker/)

## Executando o projeto com Docker

Siga as etapas abaixo para construir e executar o projeto usando Docker:

1. Clone este repositório para o seu ambiente local.

2. Navegue para o diretório do projeto:

   ```cmd
   cd nome-do-diretorio-do-projeto
   ```

3. Construa a imagem Docker usando o Dockerfile:

    ```docker
    docker build -t fullcycle .
    ```
  

4. Execute o contêiner Docker a partir da imagem criada:

    ```docker
    docker run -it --rm fullcycle
    ```
  
5. Você verá a seguinte saída no terminal:

    ```cmd
    Full Cycle Rocks!!
    ```

Isso indica que o projeto foi executado com sucesso no contêiner Docker.

## Contribuindo

Se você quiser contribuir para este projeto, fique à vontade para fazer um fork e enviar uma pull request. Ficarei feliz em revisar e mesclar suas contribuições.

## Disponibilidade no Docker Hub

Este projeto está disponível no Docker Hub para facilitar o acesso e uso. Você pode encontrá-lo em [nome-do-repositorio-no-dockerhub](link-do-repositorio-no-dockerhub). 

Para puxar a imagem do Docker Hub, execute o seguinte comando:

  ```docker
  docker pull stellaoliveirabertt/fullcycle:tagname
  ```

## Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).
