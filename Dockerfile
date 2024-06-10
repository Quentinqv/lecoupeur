# Utilisez une image de base officielle de Go pour la phase de build
FROM golang:1.21.6-alpine AS builder

# Définir le répertoire de travail à l'intérieur du conteneur
WORKDIR /app

# Copier les fichiers de votre projet
COPY . .

# Télécharger les dépendances et construire l'application
RUN go mod download
RUN go build -o lecoupeur

# Utilisez une image de base plus petite pour la phase de production
FROM alpine:latest

# Définir le répertoire de travail à l'intérieur du conteneur
WORKDIR /app

# Copier le binaire construit depuis la phase de build
COPY --from=builder /app/lecoupeur /app/lecoupeur

# Exposer le port sur lequel l'application écoute
EXPOSE 8080

# Commande pour exécuter l'application
CMD ["./lecoupeur"]
