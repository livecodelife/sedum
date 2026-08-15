func Create{{resource|exported}}(w http.ResponseWriter, r *http.Request) {
	var in db.{{resource|exported}}
	if err := decodeBody(r, &in); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	created, err := db.Insert{{resource|exported}}(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}
