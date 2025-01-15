package fonctions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Structure pour représenter un blog
type Blog struct {
	NumeroBlog  int    `json:"Blog"`
	Titre       string `json:"titre"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

// Fonction pour sauvegarder les données dans un fichier JSON
func Save_Data() {
	// Créer un tableau pour stocker les blogs
	var blogs []Blog

	// Parcourir les blogs existants et les structurer dans le format souhaité
	for i := 0; i < len(Blogs.Titre); i++ {
		// On récupère le numéro du blog (Numéro d'ordre)
		blog := Blog{
			NumeroBlog:  i, // Le numéro du blog commence à 0
			Titre:       Blogs.Titre[i],
			Image:       Blogs.Image[i],
			Description: Blogs.Descriptions[i],
		}
		blogs = append(blogs, blog)
	}

	// Encoder les blogs en JSON
	data, err := json.MarshalIndent(blogs, "", "    ") // Utiliser MarshalIndent pour un format lisible
	if err != nil {
		log.Fatal(err)
	}

	// Sauvegarder le JSON dans un fichier
	err = ioutil.WriteFile("blogs.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Données enregistrées dans le fichier blogs.json")
}

// Fonction pour lire les données du fichier JSON et les charger dans la structure Blogs
func Load_Data() {
	// Lire le fichier JSON
	data, err := ioutil.ReadFile("blogs.json")
	if err != nil {
		log.Fatal("Erreur lors de la lecture du fichier JSON:", err)
	}

	// Initialiser un tableau temporaire pour les blogs
	var blogs []Blog

	// Désérialiser le contenu JSON dans le tableau de blogs
	err = json.Unmarshal(data, &blogs)
	if err != nil {
		log.Fatal("Erreur lors de la désérialisation du JSON:", err)
	}

	// Maintenant, nous avons les données dans `blogs`, nous allons les copier dans `Blogs`
	for _, blog := range blogs {
		// Ajouter chaque blog à la structure `Blogs`
		Blogs.NumeroBlog = append(Blogs.NumeroBlog, blog.NumeroBlog)
		Blogs.Titre = append(Blogs.Titre, blog.Titre)
		Blogs.Image = append(Blogs.Image, blog.Image)
		Blogs.Descriptions = append(Blogs.Descriptions, blog.Description)
	}

	fmt.Println("Données chargées depuis le fichier blogs.json")
}
