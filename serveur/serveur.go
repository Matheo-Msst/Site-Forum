package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"ydays/serveur/fonctions"

	"github.com/google/uuid" // Vous devez importer ce package pour générer des UUID
)

const port = ":1608"

func main() {
	fonctions.Load_Data()

	// Servir les fichiers statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/ajouter", AjoutHandler)
	http.HandleFunc("/", AccueilHandler)
	http.HandleFunc("/supprimer/", SupprimerHandler) // Ajout de la route de suppression

	// Démarrer le serveur
	fmt.Println("Serveur démarré sur http://localhost:1608")
	http.ListenAndServe(port, nil)
}

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	// Charger le template
	templateAcc := "templates/accueil.html"
	t, err := template.ParseFiles(templateAcc)
	if err != nil {
		fmt.Printf("Erreur de chargement du template: %v\n", err)
		http.Error(w, fmt.Sprintf("Erreur lors du chargement du template: %s", err), http.StatusInternalServerError)
		return
	}

	// Passer la liste de blogs au template
	err = t.Execute(w, fonctions.Blogs)
	if err != nil {
		fmt.Printf("Erreur d'exécution du template: %v\n", err)
		http.Error(w, fmt.Sprintf("Erreur lors de l'exécution du template: %s", err), http.StatusInternalServerError)
	}
}

// Fonction pour gérer l'ajout de blog
func AjoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Vérifier la taille des slices avant d'ajouter de nouvelles données
		fmt.Println("Nombre avant ajout:", "Titres : ", len(fonctions.Blogs.Titre), "Images : ", len(fonctions.Blogs.Image), "Descriptions :", len(fonctions.Blogs.Descriptions))
		// Récupérer les valeurs envoyées via le formulaire
		titre := r.FormValue("title")
		description := r.FormValue("description")
		fonctions.Blogs.Titre = append(fonctions.Blogs.Titre, titre)
		fonctions.Blogs.Descriptions = append(fonctions.Blogs.Descriptions, description)

		// Parse le formulaire pour récupérer le fichier
		err := r.ParseMultipartForm(10 << 20) // Limite de 10 Mo
		if err != nil {
			http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
			return
		}

		// Récupérer le fichier téléchargé
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Erreur lors de l'obtention du fichier", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Créer un dossier pour les uploads si il n'existe pas
		if _, err := os.Stat(fonctions.CHEMIN_IMG_AJOUTEE); os.IsNotExist(err) {
			err := os.Mkdir(fonctions.CHEMIN_IMG_AJOUTEE, os.ModePerm)
			if err != nil {
				http.Error(w, "Erreur lors de la création du dossier", http.StatusInternalServerError)
				return
			}
		}

		// Utiliser un UUID pour nommer l'image de manière unique
		imageName := uuid.New().String() + ".png" // Générer un UUID unique
		imagePath := fonctions.CHEMIN_IMG_AJOUTEE + imageName

		// Créer un fichier pour enregistrer l'image
		outFile, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Erreur lors de l'enregistrement du fichier", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copier le contenu du fichier téléchargé dans le fichier local
		_, err = outFile.ReadFrom(file)
		if err != nil {
			http.Error(w, "Erreur lors de la copie du fichier", http.StatusInternalServerError)
			return
		}

		// Ajouter le nom de l'image au tableau des images
		fonctions.Blogs.Image = append(fonctions.Blogs.Image, imagePath)
		// Vérifier la taille des slices après l'ajout de nouvelles données
		fmt.Println("Nombre après ajout:", "Titres : ", len(fonctions.Blogs.Titre), "Images : ", len(fonctions.Blogs.Image), "Descriptions :", len(fonctions.Blogs.Descriptions), "\n")
		// Sauvegarder les données
		fonctions.Save_Data()

		// Rediriger vers la page d'accueil
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// Si ce n'est pas une méthode POST, afficher le formulaire
		http.ServeFile(w, r, "templates/ajouter.html")
	}
}

// Fonction pour gérer la suppression d'un blog
func SupprimerHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'index du blog à supprimer depuis l'URL
	indexStr := strings.TrimPrefix(r.URL.Path, "/supprimer/")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(fonctions.Blogs.Titre) {
		http.Error(w, "Blog non trouvé", http.StatusNotFound)
		return
	}

	// Supprimer le blog des slices (Titres, Descriptions et Images)
	fonctions.Blogs.Titre = append(fonctions.Blogs.Titre[:index], fonctions.Blogs.Titre[index+1:]...)
	fonctions.Blogs.Descriptions = append(fonctions.Blogs.Descriptions[:index], fonctions.Blogs.Descriptions[index+1:]...)
	fonctions.Blogs.Image = append(fonctions.Blogs.Image[:index], fonctions.Blogs.Image[index+1:]...)

	// Sauvegarder les données mises à jour dans le fichier JSON
	fonctions.Save_Data()

	// Rediriger vers la page d'accueil après la suppression
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Génère des routes pour chaque page de blog
func GenererRoutesBlogs() {
	// Définir le répertoire où les pages des blogs sont stockées
	outputDir := "templates/blogs"

	// Lister tous les fichiers dans le répertoire
	files, err := os.ReadDir(outputDir)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du répertoire des blogs: %v", err)
	}

	// Pour chaque fichier dans le répertoire, définir une route
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") {
			// Extraire le numéro du blog à partir du nom du fichier
			blogID := strings.TrimSuffix(file.Name(), ".html")
			http.HandleFunc("/blogs/"+blogID, func(w http.ResponseWriter, r *http.Request) {
				// Charger le template pour ce blog
				tmplPath := filepath.Join(outputDir, file.Name())
				tmpl, err := template.ParseFiles(tmplPath)
				if err != nil {
					log.Printf("Erreur de chargement du template pour %s: %v", file.Name(), err)
					http.Error(w, fmt.Sprintf("Erreur lors du chargement du template: %s", err), http.StatusInternalServerError)
					return
				}

				// Chercher l'index du blog à partir de son ID
				index, err := strconv.Atoi(blogID)
				if err != nil || index < 1 || index > len(fonctions.Blogs.Titre) {
					http.Error(w, "Blog non trouvé", http.StatusNotFound)
					return
				}

				// Récupérer les données du blog
				blog := struct {
					NumeroBlog  string
					Titre       string
					Description string
					Image       string
				}{
					NumeroBlog:  blogID,
					Titre:       fonctions.Blogs.Titre[index-1], // index-1 car l'index est basé sur 1 dans l'URL
					Description: fonctions.Blogs.Descriptions[index-1],
					Image:       fonctions.Blogs.Image[index-1],
				}

				// Exécuter le template avec les données du blog
				err = tmpl.Execute(w, blog)
				if err != nil {
					log.Printf("Erreur lors de l'exécution du template pour %s: %v", blogID, err)
					http.Error(w, fmt.Sprintf("Erreur lors de l'exécution du template: %s", err), http.StatusInternalServerError)
				}
			})
		}
	}
}
