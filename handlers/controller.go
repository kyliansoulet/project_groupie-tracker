package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"pokedex-project/models"
	"strconv"
	"strings"
)

// "vase de données" en mémoire
var AllPokemons []models.PokemonDetail
var AllTypes []string

// init se lance tout seul au démarrage
func init() {
	fmt.Println("Chargement des données API... (ça peut prendre un moment)")
	fetchData()
}

func fetchData() {
	// 1. recup les types (Endpoint 1)
	respTypes, err := http.Get("https://pokeapi.co/api/v2/type")
	if err == nil {
		defer respTypes.Body.Close()
		var typeList models.TypeListResponse
		json.NewDecoder(respTypes.Body).Decode(&typeList)
		for _, t := range typeList.Results {
			AllTypes = append(AllTypes, t.Name)
		}
	}

	// 2. recup la liste (Endpoint 2)
	// je limite à 50 pour pas faire exploser le temps de chargement
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon?limit=50")
	if err != nil {
		fmt.Println("Erreur API:", err)
		return
	}
	defer resp.Body.Close()

	var list models.PokemonListResponse
	json.NewDecoder(resp.Body).Decode(&list)

	// 3. recup les détails pour chq pokemon (Endpoint 3 avec URL)
	for _, item := range list.Results {
		respDet, err := http.Get(item.URL)
		if err == nil {
			var p models.PokemonDetail
			json.NewDecoder(respDet.Body).Decode(&p)
			AllPokemons = append(AllPokemons, p)
			respDet.Body.Close()
		}
	}
	fmt.Println("Données chargées :", len(AllPokemons), "pokemons.")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	// Récupération des paramètres URL
	search := r.URL.Query().Get("search")
	typeFilter := r.URL.Query().Get("type")

	// -- FT1 & FT2 : Filtrage manuel --
	var filtered []models.PokemonDetail
	for _, p := range AllPokemons {
		keep := true

		// Filtre Recherche
		if search != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
			keep = false
		}

		// Filtre Type
		if typeFilter != "" {
			hasType := false
			for _, t := range p.Types {
				if t.Type.Name == typeFilter {
					hasType = true
				}
			}
			if !hasType {
				keep = false
			}
		}

		if keep {
			filtered = append(filtered, p)
		}
	}

	// -- FT3 : Pagination --
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	perPage := 12 // Un multiple de 3 ou 4 c'est mieux pour la grille

	start := (page - 1) * perPage
	end := start + perPage
	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	paginated := filtered[start:end]

	data := models.CollectionPageData{
		Pokemons:   paginated,
		Types:      AllTypes,
		Page:       page,
		PrevPage:   page - 1,
		NextPage:   page + 1,
		Search:     search,
		TypeFilter: typeFilter,
	}

	renderTemplate(w, "collection.html", data)
}

func DetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	var selected models.PokemonDetail
	found := false

	// Recherche dans ma liste en mémoire
	for _, p := range AllPokemons {
		if p.ID == id {
			selected = p
			found = true
			break
		}
	}

	if !found {
		http.Redirect(w, r, "/collection", 302)
		return
	}

	renderTemplate(w, "details.html", selected)
}

// FT4 : Favori
func FavoritesHandler(w http.ResponseWriter, r *http.Request) {
	favs := getFavorites()
	renderTemplate(w, "favorites.html", favs)
}

func AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		img := r.FormValue("image")

		favs := getFavorites()

		// Vérif doublon
		exists := false
		for _, f := range favs {
			if f.ID == id {
				exists = true
			}
		}

		if !exists {
			favs = append(favs, models.Favorite{ID: id, Name: name, Image: img})
			saveFavorites(favs)
		}
	}
	http.Redirect(w, r, "/favoris", 303)
}

func RemoveFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id, _ := strconv.Atoi(r.FormValue("id"))
		favs := getFavorites()
		var newFavs []models.Favorite

		for _, f := range favs {
			if f.ID != id {
				newFavs = append(newFavs, f)
			}
		}
		saveFavorites(newFavs)
	}
	http.Redirect(w, r, "/favoris", 303)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html", nil)
}

// helpers pour JSON et Templates

func getFavorites() []models.Favorite {
	file, err := ioutil.ReadFile("favorites.json")
	if err != nil {
		return []models.Favorite{} // Retourne vide si pas de fichier
	}
	var favs []models.Favorite
	json.Unmarshal(file, &favs)
	return favs
}

func saveFavorites(favs []models.Favorite) {
	data, _ := json.MarshalIndent(favs, "", " ")
	ioutil.WriteFile("favorites.json", data, 0644)
}

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	t, err := template.ParseFiles("templates/"+tmplName, "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, "Erreur Template: "+err.Error(), 500)
		return
	}
	t.Execute(w, data)
}
