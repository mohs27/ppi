package pages

import (
	"html/template"
	"log"
	"net/http"

	"codeberg.org/imabritishcow/librarian/api"
	"codeberg.org/imabritishcow/librarian/templates"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	videoData := api.GetVideo(vars["channel"], vars["video"])
	videoStream := api.GetVideoStream(videoData.LbryUrl)
	comments := api.GetComments(videoData.ClaimId, videoData.Channel.Id, videoData.Channel.Name)

	videoTemplate, _ := template.ParseFS(templates.GetFiles(), "video.html")
	err := videoTemplate.Execute(w, map[string]interface{}{
		"stream":         videoStream,
		"video":          videoData,
		"comments":       comments,
		"commentsLength": len(comments),
		"config":         viper.AllSettings(),
	})
	if err != nil {
		log.Fatal(err)
	}
}
