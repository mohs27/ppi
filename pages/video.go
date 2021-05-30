package pages

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imabritishcow/librarian/api"
	"github.com/imabritishcow/librarian/config"
	"github.com/imabritishcow/librarian/templates"
)


func VideoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)

	videoData := api.GetVideo(vars["channel"], vars["video"])

	videoTemplate, _ := template.ParseFS(templates.GetFiles(), "video.html")
	err := videoTemplate.Execute(w, map[string]interface{}{
		"videos": videoData.Videos,
		"stream": videoData.StreamUrl,
		"video": videoData.Videos[0],
		"config": config.GetConfig(),
	})
	if err != nil {
		log.Fatal(err)
	}
}