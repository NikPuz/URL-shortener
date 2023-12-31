package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
	"url-shortner/internal/api/dto"
	routerMiddleware "url-shortner/internal/api/middleware"
	"url-shortner/internal/entity"
)

type urlHandler struct {
	URLService entity.IURLService
	Logger     *zap.Logger
}

func RegisterURLHandlers(r *chi.Mux, service entity.IURLService, logger *zap.Logger, routerMiddleware routerMiddleware.IMiddleware) {
	URLHandler := new(urlHandler)
	URLHandler.URLService = service
	URLHandler.Logger = logger

	r.Use(routerMiddleware.PanicRecovery)
	r.Use(middleware.Timeout(time.Second * 10))
	r.Use(middleware.RequestID)
	r.Use(routerMiddleware.ContentTypeJSON)

	r.Get("/a/", routerMiddleware.DebugLogger(URLHandler.CreateShortURL))
	r.Get("/s/{shortURL:[A-z0-9]{8}}", routerMiddleware.DebugLogger(URLHandler.RedirectByShortURL))
}

func (h urlHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) ([]byte, int) {

	longURL := r.URL.Query().Get("url")
	if len(longURL) == 0 {
		resp, code := entity.HandleError(r.Context(), h.Logger, entity.URLNotEnteredError)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, entity.NotURLError)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	shortURL, err := h.URLService.CreateShortURL(r.Context(), longURL)
	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, err)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	resp, _ := json.Marshal(dto.ShortURL{ShortURL: shortURL})

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return resp, http.StatusOK
}

func (h urlHandler) RedirectByShortURL(w http.ResponseWriter, r *http.Request) ([]byte, int) {

	shortURL := chi.URLParam(r, "shortURL")
	if len(shortURL) == 0 {
		resp, code := entity.HandleError(r.Context(), h.Logger, entity.ShortURLNotEnteredError)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	longURL, err := h.URLService.GetLongURL(r.Context(), shortURL)
	if len(longURL) == 0 {
		err = entity.NotFoundError
	}

	if err != nil {
		resp, code := entity.HandleError(r.Context(), h.Logger, err)
		w.WriteHeader(code)
		w.Write(resp)
		return resp, code
	}

	http.Redirect(w, r, longURL, http.StatusFound)
	return nil, http.StatusFound
}
