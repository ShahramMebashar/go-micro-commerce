package api

import (
	"encoding/json"
	"errors"
	"microservice/services/product-service/internal/domain"
	"microservice/services/product-service/internal/infrastructure/validator"
	"microservice/services/product-service/internal/interfaces"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service   interfaces.Service
	validator *validator.Validator
}

func NewProductHandler(service interfaces.Service, validator *validator.Validator) *ProductHandler {
	return &ProductHandler{
		service:   service,
		validator: validator,
	}
}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.ListProducts)
		r.Post("/", h.CreateProduct)
		r.Get("/{id}", h.GetProduct)
		r.Put("/{id}", h.UpdateProduct)
		r.Delete("/{id}", h.DeleteProduct)
		r.Get("/category/{categoryID}", h.GetProductsByCategory)
		r.Get("/search", h.SearchProducts)
	})
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {

	params := ParseQueryParams(r)

	products, _, err := h.service.GetAll(r.Context(), params.GetLimit(), params.GetOffset())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "failed to list products"}`))
		return
	}

	RespondWithJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid product ID format"}`))
		return
	}

	product, err := h.service.GetByID(r.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "product not found"}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to get product"}`))
		}
		return
	}

	response := ProductResponse{
		ID:          product.ID.String(),
		SKU:         product.SKU,
		Name:        product.Name,
		Price:       float64(product.Price),
		Description: product.Description,
		CategoryID:  product.CategoryID.String(),
		CreatedAt:   product.CreatedAt.Format(time.RFC1123),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC1123),
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)

	var req ProductRequest
	if err := d.Decode(&req); err != nil {
		RespondWithError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	v := h.validator
	// check if category exists
	categoryID := req.Validate(h.validator)

	exists, err := h.service.CategoryExists(r.Context(), categoryID)
	if err != nil {
		RespondWithError(w, "failed to check if category exists", http.StatusInternalServerError)
		return
	}

	v.Check(exists, "category_id", "Category does not exist")

	if !v.Valid() {
		RespondWithValidationErrors(w, v.Errors)
		return
	}

	product, err := req.ToModel()
	if err != nil {
		RespondWithError(w, "invalid request data", http.StatusBadRequest)
		return
	}

	err = h.service.Create(r.Context(), product)
	if err != nil {
		RespondWithError(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	response := ProductResponseFromModel(product)

	RespondWithJSON(w, http.StatusCreated, response)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {

}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {

}

func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {

}

func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {

}
