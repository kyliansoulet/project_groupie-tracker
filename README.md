Projet Groupie Tracker - Version Pokémon

Salut, voici mon projet pour le cours de Go. C'est un Pokédex qui utilise l'API PokéAPI.

Comment lancer le projet
- Lancez la commande : "go run main.go"
- Ouvrez le navigateur sur "http://localhost:8080"

Comment ça fonctionne:
- Collection: Liste des pokémons
- Recherche: On peut chercher par nom
- Filtres: On peut filtrer par Type 
- Pagination: 10 par page 
- Favoris: On peut ajouter des pokémons en favoris, c'est sauvegardé dans "favorites.json".

API Utilisée
J'ai utilisé l'api "https://pokeapi.co/".

Endpoints :
- /pokemon (liste)
- /pokemon/{id} (détails)
- /type (pour les filtres)

