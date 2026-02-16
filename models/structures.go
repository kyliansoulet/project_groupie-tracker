package models

type PokemonListResponse struct { // struc en liste pour la réponse de l'API /pokemon
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type TypeListResponse struct { // pareil pour la réponse de l'API /type
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

type PokemonDetail struct { // structure complète d'un Pokemon (dans details)
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
		Other        struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
			} `json:"dream_world"`
		} `json:"other"`
	} `json:"sprites"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
}

type Favorite struct { // structure pour nos favoris (JSON)
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type CollectionPageData struct { // structure envoyée à la page Collection (HTML)
	Pokemons   []PokemonDetail
	Types      []string // pour remplir le select
	Page       int
	PrevPage   int
	NextPage   int
	TotalPages int
	Search     string
	TypeFilter string
}
