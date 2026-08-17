  def update
    # Only what the caller sent. Strong parameters drop everything unpermitted,
    # so an omitted attribute is simply absent here and keeps its stored value.
    attributes = {{resource|record}}_params

    # Write first, then read what the write produced. Loading the row, mutating
    # it and rendering from memory would spend a read before the write and none
    # after, which describes the same change as one interaction where the
    # contract names two.
    changed = {{resource|model}}.where(id: params[:id]).update_all(
      attributes.to_h.merge("updated_at" => Time.current)
    )
    raise ActiveRecord::RecordNotFound if changed.zero?

    render json: {{resource|model}}.find(params[:id])
  end
