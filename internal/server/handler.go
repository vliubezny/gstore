package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vliubezny/gstore/internal/service"
)

func (s *server) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categories, err := s.s.GetCategories(r.Context())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get categories")
		return
	}

	resp := make([]category, len(categories))

	for i, c := range categories {
		resp[i] = fromCategoryModel(c)
	}

	writeOK(l, w, resp)
}

func (s *server) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categoryID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	c, err := s.s.GetCategory(r.Context(), categoryID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "category not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to get category")
		return
	}

	writeOK(l, w, fromCategoryModel(c))
}

func (s *server) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	var req category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	c, err := s.s.CreateCategory(r.Context(), req.toModel())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to create category")
		return
	}

	writeOK(l, w, fromCategoryModel(c))
}

func (s *server) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categoryID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	var req category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	c := req.toModel()
	c.ID = categoryID

	if err := s.s.UpdateCategory(r.Context(), c); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "category not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to update category")
		return
	}

	writeOK(l, w, fromCategoryModel(c))
}

func (s *server) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categoryID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	err = s.s.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "category not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to delete category")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getStoresHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	stores, err := s.s.GetStores(r.Context())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get stores")
		return
	}

	resp := make([]store, len(stores))

	for i, str := range stores {
		resp[i] = fromStoreModel(str)
	}

	writeOK(l, w, resp)
}

func (s *server) getStoreHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	storeID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	str, err := s.s.GetStore(r.Context(), storeID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "store not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to get store")
		return
	}

	writeOK(l, w, fromStoreModel(str))
}

func (s *server) createStoreHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	var req store
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	str, err := s.s.CreateStore(r.Context(), req.toModel())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to create store")
		return
	}

	writeOK(l, w, fromStoreModel(str))
}

func (s *server) updateStoreHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	storeID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	var req store
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	str := req.toModel()
	str.ID = storeID

	if err := s.s.UpdateStore(r.Context(), str); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "store not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to update store")
		return
	}

	writeOK(l, w, fromStoreModel(str))
}

func (s *server) deleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	storeID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	err = s.s.DeleteStore(r.Context(), storeID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "store not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to delete store")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getStorePositionsHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	storeID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	positions, err := s.s.GetStorePositions(r.Context(), storeID)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get store positions")
		return
	}

	resp := make([]position, len(positions))

	for i, p := range positions {
		resp[i] = fromPositionModel(p)
	}

	writeOK(l, w, resp)
}

func (s *server) satPositionHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	storeID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	productID, err := getIDFromURL(r, "productId")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var req position
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	p := req.toModel()
	p.ProductID = productID
	p.StoreID = storeID

	if err := s.s.SetPosition(r.Context(), p); err != nil {
		switch {
		case errors.Is(err, service.ErrUnknownProduct):
			writeError(l.WithError(err), w, http.StatusNotFound, "product not found")
		case errors.Is(err, service.ErrUnknownStore):
			writeError(l.WithError(err), w, http.StatusNotFound, "store not found")
		default:
			writeInternalError(l.WithError(err), w, "fail to set position")
		}
		return
	}

	writeOK(l, w, fromPositionModel(p))
}

func (s *server) deletePositionHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	storeID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	productID, err := getIDFromURL(r, "productId")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid product ID")
		return
	}

	err = s.s.DeletePosition(r.Context(), productID, storeID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "product not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getCategoryProductsHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categoryID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	products, err := s.s.GetProducts(r.Context(), categoryID)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get products")
		return
	}

	resp := make([]product, len(products))

	for i, p := range products {
		resp[i] = fromProductModel(p)
	}

	writeOK(l, w, resp)
}

func (s *server) getProductHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	productID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid product ID")
		return
	}

	p, err := s.s.GetProduct(r.Context(), productID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "product not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to get product")
		return
	}

	writeOK(l, w, fromProductModel(p))
}

func (s *server) createProductHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	var req product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	p, err := s.s.CreateProduct(r.Context(), req.toModel())
	if err != nil {
		if errors.Is(err, service.ErrUnknownCategory) {
			writeError(l.WithError(err), w, http.StatusBadRequest, "unknown category")
		} else {
			writeInternalError(l.WithError(err), w, "fail to create product")
		}
		return
	}

	writeOK(l, w, fromProductModel(p))
}

func (s *server) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	productID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid product ID")
		return
	}

	var req product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	p := req.toModel()
	p.ID = productID

	if err := s.s.UpdateProduct(r.Context(), p); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			writeError(l.WithError(err), w, http.StatusNotFound, "product not found")
		case errors.Is(err, service.ErrUnknownCategory):
			writeError(l.WithError(err), w, http.StatusBadRequest, "unknown category")
		default:
			writeInternalError(l.WithError(err), w, "fail to update product")
		}
		return
	}

	writeOK(l, w, fromProductModel(p))
}

func (s *server) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	productID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid product ID")
		return
	}

	err = s.s.DeleteProduct(r.Context(), productID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "product not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getProductOffersHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	productID, err := getIDFromURL(r, "id")
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid product ID")
		return
	}

	positions, err := s.s.GetProductPositions(r.Context(), productID)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get product offers")
		return
	}

	resp := make([]position, len(positions))

	for i, p := range positions {
		resp[i] = fromPositionModel(p)
	}

	writeOK(l, w, resp)
}
