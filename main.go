package main

import (
	"fmt"
	"net/http"
	"pokedex-project/handlers"
	"time"
)

func main() {
	// charge mes fichiers CSS/Images
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Mes routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/collection", handlers.CollectionHandler)     // le principal (FT1, FT2, FT3)
	http.HandleFunc("/pokemon", handlers.DetailHandler)            // Page détail
	http.HandleFunc("/favoris", handlers.FavoritesHandler)         // FT4
	http.HandleFunc("/add-fav", handlers.AddFavoriteHandler)       // ajouter en fav
	http.HandleFunc("/remove-fav", handlers.RemoveFavoriteHandler) // supprimer des fav
	http.HandleFunc("/about", handlers.AboutHandler)               // La FAQ

	fmt.Println("Serveur lancé sur http://localhost:8080 (Attends quelques secondes que l'API charge...)")

	// ptit timeout pour pas que ça plante direct
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server.ListenAndServe()
}
