func List{{resource|collection}}(w http.ResponseWriter, r *http.Request) {
	all, err := db.SelectAll{{resource|collection}}(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, all)
}
