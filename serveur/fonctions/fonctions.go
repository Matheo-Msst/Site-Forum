package fonctions

import (
	"encoding/json"
	"fmt"
)

const CHEMIN_IMG_AJOUTEE string = "static/images/uploads/"
const EXTENSION_IMAGE string = ".png"

// Structure pour stocker les informations des blogs
type Informations struct {
	NumeroBlog   []int    `json:"Blog"`
	Titre        []string `json:"titre"`
	Image        []string `json:"image"`
	Descriptions []string `json:"descriptions"`
}

// Variable globale pour stocker les blogs
var Blogs Informations

// Ajouter un blog à la structure
func AjouterBlog(title, description, image string) {
	Blogs.Titre = append(Blogs.Titre, title)
	Blogs.Descriptions = append(Blogs.Descriptions, description)
	Blogs.Image = append(Blogs.Image, image)
}

// Ajouter une image à la structure (facultatif)
func AjouterImage(image string) {
	Blogs.Image = append(Blogs.Image, image)
}

// Fonction pour convertir la structure en JSON
func ConvertToJSON(info Informations) (string, error) {
	// On va créer une nouvelle structure pour formater le JSON de manière appropriée
	var blogs []map[string]string

	// Assurons-nous que le nombre de titres, images et descriptions est le même
	if len(info.Titre) != len(info.Image) || len(info.Titre) != len(info.Descriptions) {
		return "", fmt.Errorf("les longueurs des titres, images et descriptions ne sont pas égales")
	}

	// Construction des blogs
	for i := range info.Titre {
		blog := map[string]string{
			"Titre":       info.Titre[i],
			"Image":       info.Image[i], // Ajout de l'image
			"Description": info.Descriptions[i],
		}
		blogs = append(blogs, blog)
	}

	// Convertir les blogs au format JSON
	blogJSON, err := json.MarshalIndent(blogs, "", "  ")
	if err != nil {
		return "", err
	}

	return string(blogJSON), nil
}
