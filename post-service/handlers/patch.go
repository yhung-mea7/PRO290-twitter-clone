package handlers

import (
	"fmt"
	"net/http"

	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/amqp"
	"github.com/yhung-mea7/PRO290-twitter-clone/tree/main/post-service/data"
)

func (ph *PostHandler) LikePost() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userInfo, err := ph.getUserInformation(r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"unable to connect to user service"}, rw)
			return
		}
		post := ph.repo.GetPost(uint(getPostId(r)))
		if post.Author == "" {
			ph.log.Println("[ERROR] No post found")
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&generalMesage{"No post found"}, rw)
			return
		}

		if err := ph.repo.LikePost(post.ID); err != nil {
			ph.log.Println("[ERROR] unable to like post")
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&generalMesage{"Unable to like post"}, rw)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
		ph.messenger.SubmitToMessageBroker(&amqp.Message{
			Username: post.Author,
			Message:  fmt.Sprintf("%s liked your post!", userInfo.Username),
		})

	}
}
