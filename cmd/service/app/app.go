package app

import (
	"encoding/json"
	"github.com/DaniilOr/newRest/pkg/offers"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

const OK = 200

type Server struct {
	offersSvc *offers.Service
	router    chi.Router
}
type Result struct {
	Resulg  string
	Comment string `json:"comment", omitempty`
}

func NewServer(offersSvc *offers.Service, router chi.Router) *Server {
	return &Server{offersSvc: offersSvc, router: router}
}

func (s *Server) Init() error {
	s.router.Get("/offers", s.handleGetOffers)
	s.router.Get("/offers/{id}", s.handleGetOfferByID)
	s.router.Post("/offers", s.handleSaveOffer)
	s.router.Delete("/offers/{id}", s.handleRemoveOfferByID)

	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleGetOffers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.offersSvc.All(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = makeResponse(items, OK, writer)
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) handleGetOfferByID(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.offersSvc.ByID(request.Context(), id)
	err = makeResponse(item, OK, writer)
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) handleSaveOffer(writer http.ResponseWriter, request *http.Request) {
	itemToSave := &offers.Offer{}
	err := json.NewDecoder(request.Body).Decode(&itemToSave)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.offersSvc.Save(request.Context(), itemToSave)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = makeResponse(item, OK, writer)
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) handleRemoveOfferByID(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	offer, err := s.offersSvc.Delete(request.Context(), id)
	if offer.ID == 0 && offer.Percent == "" && offer.Company == "" && offer.Comment == "" {
		res := Result{Resulg: "Error", Comment: "No such offer"}
		writer.Header().Set("Content-Type", "application/json")
		err := makeResponse(res, OK, writer)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = makeResponse(offer, OK, writer)
	if err != nil {
		log.Println(err)
		return
	}

}
func makeResponse(resp interface{}, status int, w http.ResponseWriter) error {
	if resp != nil {
		body, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}
	w.WriteHeader(status)
	return nil
}
