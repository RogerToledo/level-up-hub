#!/bin/bash

echo "=== Limpando instalação anterior ==="
asdf uninstall golang 1.26.1

echo "=== Atualizando plugin asdf-golang ==="
asdf plugin update golang

echo "=== Instalando Go 1.26.1 ==="
asdf install golang 1.26.1

echo "=== Definindo como versão global ==="
asdf global golang 1.26.1

echo "=== Limpando caches ==="
go clean -cache -modcache

echo "=== Removendo vendor ==="
rm -rf vendor/

echo "=== Testando compilação ==="
go build ./cmd/api

echo "=== Verificando instalação ==="
go version
go env GOROOT

echo "=== Concluído! ==="
