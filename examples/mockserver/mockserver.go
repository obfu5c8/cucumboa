package mockserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func New() *MockServer {
	return &MockServer{
		pets: map[int]Pet{
			1234: Pet{
				Id:   1234,
				Name: "doggie",
				Category: PetCategory{
					Id:   1,
					Name: "Dogs",
				},
				Status:    "available",
				PhotoUrls: []string{},
			},
		},
	}
}

type MockServer struct {
	pets map[int]Pet
}

func (m *MockServer) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(defaultHeadersMiddleware)

	r.Get("/pet/{petId}", func(w http.ResponseWriter, r *http.Request) {
		petId, err := strconv.Atoi(chi.URLParamFromCtx(r.Context(), "petId"))
		if err != nil {
			sendJsonResponse(w, 500, err.Error())
			return
		}

		pet, ok := m.pets[petId]
		if !ok {
			sendJsonResponse(w, 404, nil)
			return
		}

		sendJsonResponse(w, 200, pet)
	})

	r.Delete("/pet/{petId}", func(w http.ResponseWriter, r *http.Request) {
		petId, err := strconv.Atoi(chi.URLParamFromCtx(r.Context(), "petId"))
		if err != nil {
			sendJsonResponse(w, 500, err.Error())
			return
		}

		delete(m.pets, petId)
		sendJsonResponse(w, 204, nil)

	})

	return r
}

func (m *MockServer) SetPets(pets map[int]Pet) {
	m.pets = pets
}

func sendJsonResponse(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func defaultHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)

	})
}
